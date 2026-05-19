package service

import (
	"context"
	"errors"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type AuditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) ListAuditLogs(ctx context.Context, currentUserRole string, limit int) ([]*model.AuditLog, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view audit logs")
	}
	return s.auditRepo.ListAuditLogs(ctx, limit)
}

func (s *AuditService) ListRefundAuditLogs(ctx context.Context, currentUserRole string, orderID uint) ([]*model.RefundAuditLog, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view refund audit logs")
	}
	return s.auditRepo.ListRefundAuditLogs(ctx, orderID)
}
