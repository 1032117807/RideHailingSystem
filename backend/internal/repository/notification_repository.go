package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	ListByUserID(ctx context.Context, userID uint, limit int) ([]*model.Notification, error)
	CountUnreadByUserID(ctx context.Context, userID uint) (int, error)
	GetByID(ctx context.Context, id uint) (*model.Notification, error)
	MarkRead(ctx context.Context, id uint) error
	MarkAllReadByUserID(ctx context.Context, userID uint) error
}

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewGormNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *GormNotificationRepository) ListByUserID(ctx context.Context, userID uint, limit int) ([]*model.Notification, error) {
	var notifications []*model.Notification

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_read ASC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *GormNotificationRepository) CountUnreadByUserID(ctx context.Context, userID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GormNotificationRepository) GetByID(ctx context.Context, id uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.WithContext(ctx).First(&notification, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *GormNotificationRepository) MarkRead(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *GormNotificationRepository) MarkAllReadByUserID(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read": true,
			"read_at": &now,
		}).Error
}
