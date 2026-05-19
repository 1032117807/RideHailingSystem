package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ridehailing/backend/internal/model"
)

type ElectronicTicketRepository interface {
	Create(ctx context.Context, ticket *model.ElectronicTicket) error
	GetByOrderID(ctx context.Context, orderID uint) (*model.ElectronicTicket, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.ElectronicTicket, error)
	MarkVerified(ctx context.Context, ticketID uint, driverID uint, verifiedAt time.Time) (*model.ElectronicTicket, error)
	VoidByOrderID(ctx context.Context, orderID uint) error
	CreateVerification(ctx context.Context, verification *model.TicketVerification) error
	ListVerificationsByOrderID(ctx context.Context, orderID uint) ([]*model.TicketVerification, error)
}

type GormElectronicTicketRepository struct {
	db *gorm.DB
}

func NewGormElectronicTicketRepository(db *gorm.DB) *GormElectronicTicketRepository {
	return &GormElectronicTicketRepository{db: db}
}

func (r *GormElectronicTicketRepository) Create(ctx context.Context, ticket *model.ElectronicTicket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *GormElectronicTicketRepository) GetByOrderID(ctx context.Context, orderID uint) (*model.ElectronicTicket, error) {
	var ticket model.ElectronicTicket
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Trip").
		Where("order_id = ?", orderID).
		First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *GormElectronicTicketRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.ElectronicTicket, error) {
	var ticket model.ElectronicTicket
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Trip").
		Where("token_hash = ?", tokenHash).
		First(&ticket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *GormElectronicTicketRepository) MarkVerified(ctx context.Context, ticketID uint, driverID uint, verifiedAt time.Time) (*model.ElectronicTicket, error) {
	var updated model.ElectronicTicket
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ticket model.ElectronicTicket
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&ticket, ticketID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("electronic ticket not found")
		}
		if err != nil {
			return err
		}
		if ticket.Status == model.ElectronicTicketStatusVerified {
			updated = ticket
			return nil
		}
		if ticket.Status != model.ElectronicTicketStatusIssued {
			return errors.New("electronic ticket is not valid")
		}

		ticket.Status = model.ElectronicTicketStatusVerified
		ticket.VerifiedAt = &verifiedAt
		ticket.VerifiedByDriverID = &driverID
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Order{}).
			Where("id = ? AND order_status = ?", ticket.OrderID, model.OrderStatusPendingVerification).
			Update("order_status", model.OrderStatusCompleted).Error; err != nil {
			return err
		}

		updated = ticket
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *GormElectronicTicketRepository) CreateVerification(ctx context.Context, verification *model.TicketVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

func (r *GormElectronicTicketRepository) VoidByOrderID(ctx context.Context, orderID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.ElectronicTicket{}).
		Where("order_id = ? AND status = ?", orderID, model.ElectronicTicketStatusIssued).
		Update("status", model.ElectronicTicketStatusVoided).Error
}

func (r *GormElectronicTicketRepository) ListVerificationsByOrderID(ctx context.Context, orderID uint) ([]*model.TicketVerification, error) {
	var records []*model.TicketVerification
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
