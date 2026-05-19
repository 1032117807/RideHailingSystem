package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	limit, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("limit")))
	logs, err := h.auditService.ListAuditLogs(r.Context(), currentUser.Role, limit)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, logs)
}

func (h *AuditHandler) ListRefundAuditLogs(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	orderID, err := strconv.ParseUint(strings.TrimSpace(r.PathValue("orderId")), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}
	logs, err := h.auditService.ListRefundAuditLogs(r.Context(), currentUser.Role, uint(orderID))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, logs)
}
