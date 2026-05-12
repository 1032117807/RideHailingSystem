package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type createPaymentRequest struct {
	OrderID uint   `json:"orderId"`
	Channel string `json:"channel"`
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createPaymentRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := h.paymentService.CreatePayment(r.Context(), currentUser.ID, currentUser.Role, service.CreatePaymentInput{
		OrderID: req.OrderID,
		Channel: req.Channel,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, payment)
}

func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	paymentID, err := parsePaymentPathUint(r.PathValue("paymentId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid paymentId")
		return
	}

	payment, err := h.paymentService.GetPaymentStatus(r.Context(), currentUser.ID, currentUser.Role, paymentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, payment)
}

func (h *PaymentHandler) MockPaySuccess(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	paymentID, err := parsePaymentPathUint(r.PathValue("paymentId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid paymentId")
		return
	}

	payment, err := h.paymentService.MockPaySuccess(r.Context(), currentUser.ID, currentUser.Role, paymentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, payment)
}

func parsePaymentPathUint(raw string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
