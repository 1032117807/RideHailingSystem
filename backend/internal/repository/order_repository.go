package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ridehailing/backend/internal/model"
)

type OrderRepository interface {
	CreateWithSeatDeduction(ctx context.Context, order *model.Order) error
	ListByUserID(ctx context.Context, userID uint) ([]*model.Order, error)
	ListByTripID(ctx context.Context, tripID uint) ([]*model.Order, error)
	ListForAdmin(ctx context.Context, refundStatus string) ([]*model.Order, error)
	CountRefundSummary(ctx context.Context) (int, int, int, error)
	GetByID(ctx context.Context, id uint) (*model.Order, error)
	CountSummaryByTripID(ctx context.Context, tripID uint) (int, int, error)
	ListExpiredPendingOrders(ctx context.Context, now time.Time, limit int) ([]*model.Order, error)
	ExpirePendingOrder(ctx context.Context, orderID uint, now time.Time) (bool, error)
	CancelAndReleaseSeats(ctx context.Context, order *model.Order) error
	RequestRefund(ctx context.Context, order *model.Order) error
	MarkCompleted(ctx context.Context, order *model.Order) error
	ApproveRefund(ctx context.Context, order *model.Order, reviewNote string) error
	RejectRefund(ctx context.Context, order *model.Order, reviewNote string) error
}

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateWithSeatDeduction(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var trip model.Trip
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, order.TripID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("trip not found")
		}
		if err != nil {
			return err
		}

		if trip.Status != model.TripStatusPublished {
			return errors.New("trip is not available")
		}
		if order.TicketCount <= 0 {
			return errors.New("ticket count must be greater than 0")
		}
		if trip.SeatAvailable < order.TicketCount {
			return errors.New("not enough seats available")
		}

		trip.SeatAvailable -= order.TicketCount
		if err := tx.Save(&trip).Error; err != nil {
			return err
		}

		return tx.Create(order).Error
	})
}

func (r *GormOrderRepository) ListByUserID(ctx context.Context, userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Trip").
		Preload("Trip.Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order ASC")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) ListByTripID(ctx context.Context, tripID uint) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Trip").
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) ListForAdmin(ctx context.Context, refundStatus string) ([]*model.Order, error) {
	var orders []*model.Order
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Trip").
		Order("updated_at DESC, created_at DESC")

	if refundStatus != "" {
		query = query.Where("refund_status = ?", refundStatus)
	}

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) CountRefundSummary(ctx context.Context) (int, int, int, error) {
	var requested int64
	if err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("refund_status = ?", model.RefundStatusRequested).
		Count(&requested).Error; err != nil {
		return 0, 0, 0, err
	}

	var refunded int64
	if err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("refund_status = ?", model.RefundStatusRefunded).
		Count(&refunded).Error; err != nil {
		return 0, 0, 0, err
	}

	var rejected int64
	if err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("refund_status = ?", model.RefundStatusRejected).
		Count(&rejected).Error; err != nil {
		return 0, 0, 0, err
	}

	return int(requested), int(refunded), int(rejected), nil
}

func (r *GormOrderRepository) GetByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Trip").
		Preload("Trip.Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order ASC")
		}).
		First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderRepository) CountSummaryByTripID(ctx context.Context, tripID uint) (int, int, error) {
	var pendingVerificationCount int64
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("trip_id = ? AND order_status = ?", tripID, model.OrderStatusPendingVerification).
		Select("COALESCE(SUM(ticket_count), 0)").
		Scan(&pendingVerificationCount).Error
	if err != nil {
		return 0, 0, err
	}

	var refundRequestCount int64
	err = r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("trip_id = ? AND refund_status = ?", tripID, model.RefundStatusRequested).
		Count(&refundRequestCount).Error
	if err != nil {
		return 0, 0, err
	}

	return int(pendingVerificationCount), int(refundRequestCount), nil
}

func (r *GormOrderRepository) ListExpiredPendingOrders(ctx context.Context, now time.Time, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	query := r.db.WithContext(ctx).
		Where("order_status = ?", model.OrderStatusPendingPayment).
		Where("pay_status = ?", model.PayStatusUnpaid).
		Where("payment_expire_at IS NOT NULL").
		Where("payment_expire_at <= ?", now).
		Order("payment_expire_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) ExpirePendingOrder(ctx context.Context, orderID uint, now time.Time) (bool, error) {
	expired := false

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, orderID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		if order.OrderStatus != model.OrderStatusPendingPayment || order.PayStatus != model.PayStatusUnpaid {
			return nil
		}
		if order.PaymentExpireAt == nil || order.PaymentExpireAt.After(now) {
			return nil
		}

		var trip model.Trip
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, order.TripID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("trip not found")
		}
		if err != nil {
			return err
		}

		order.OrderStatus = model.OrderStatusCancelled
		order.RefundStatus = model.RefundStatusNone

		trip.SeatAvailable += order.TicketCount
		if trip.SeatAvailable > trip.SeatTotal {
			trip.SeatAvailable = trip.SeatTotal
		}

		if err := tx.Save(&trip).Error; err != nil {
			return err
		}
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		expired = true
		return nil
	})

	return expired, err
}

func (r *GormOrderRepository) CancelAndReleaseSeats(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var trip model.Trip
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, order.TripID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("trip not found")
		}
		if err != nil {
			return err
		}

		order.OrderStatus = model.OrderStatusCancelled
		order.RefundStatus = model.RefundStatusNone

		trip.SeatAvailable += order.TicketCount
		if trip.SeatAvailable > trip.SeatTotal {
			trip.SeatAvailable = trip.SeatTotal
		}

		if err := tx.Save(&trip).Error; err != nil {
			return err
		}
		return tx.Save(order).Error
	})
}

func (r *GormOrderRepository) RequestRefund(ctx context.Context, order *model.Order) error {
	order.RefundStatus = model.RefundStatusRequested
	order.RefundReviewNote = ""
	order.RefundReviewedAt = nil
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *GormOrderRepository) MarkCompleted(ctx context.Context, order *model.Order) error {
	order.OrderStatus = model.OrderStatusCompleted
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *GormOrderRepository) ApproveRefund(ctx context.Context, order *model.Order, reviewNote string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedOrder model.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedOrder, order.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		if err != nil {
			return err
		}

		if lockedOrder.RefundStatus != model.RefundStatusRequested {
			return errors.New("order is not waiting for refund review")
		}

		now := time.Now()
		lockedOrder.RefundStatus = model.RefundStatusRefunded
		lockedOrder.RefundReviewNote = reviewNote
		lockedOrder.RefundReviewedAt = &now

		if lockedOrder.OrderStatus == model.OrderStatusPendingVerification {
			var trip model.Trip
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, lockedOrder.TripID).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("trip not found")
			}
			if err != nil {
				return err
			}

			trip.SeatAvailable += lockedOrder.TicketCount
			if trip.SeatAvailable > trip.SeatTotal {
				trip.SeatAvailable = trip.SeatTotal
			}
			lockedOrder.OrderStatus = model.OrderStatusCancelled

			if err := tx.Save(&trip).Error; err != nil {
				return err
			}
		}

		return tx.Save(&lockedOrder).Error
	})
}

func (r *GormOrderRepository) RejectRefund(ctx context.Context, order *model.Order, reviewNote string) error {
	now := time.Now()
	order.RefundStatus = model.RefundStatusRejected
	order.RefundReviewNote = reviewNote
	order.RefundReviewedAt = &now
	return r.db.WithContext(ctx).Save(order).Error
}
