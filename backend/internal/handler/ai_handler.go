package handler

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/repository"
	"ridehailing/backend/internal/service"
)

type AIHandler struct {
	aiService          *service.AIService
	aiLimiter          repository.AIRateLimiter
	passengerChatLimit int64
	driverDraftLimit   int64
	rateLimitWindow    time.Duration
	riskService        *service.RiskService
}

func NewAIHandler(
	aiService *service.AIService,
	aiLimiter repository.AIRateLimiter,
	passengerChatLimit int,
	driverDraftLimit int,
	rateLimitWindow time.Duration,
	riskService *service.RiskService,
) *AIHandler {
	return &AIHandler{
		aiService:          aiService,
		aiLimiter:          aiLimiter,
		passengerChatLimit: int64(passengerChatLimit),
		driverDraftLimit:   int64(driverDraftLimit),
		rateLimitWindow:    rateLimitWindow,
		riskService:        riskService,
	}
}

type createDriverTripDraftRequest struct {
	Prompt string `json:"prompt"`
}

type passengerAIChatRequest struct {
	Messages []service.AIChatMessage `json:"messages"`
}

func (h *AIHandler) CreateDriverTripDraft(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	subject := fmt.Sprintf("user:%d", currentUser.ID)
	if blocked := h.rejectIfRateLimited(w, r, "driver-create-trip", subject, h.driverDraftLimit); blocked {
		if h.riskService != nil {
			_ = h.riskService.RecordAIRateLimitEvent(
				r.Context(),
				currentUser.ID,
				currentUser.Role,
				"driver-create-trip",
				subject,
				h.driverDraftLimit,
				h.rateLimitWindow,
			)
		}
		return
	}

	var req createDriverTripDraftRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	draft, err := h.aiService.GenerateDriverTripDraft(r.Context(), currentUser.ID, currentUser.Role, req.Prompt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, draft)
}

func (h *AIHandler) ChatPassenger(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := middleware.CurrentUser(r)

	subject := "ip:" + clientIP(r)
	var currentUserID uint
	var currentUserRole string
	if currentUser != nil {
		currentUserID = currentUser.ID
		currentUserRole = currentUser.Role
		subject = fmt.Sprintf("user:%d", currentUser.ID)
	}

	if blocked := h.rejectIfRateLimited(w, r, "passenger-chat", subject, h.passengerChatLimit); blocked {
		if h.riskService != nil && currentUser != nil {
			_ = h.riskService.RecordAIRateLimitEvent(
				r.Context(),
				currentUser.ID,
				currentUser.Role,
				"passenger-chat",
				subject,
				h.passengerChatLimit,
				h.rateLimitWindow,
			)
		}
		return
	}

	var req passengerAIChatRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.aiService.ChatPassenger(r.Context(), currentUserID, currentUserRole, req.Messages)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, result)
}

func (h *AIHandler) rejectIfRateLimited(w http.ResponseWriter, r *http.Request, scope string, subject string, limit int64) bool {
	if h.aiLimiter == nil || limit <= 0 || h.rateLimitWindow <= 0 {
		return false
	}

	allowed, remaining, retryAfter, err := h.aiLimiter.Allow(r.Context(), scope, subject, limit, h.rateLimitWindow)
	if err != nil {
		log.Printf("ai rate limit failed: scope=%s subject=%s err=%v", scope, subject, err)
		return false
	}

	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
	if retryAfter > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
	}

	if !allowed {
		response.Error(w, http.StatusTooManyRequests, "too many AI requests, please retry later")
		return true
	}
	return false
}

func clientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}

	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	if strings.TrimSpace(r.RemoteAddr) != "" {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return "unknown"
}
