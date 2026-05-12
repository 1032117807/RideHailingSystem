package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type TokenUsageSummaryRow struct {
	UserID           uint
	Role             string
	Feature          string
	Model            string
	RequestCount     int
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	LastUsedAt       time.Time
}

type TokenUsageUserSummaryRow struct {
	UserID           uint
	Role             string
	RequestCount     int
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	LastUsedAt       time.Time
}

type TokenUsageRepository struct {
	db *gorm.DB
}

func NewTokenUsageRepository(db *gorm.DB) *TokenUsageRepository {
	return &TokenUsageRepository{db: db}
}

func (r *TokenUsageRepository) Create(ctx context.Context, item *model.TokenUsage) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *TokenUsageRepository) ListSummary(
	ctx context.Context,
	start time.Time,
	end time.Time,
	role string,
	feature string,
) ([]*TokenUsageSummaryRow, error) {
	var rows []*TokenUsageSummaryRow

	query := r.db.WithContext(ctx).
		Table("token_usages").
		Select(`
			user_id,
			role,
			feature,
			model,
			SUM(request_count) as request_count,
			SUM(prompt_tokens) as prompt_tokens,
			SUM(completion_tokens) as completion_tokens,
			SUM(total_tokens) as total_tokens,
			MAX(created_at) as last_used_at
		`).
		Where("created_at >= ? AND created_at < ?", start, end)

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if feature != "" {
		query = query.Where("feature = ?", feature)
	}

	err := query.
		Group("user_id, role, feature, model").
		Order("total_tokens DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *TokenUsageRepository) ListUserSummary(
	ctx context.Context,
	start time.Time,
	end time.Time,
	role string,
) ([]*TokenUsageUserSummaryRow, error) {
	var rows []*TokenUsageUserSummaryRow

	query := r.db.WithContext(ctx).
		Table("token_usages").
		Select(`
			user_id,
			role,
			SUM(request_count) as request_count,
			SUM(prompt_tokens) as prompt_tokens,
			SUM(completion_tokens) as completion_tokens,
			SUM(total_tokens) as total_tokens,
			MAX(created_at) as last_used_at
		`).
		Where("created_at >= ? AND created_at < ?", start, end)

	if role != "" {
		query = query.Where("role = ?", role)
	}

	err := query.
		Group("user_id, role").
		Order("total_tokens DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
