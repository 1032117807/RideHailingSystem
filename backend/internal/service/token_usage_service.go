package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type RecordTokenUsageInput struct {
	UserID           uint
	Role             string
	Feature          string
	RequestKind      string
	Provider         string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type AdminTokenUsageItem struct {
	UserID           uint      `json:"userId"`
	Role             string    `json:"role"`
	Feature          string    `json:"feature"`
	Model            string    `json:"model"`
	RequestCount     int       `json:"requestCount"`
	PromptTokens     int       `json:"promptTokens"`
	CompletionTokens int       `json:"completionTokens"`
	TotalTokens      int       `json:"totalTokens"`
	LastUsedAt       time.Time `json:"lastUsedAt"`
}

type AdminTokenUserItem struct {
	UserID           uint      `json:"userId"`
	Role             string    `json:"role"`
	RequestCount     int       `json:"requestCount"`
	PromptTokens     int       `json:"promptTokens"`
	CompletionTokens int       `json:"completionTokens"`
	TotalTokens      int       `json:"totalTokens"`
	LastUsedAt       time.Time `json:"lastUsedAt"`
}

type AdminTokenUsageResult struct {
	Start string                 `json:"start"`
	End   string                 `json:"end"`
	Role  string                 `json:"role"`
	Items []*AdminTokenUsageItem `json:"items"`
	Users []*AdminTokenUserItem  `json:"users"`
}

type TokenUsageService struct {
	tokenUsageRepo *repository.TokenUsageRepository
}

func NewTokenUsageService(tokenUsageRepo *repository.TokenUsageRepository) *TokenUsageService {
	return &TokenUsageService{tokenUsageRepo: tokenUsageRepo}
}

func (s *TokenUsageService) Record(ctx context.Context, input RecordTokenUsageInput) error {
	if input.UserID == 0 {
		return nil
	}
	input.Role = strings.TrimSpace(input.Role)
	input.Feature = strings.TrimSpace(input.Feature)
	input.RequestKind = strings.TrimSpace(input.RequestKind)
	input.Provider = strings.TrimSpace(input.Provider)
	input.Model = strings.TrimSpace(input.Model)

	if input.Role == "" || input.Feature == "" || input.RequestKind == "" || input.Model == "" {
		return nil
	}
	if input.TotalTokens <= 0 && input.PromptTokens <= 0 && input.CompletionTokens <= 0 {
		return nil
	}
	if input.TotalTokens <= 0 {
		input.TotalTokens = input.PromptTokens + input.CompletionTokens
	}

	return s.tokenUsageRepo.Create(ctx, &model.TokenUsage{
		UserID:           input.UserID,
		Role:             input.Role,
		Feature:          input.Feature,
		RequestKind:      input.RequestKind,
		Provider:         input.Provider,
		Model:            input.Model,
		PromptTokens:     input.PromptTokens,
		CompletionTokens: input.CompletionTokens,
		TotalTokens:      input.TotalTokens,
		RequestCount:     1,
	})
}

func (s *TokenUsageService) GetAdminUsage(
	ctx context.Context,
	currentUserRole string,
	start time.Time,
	end time.Time,
	role string,
	feature string,
) (*AdminTokenUsageResult, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view token usage")
	}

	rows, err := s.tokenUsageRepo.ListSummary(ctx, start, end, strings.TrimSpace(role), strings.TrimSpace(feature))
	if err != nil {
		return nil, err
	}
	userRows, err := s.tokenUsageRepo.ListUserSummary(ctx, start, end, strings.TrimSpace(role))
	if err != nil {
		return nil, err
	}

	items := make([]*AdminTokenUsageItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &AdminTokenUsageItem{
			UserID:           row.UserID,
			Role:             row.Role,
			Feature:          row.Feature,
			Model:            row.Model,
			RequestCount:     row.RequestCount,
			PromptTokens:     row.PromptTokens,
			CompletionTokens: row.CompletionTokens,
			TotalTokens:      row.TotalTokens,
			LastUsedAt:       row.LastUsedAt,
		})
	}

	users := make([]*AdminTokenUserItem, 0, len(userRows))
	for _, row := range userRows {
		users = append(users, &AdminTokenUserItem{
			UserID:           row.UserID,
			Role:             row.Role,
			RequestCount:     row.RequestCount,
			PromptTokens:     row.PromptTokens,
			CompletionTokens: row.CompletionTokens,
			TotalTokens:      row.TotalTokens,
			LastUsedAt:       row.LastUsedAt,
		})
	}

	return &AdminTokenUsageResult{
		Start: start.Format(time.RFC3339),
		End:   end.Format(time.RFC3339),
		Role:  role,
		Items: items,
		Users: users,
	}, nil
}
