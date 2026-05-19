package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type ElectronicTicketHandler struct {
	ticketService *service.ElectronicTicketService
}

func NewElectronicTicketHandler(ticketService *service.ElectronicTicketService) *ElectronicTicketHandler {
	return &ElectronicTicketHandler{ticketService: ticketService}
}

type verifyElectronicTicketRequest struct {
	Token string `json:"token"`
}

func (h *ElectronicTicketHandler) GetMyOrderTicket(w http.ResponseWriter, r *http.Request) {
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

	ticket, err := h.ticketService.GetMyOrderTicket(r.Context(), currentUser.ID, currentUser.Role, uint(orderID))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, ticket)
}

func (h *ElectronicTicketHandler) VerifyByDriver(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req verifyElectronicTicketRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ticket, err := h.ticketService.VerifyByDriver(r.Context(), currentUser.ID, currentUser.Role, req.Token)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, ticket)
}
