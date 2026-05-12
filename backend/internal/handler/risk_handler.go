package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type RiskHandler struct {
	riskService *service.RiskService
}

func NewRiskHandler(riskService *service.RiskService) *RiskHandler {
	return &RiskHandler{riskService: riskService}
}

func (h *RiskHandler) ListAdminRisks(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	severity := strings.TrimSpace(r.URL.Query().Get("severity"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
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

	result, err := h.riskService.GetAdminRisks(r.Context(), currentUser.Role, severity, status, start, end)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, result)
}

func (h *RiskHandler) UpdateRiskStatus(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := strconv.ParseUint(strings.TrimSpace(r.PathValue("eventId")), 10, 64)
	if err != nil || eventID == 0 {
		response.Error(w, http.StatusBadRequest, "invalid eventId")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.riskService.UpdateRiskStatus(r.Context(), currentUser.Role, uint(eventID), req.Status); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]interface{}{
		"eventId": uint(eventID),
		"status":  strings.TrimSpace(req.Status),
	})
}
