package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type PassengerHandler struct {
	passengerService *service.PassengerService
}

func NewPassengerHandler(passengerService *service.PassengerService) *PassengerHandler {
	return &PassengerHandler{passengerService: passengerService}
}

func (h *PassengerHandler) Create(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.Passenger
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.passengerService.Create(r.Context(), currentUser.ID, currentUser.Role, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, item)
}

func (h *PassengerHandler) List(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.passengerService.List(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, items)
}

func (h *PassengerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(r.PathValue("passengerId")), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid passengerId")
		return
	}
	if err := h.passengerService.Delete(r.Context(), currentUser.ID, currentUser.Role, uint(id)); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, map[string]bool{"deleted": true})
}

type PriceAlertHandler struct {
	priceAlertService *service.PriceAlertService
}

func NewPriceAlertHandler(priceAlertService *service.PriceAlertService) *PriceAlertHandler {
	return &PriceAlertHandler{priceAlertService: priceAlertService}
}

func (h *PriceAlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.PriceAlert
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.priceAlertService.Create(r.Context(), currentUser.ID, currentUser.Role, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, item)
}

func (h *PriceAlertHandler) List(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.priceAlertService.List(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, items)
}

func (h *PriceAlertHandler) Disable(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(r.PathValue("alertId")), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid alertId")
		return
	}
	if err := h.priceAlertService.Disable(r.Context(), currentUser.ID, currentUser.Role, uint(id)); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, map[string]bool{"disabled": true})
}

type DriverProfileHandler struct {
	driverProfileService *service.DriverProfileService
}

func NewDriverProfileHandler(driverProfileService *service.DriverProfileService) *DriverProfileHandler {
	return &DriverProfileHandler{driverProfileService: driverProfileService}
}

func (h *DriverProfileHandler) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.DriverProfile
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.driverProfileService.UpsertProfile(r.Context(), currentUser.ID, currentUser.Role, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *DriverProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	item, err := h.driverProfileService.GetProfile(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *DriverProfileHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.Vehicle
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.driverProfileService.CreateVehicle(r.Context(), currentUser.ID, currentUser.Role, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, item)
}

func (h *DriverProfileHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.driverProfileService.ListVehicles(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, items)
}
