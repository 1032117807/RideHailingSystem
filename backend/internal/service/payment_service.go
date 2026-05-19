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

type CreatePaymentInput struct {
	OrderID uint
	Channel string
}

type PaymentService struct {
	paymentRepo   repository.PaymentRepository
	orderRepo     repository.OrderRepository
	ticketService *ElectronicTicketService
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, ticketService *ElectronicTicketService) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo, ticketService: ticketService}
}

func (s *PaymentService) CreatePayment(ctx context.Context, currentUserID uint, currentUserRole string, input CreatePaymentInput) (*model.Payment, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can create payment")
	}
	if input.OrderID == 0 {
		return nil, errors.New("order ID is required")
	}

	order, err := s.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID {
		return nil, errors.New("order does not belong to current user")
	}
	if order.OrderStatus == model.OrderStatusCancelled {
		return nil, errors.New("cancelled order cannot be paid")
	}
	if order.PayStatus == model.PayStatusPaid {
		return nil, errors.New("order already paid")
	}
	if order.PaymentExpireAt != nil && time.Now().After(*order.PaymentExpireAt) {
		return nil, errors.New("order payment expired")
	}

	existing, err := s.paymentRepo.GetLatestByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == model.PaymentStatusPending {
		return existing, nil
	}

	channel := strings.TrimSpace(input.Channel)
	if channel == "" {
		channel = model.PaymentChannelMock
	}

	payment := &model.Payment{
		PaymentNo: generatePaymentNo(),
		OrderID:   order.ID,
		UserID:    currentUserID,
		Amount:    order.Amount,
		Channel:   channel,
		Status:    model.PaymentStatusPending,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}
	return s.paymentRepo.GetByID(ctx, payment.ID)
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, currentUserID uint, currentUserRole string, paymentID uint) (*model.Payment, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can view payment status")
	}

	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	if payment.UserID != currentUserID {
		return nil, errors.New("payment does not belong to current user")
	}
	return payment, nil
}

func (s *PaymentService) MockPaySuccess(ctx context.Context, currentUserID uint, currentUserRole string, paymentID uint) (*model.Payment, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can complete mock payment")
	}

	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	if payment.UserID != currentUserID {
		return nil, errors.New("payment does not belong to current user")
	}

	order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.PaymentExpireAt != nil && time.Now().After(*order.PaymentExpireAt) {
		return nil, errors.New("order payment expired")
	}

	if err := s.paymentRepo.MarkPaid(ctx, paymentID); err != nil {
		return nil, err
	}
	if s.ticketService != nil {
		paidOrder, orderErr := s.orderRepo.GetByID(ctx, payment.OrderID)
		if orderErr != nil {
			return nil, orderErr
		}
		if paidOrder != nil {
			if _, ticketErr := s.ticketService.EnsureForOrder(ctx, paidOrder); ticketErr != nil {
				return nil, ticketErr
			}
		}
	}
	return s.paymentRepo.GetByID(ctx, paymentID)
}

func generatePaymentNo() string {
	return fmt.Sprintf("PAY%d", time.Now().UnixNano())
}
