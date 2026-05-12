package handler

import (
	"net/http"
	"strings"
	"time"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type TokenUsageHandler struct {
	tokenUsageService *service.TokenUsageService
}

func NewTokenUsageHandler(tokenUsageService *service.TokenUsageService) *TokenUsageHandler {
	return &TokenUsageHandler{tokenUsageService: tokenUsageService}
}

func (h *TokenUsageHandler) ListAdminUsage(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := strings.TrimSpace(r.URL.Query().Get("role"))
	feature := strings.TrimSpace(r.URL.Query().Get("feature"))
	rangeDays := 30

	if rawDays := strings.TrimSpace(r.URL.Query().Get("days")); rawDays != "" {
		if rawDays == "7" {
			rangeDays = 7
		} else if rawDays == "90" {
			rangeDays = 90
		}
	}

	end := time.Now()
	start := end.AddDate(0, 0, -rangeDays)

	result, err := h.tokenUsageService.GetAdminUsage(
		r.Context(),
		currentUser.Role,
		start,
		end,
		role,
		feature,
	)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, result)
}
