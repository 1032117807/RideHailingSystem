package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type RiskEventSummaryRow struct {
	Severity string
	Count    int
}

type RiskEventListFilter struct {
	Severity string
	Status   string
	Start    time.Time
	End      time.Time
}

type RiskEventRepository struct {
	db *gorm.DB
}

func NewRiskEventRepository(db *gorm.DB) *RiskEventRepository {
	return &RiskEventRepository{db: db}
}

func (r *RiskEventRepository) Create(ctx context.Context, event *model.RiskEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *RiskEventRepository) FindOpenByFingerprintSince(ctx context.Context, fingerprint string, since time.Time) (*model.RiskEvent, error) {
	var item model.RiskEvent
	err := r.db.WithContext(ctx).
		Where("fingerprint = ? AND status = ? AND created_at >= ?", fingerprint, model.RiskStatusOpen, since).
		Order("id desc").
		First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *RiskEventRepository) ListForAdmin(ctx context.Context, filter RiskEventListFilter) ([]*model.RiskEvent, error) {
	var items []*model.RiskEvent
	query := r.db.WithContext(ctx).
		Model(&model.RiskEvent{}).
		Where("created_at >= ? AND created_at < ?", filter.Start, filter.End).
		Order("created_at desc")

	if filter.Severity != "" {
		query = query.Where("severity = ?", filter.Severity)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Limit(100).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RiskEventRepository) CountBySeverity(ctx context.Context, start time.Time, end time.Time) ([]*RiskEventSummaryRow, error) {
	var rows []*RiskEventSummaryRow
	err := r.db.WithContext(ctx).
		Model(&model.RiskEvent{}).
		Select("severity, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("severity").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *RiskEventRepository) CountOpen(ctx context.Context, start time.Time, end time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RiskEvent{}).
		Where("status = ? AND created_at >= ? AND created_at < ?", model.RiskStatusOpen, start, end).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RiskEventRepository) UpdateStatus(ctx context.Context, eventID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.RiskEvent{}).
		Where("id = ?", eventID).
		Update("status", status).
		Error
}
