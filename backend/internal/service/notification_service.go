package service

import (
	"context"
	"errors"
	"strings"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type CreateNotificationInput struct {
	UserID         uint
	Type           string
	Title          string
	Content        string
	RelatedOrderID *uint
}

type NotificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) CreateSystemNotification(ctx context.Context, input CreateNotificationInput) error {
	if input.UserID == 0 {
		return errors.New("user ID is required")
	}

	input.Type = strings.TrimSpace(input.Type)
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)

	if input.Type == "" || input.Title == "" || input.Content == "" {
		return errors.New("type, title and content are required")
	}

	return s.notificationRepo.Create(ctx, &model.Notification{
		UserID:         input.UserID,
		Type:           input.Type,
		Title:          input.Title,
		Content:        input.Content,
		RelatedOrderID: input.RelatedOrderID,
		IsRead:         false,
	})
}

func (s *NotificationService) ListMyNotifications(ctx context.Context, currentUserID uint, currentUserRole string, limit int) ([]*model.Notification, error) {
	if currentUserRole != model.RolePassenger && currentUserRole != model.RoleDriver && currentUserRole != model.RoleAdmin {
		return nil, errors.New("invalid role")
	}
	return s.notificationRepo.ListByUserID(ctx, currentUserID, limit)
}

func (s *NotificationService) CountMyUnreadNotifications(ctx context.Context, currentUserID uint, currentUserRole string) (int, error) {
	if currentUserRole != model.RolePassenger && currentUserRole != model.RoleDriver && currentUserRole != model.RoleAdmin {
		return 0, errors.New("invalid role")
	}
	return s.notificationRepo.CountUnreadByUserID(ctx, currentUserID)
}

func (s *NotificationService) MarkMyNotificationRead(ctx context.Context, currentUserID uint, currentUserRole string, notificationID uint) error {
	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		return errors.New("notification not found")
	}
	if notification.UserID != currentUserID {
		return errors.New("notification does not belong to current user")
	}
	return s.notificationRepo.MarkRead(ctx, notificationID)
}

func (s *NotificationService) MarkAllMyNotificationsRead(ctx context.Context, currentUserID uint, currentUserRole string) error {
	if currentUserRole != model.RolePassenger && currentUserRole != model.RoleDriver && currentUserRole != model.RoleAdmin {
		return errors.New("invalid role")
	}
	return s.notificationRepo.MarkAllReadByUserID(ctx, currentUserID)
}
