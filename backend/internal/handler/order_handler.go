package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type createOrderRequest struct {
	TripID      uint   `json:"tripId"`
	TicketCount int    `json:"ticketCount"`
	SeatType    string `json:"seatType"`
}

type adminRefundReviewRequest struct {
	ReviewNote string `json:"reviewNote"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createOrderRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), currentUser.ID, currentUser.Role, service.CreateOrderInput{
		TripID:      req.TripID,
		TicketCount: req.TicketCount,
		SeatType:    req.SeatType,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, order)
}

func (h *OrderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := h.orderService.ListMyOrders(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, orders)
}

func (h *OrderHandler) GetMyOrderDetail(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	order, err := h.orderService.GetMyOrderDetail(r.Context(), currentUser.ID, currentUser.Role, orderID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) CancelMyOrder(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	order, err := h.orderService.CancelMyOrder(r.Context(), currentUser.ID, currentUser.Role, orderID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) RequestRefund(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	order, err := h.orderService.RequestRefund(r.Context(), currentUser.ID, currentUser.Role, orderID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) VerifyOrderByDriver(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	order, err := h.orderService.VerifyOrderByDriver(r.Context(), currentUser.ID, currentUser.Role, orderID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) ListRefundOrdersForAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	refundStatus := strings.TrimSpace(r.URL.Query().Get("refundStatus"))
	orders, err := h.orderService.ListRefundOrdersForAdmin(r.Context(), currentUser.ID, currentUser.Role, refundStatus)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, orders)
}

func (h *OrderHandler) GetAdminDashboardSummary(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.orderService.GetAdminDashboardSummary(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, summary)
}

func (h *OrderHandler) ApproveRefundByAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	var req adminRefundReviewRequest
	if err := decodeJSONBody(r.Body, &req); err != nil && !strings.Contains(err.Error(), "EOF") {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.ApproveRefundByAdmin(r.Context(), currentUser.ID, currentUser.Role, orderID, req.ReviewNote)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) RejectRefundByAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, err := parseOrderPathUint(r.PathValue("orderId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid orderId")
		return
	}

	var req adminRefundReviewRequest
	if err := decodeJSONBody(r.Body, &req); err != nil && !strings.Contains(err.Error(), "EOF") {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.RejectRefundByAdmin(r.Context(), currentUser.ID, currentUser.Role, orderID, req.ReviewNote)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, order)
}

func parseOrderPathUint(raw string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
