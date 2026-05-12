package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ridehailing/backend/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByID(ctx context.Context, id uint) (*model.Payment, error)
	GetLatestByOrderID(ctx context.Context, orderID uint) (*model.Payment, error)
	MarkPaid(ctx context.Context, paymentID uint) error
	ClosePendingByOrderID(ctx context.Context, orderID uint) error
}

type GormPaymentRepository struct {
	db *gorm.DB
}

func NewGormPaymentRepository(db *gorm.DB) *GormPaymentRepository {
	return &GormPaymentRepository{db: db}
}

func (r *GormPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *GormPaymentRepository) GetByID(ctx context.Context, id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).First(&payment, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *GormPaymentRepository) GetLatestByOrderID(ctx context.Context, orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC, id DESC").
		First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *GormPaymentRepository) MarkPaid(ctx context.Context, paymentID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment model.Payment
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, paymentID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		if err != nil {
			return err
		}

		if payment.Status == model.PaymentStatusPaid {
			return nil
		}
		if payment.Status != model.PaymentStatusPending {
			return errors.New("payment is not payable")
		}

		var order model.Order
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, payment.OrderID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		if err != nil {
			return err
		}

		if order.OrderStatus == model.OrderStatusCancelled {
			return errors.New("cancelled order cannot be paid")
		}
		if order.PayStatus == model.PayStatusPaid {
			return errors.New("order already paid")
		}
		if order.PaymentExpireAt != nil && time.Now().After(*order.PaymentExpireAt) {
			return errors.New("order payment expired")
		}

		now := time.Now()
		payment.Status = model.PaymentStatusPaid
		payment.PaidAt = &now

		order.PayStatus = model.PayStatusPaid
		if order.OrderStatus == model.OrderStatusPendingPayment {
			order.OrderStatus = model.OrderStatusPendingVerification
		}

		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return tx.Save(&payment).Error
	})
}

func (r *GormPaymentRepository) ClosePendingByOrderID(ctx context.Context, orderID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("order_id = ? AND status = ?", orderID, model.PaymentStatusPending).
		Update("status", model.PaymentStatusClosed).Error
}
