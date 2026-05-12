package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) ListMyNotifications(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := 20
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	notifications, err := h.notificationService.ListMyNotifications(r.Context(), currentUser.ID, currentUser.Role, limit)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, notifications)
}

func (h *NotificationHandler) CountMyUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	count, err := h.notificationService.CountMyUnreadNotifications(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]any{
		"unreadCount": count,
	})
}

func (h *NotificationHandler) MarkMyNotificationRead(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notificationID, err := parseNotificationPathUint(r.PathValue("notificationId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid notificationId")
		return
	}

	if err := h.notificationService.MarkMyNotificationRead(r.Context(), currentUser.ID, currentUser.Role, notificationID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]any{
		"message": "notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllMyNotificationsRead(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.notificationService.MarkAllMyNotificationsRead(r.Context(), currentUser.ID, currentUser.Role); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]any{
		"message": "all notifications marked as read",
	})
}

func parseNotificationPathUint(raw string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
