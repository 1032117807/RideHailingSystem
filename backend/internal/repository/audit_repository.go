package repository

import (
	"context"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type AuditRepository interface {
	CreateAuditLog(ctx context.Context, log *model.AuditLog) error
	ListAuditLogs(ctx context.Context, limit int) ([]*model.AuditLog, error)
	CreateRefundAuditLog(ctx context.Context, log *model.RefundAuditLog) error
	ListRefundAuditLogs(ctx context.Context, orderID uint) ([]*model.RefundAuditLog, error)
}

type GormAuditRepository struct {
	db *gorm.DB
}

func NewGormAuditRepository(db *gorm.DB) *GormAuditRepository {
	return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *GormAuditRepository) ListAuditLogs(ctx context.Context, limit int) ([]*model.AuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	var logs []*model.AuditLog
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *GormAuditRepository) CreateRefundAuditLog(ctx context.Context, log *model.RefundAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *GormAuditRepository) ListRefundAuditLogs(ctx context.Context, orderID uint) ([]*model.RefundAuditLog, error) {
	var logs []*model.RefundAuditLog
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
