package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type CreateOrderInput struct {
	TripID      uint
	TicketCount int
	SeatType    string
}

const orderPaymentExpireDuration = 15 * time.Minute

type AdminDashboardSummary struct {
	TotalUsers          int `json:"totalUsers"`
	PassengerCount      int `json:"passengerCount"`
	DriverCount         int `json:"driverCount"`
	AdminCount          int `json:"adminCount"`
	ActiveCount         int `json:"activeCount"`
	FrozenCount         int `json:"frozenCount"`
	DisabledCount       int `json:"disabledCount"`
	VerifiedCount       int `json:"verifiedCount"`
	PendingRefundCount  int `json:"pendingRefundCount"`
	RefundedCount       int `json:"refundedCount"`
	RejectedRefundCount int `json:"rejectedRefundCount"`
}

type OrderService struct {
	orderRepo        repository.OrderRepository
	tripRepo         repository.TripRepository
	userRepo         repository.UserRepository
	paymentRepo      repository.PaymentRepository
	notificationRepo repository.NotificationRepository
	auditRepo        repository.AuditRepository
	ticketService    *ElectronicTicketService
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	tripRepo repository.TripRepository,
	userRepo repository.UserRepository,
	paymentRepo repository.PaymentRepository,
	notificationRepo repository.NotificationRepository,
	auditRepo repository.AuditRepository,
	ticketService *ElectronicTicketService,
) *OrderService {
	return &OrderService{
		orderRepo:        orderRepo,
		tripRepo:         tripRepo,
		userRepo:         userRepo,
		paymentRepo:      paymentRepo,
		notificationRepo: notificationRepo,
		auditRepo:        auditRepo,
		ticketService:    ticketService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, currentUserID uint, currentUserRole string, input CreateOrderInput) (*model.Order, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can create order")
	}
	if input.TicketCount <= 0 {
		return nil, errors.New("ticket count must be greater than 0")
	}
	if input.TripID == 0 {
		return nil, errors.New("trip ID is required")
	}

	trip, err := s.tripRepo.GetByID(ctx, input.TripID)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, errors.New("trip not found")
	}
	if trip.Status != model.TripStatusPublished {
		return nil, errors.New("trip is not available")
	}
	if trip.SeatAvailable < input.TicketCount {
		return nil, errors.New("not enough seats available")
	}

	seatType := strings.TrimSpace(input.SeatType)
	if seatType == "" {
		seatType = model.SeatTypeStandard
	}

	expireAt := time.Now().Add(orderPaymentExpireDuration)

	order := &model.Order{
		OrderNo:         generateOrderNo(),
		UserID:          currentUserID,
		TripID:          input.TripID,
		TicketCount:     input.TicketCount,
		SeatType:        seatType,
		Amount:          trip.PriceCent * input.TicketCount,
		PayStatus:       model.PayStatusUnpaid,
		OrderStatus:     model.OrderStatusPendingPayment,
		RefundStatus:    model.RefundStatusNone,
		PaymentExpireAt: &expireAt,
	}

	if err := s.orderRepo.CreateWithSeatDeduction(ctx, order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(ctx, order.ID)
}

func (s *OrderService) ListMyOrders(ctx context.Context, currentUserID uint, currentUserRole string) ([]*model.Order, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can view own orders")
	}
	return s.orderRepo.ListByUserID(ctx, currentUserID)
}

func (s *OrderService) GetMyOrderDetail(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint) (*model.Order, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can view own order detail")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID {
		return nil, errors.New("order does not belong to current user")
	}
	return order, nil
}

func (s *OrderService) CancelMyOrder(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint) (*model.Order, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can cancel own order")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID {
		return nil, errors.New("order does not belong to current user")
	}
	if order.OrderStatus != model.OrderStatusPendingPayment {
		return nil, errors.New("only pending payment order can be cancelled")
	}

	if err := s.orderRepo.CancelAndReleaseSeats(ctx, order); err != nil {
		return nil, err
	}
	return s.orderRepo.GetByID(ctx, orderID)
}

func (s *OrderService) RequestRefund(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint) (*model.Order, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can request refund")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID {
		return nil, errors.New("order does not belong to current user")
	}
	if order.PayStatus != model.PayStatusPaid {
		return nil, errors.New("only paid order can request refund")
	}
	if order.OrderStatus != model.OrderStatusPendingVerification && order.OrderStatus != model.OrderStatusCompleted {
		return nil, errors.New("current order does not support refund request")
	}
	if order.OrderStatus == model.OrderStatusCancelled {
		return nil, errors.New("cancelled order does not support refund request")
	}
	if order.RefundStatus == model.RefundStatusRequested {
		return nil, errors.New("refund already requested")
	}
	if order.RefundStatus == model.RefundStatusRefunded {
		return nil, errors.New("order already refunded")
	}

	if err := s.orderRepo.RequestRefund(ctx, order); err != nil {
		return nil, err
	}
	return s.orderRepo.GetByID(ctx, orderID)
}

func (s *OrderService) VerifyOrderByDriver(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint) (*model.Order, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can verify orders")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.Trip.ID == 0 {
		trip, tripErr := s.tripRepo.GetByID(ctx, order.TripID)
		if tripErr != nil {
			return nil, tripErr
		}
		if trip == nil {
			return nil, errors.New("trip not found")
		}
		order.Trip = *trip
	}
	if order.Trip.DriverID != currentUserID {
		return nil, errors.New("order does not belong to current driver trip")
	}
	if order.PayStatus != model.PayStatusPaid {
		return nil, errors.New("only paid order can be verified")
	}
	if order.OrderStatus != model.OrderStatusPendingVerification {
		return nil, errors.New("only pending verification order can be verified")
	}
	if order.RefundStatus == model.RefundStatusRequested || order.RefundStatus == model.RefundStatusRefunded {
		return nil, errors.New("refund processing order cannot be verified")
	}

	if err := s.orderRepo.MarkCompleted(ctx, order); err != nil {
		return nil, err
	}
	return s.orderRepo.GetByID(ctx, orderID)
}

func (s *OrderService) CloseExpiredPendingOrders(ctx context.Context) (int, error) {
	now := time.Now()

	orders, err := s.orderRepo.ListExpiredPendingOrders(ctx, now, 100)
	if err != nil {
		return 0, err
	}

	closedCount := 0
	for _, order := range orders {
		if order == nil {
			continue
		}

		expired, err := s.orderRepo.ExpirePendingOrder(ctx, order.ID, now)
		if err != nil {
			return closedCount, err
		}
		if !expired {
			continue
		}
		if s.paymentRepo != nil {
			if err := s.paymentRepo.ClosePendingByOrderID(ctx, order.ID); err != nil {
				return closedCount, err
			}
		}

		orderID := order.ID
		s.createPassengerNotification(
			ctx,
			order.UserID,
			model.NotificationTypeOrderExpired,
			"订单已超时取消",
			"你的订单未在规定时间内完成支付，系统已自动取消该订单并释放座位。",
			&orderID,
		)

		closedCount++
	}

	return closedCount, nil
}

func (s *OrderService) ListRefundOrdersForAdmin(ctx context.Context, currentUserID uint, currentUserRole string, refundStatus string) ([]*model.Order, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view refund review list")
	}
	return s.orderRepo.ListForAdmin(ctx, refundStatus)
}

func (s *OrderService) GetAdminDashboardSummary(ctx context.Context, currentUserID uint, currentUserRole string) (*AdminDashboardSummary, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view admin dashboard summary")
	}

	userSummary, err := s.userRepo.CountSummary(ctx)
	if err != nil {
		return nil, err
	}

	requested, refunded, rejected, err := s.orderRepo.CountRefundSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &AdminDashboardSummary{
		TotalUsers:          userSummary.TotalUsers,
		PassengerCount:      userSummary.PassengerCount,
		DriverCount:         userSummary.DriverCount,
		AdminCount:          userSummary.AdminCount,
		ActiveCount:         userSummary.ActiveCount,
		FrozenCount:         userSummary.FrozenCount,
		DisabledCount:       userSummary.DisabledCount,
		VerifiedCount:       userSummary.VerifiedCount,
		PendingRefundCount:  requested,
		RefundedCount:       refunded,
		RejectedRefundCount: rejected,
	}, nil
}

func (s *OrderService) ApproveRefundByAdmin(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint, reviewNote string) (*model.Order, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can approve refunds")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.RefundStatus != model.RefundStatusRequested {
		return nil, errors.New("order is not waiting for refund review")
	}

	reviewNote = strings.TrimSpace(reviewNote)
	if reviewNote == "" {
		reviewNote = "管理员已批准退款申请"
	}

	if err := s.orderRepo.ApproveRefund(ctx, order, reviewNote); err != nil {
		return nil, err
	}
	if s.ticketService != nil {
		_ = s.ticketService.VoidByOrderID(ctx, order.ID)
	}
	s.createRefundAuditLog(ctx, order.ID, model.RefundStatusRefunded, currentUserID, reviewNote)
	s.createAuditLog(ctx, currentUserID, currentUserRole, "refund.approve", "order", fmt.Sprintf("%d", order.ID), reviewNote)

	updatedOrder, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	relatedOrderID := updatedOrder.ID
	s.createPassengerNotification(
		ctx,
		updatedOrder.UserID,
		model.NotificationTypeRefundApproved,
		"退款审核已通过",
		"你的退款申请已审核通过，退款状态已更新，你可以在订单详情页查看最新结果。",
		&relatedOrderID,
	)

	return updatedOrder, nil
}

func (s *OrderService) RejectRefundByAdmin(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint, reviewNote string) (*model.Order, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can reject refunds")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.RefundStatus != model.RefundStatusRequested {
		return nil, errors.New("order is not waiting for refund review")
	}

	reviewNote = strings.TrimSpace(reviewNote)
	if reviewNote == "" {
		reviewNote = "管理员驳回了本次退款申请"
	}

	if err := s.orderRepo.RejectRefund(ctx, order, reviewNote); err != nil {
		return nil, err
	}
	s.createRefundAuditLog(ctx, order.ID, model.RefundStatusRejected, currentUserID, reviewNote)
	s.createAuditLog(ctx, currentUserID, currentUserRole, "refund.reject", "order", fmt.Sprintf("%d", order.ID), reviewNote)

	updatedOrder, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	relatedOrderID := updatedOrder.ID
	s.createPassengerNotification(
		ctx,
		updatedOrder.UserID,
		model.NotificationTypeRefundRejected,
		"退款申请被驳回",
		fmt.Sprintf("你的退款申请未通过审核。审核备注：%s", reviewNote),
		&relatedOrderID,
	)

	return updatedOrder, nil
}

func (s *OrderService) createRefundAuditLog(ctx context.Context, orderID uint, refundStatus string, reviewerID uint, reviewNote string) {
	if s.auditRepo == nil {
		return
	}
	_ = s.auditRepo.CreateRefundAuditLog(ctx, &model.RefundAuditLog{
		OrderID:      orderID,
		RefundStatus: refundStatus,
		ReviewNote:   reviewNote,
		ReviewerID:   reviewerID,
	})
}

func (s *OrderService) createAuditLog(ctx context.Context, actorUserID uint, actorRole, action, resourceType, resourceID, detail string) {
	if s.auditRepo == nil {
		return
	}
	_ = s.auditRepo.CreateAuditLog(ctx, &model.AuditLog{
		ActorUserID:  actorUserID,
		ActorRole:    actorRole,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detail,
	})
}

func (s *OrderService) createPassengerNotification(ctx context.Context, userID uint, typ, title, content string, relatedOrderID *uint) {
	if s.notificationRepo == nil || userID == 0 {
		return
	}

	_ = s.notificationRepo.Create(ctx, &model.Notification{
		UserID:         userID,
		Type:           typ,
		Title:          title,
		Content:        content,
		RelatedOrderID: relatedOrderID,
		IsRead:         false,
	})
}

func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}
