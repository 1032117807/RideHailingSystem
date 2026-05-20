package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) Search(w http.ResponseWriter, r *http.Request) {
	startCity := strings.TrimSpace(r.URL.Query().Get("startCity"))
	endCity := strings.TrimSpace(r.URL.Query().Get("endCity"))
	dateRaw := strings.TrimSpace(r.URL.Query().Get("date"))
	dateFromRaw := strings.TrimSpace(r.URL.Query().Get("dateFrom"))
	dateToRaw := strings.TrimSpace(r.URL.Query().Get("dateTo"))
	allowTransfer := true
	if _, exists := r.URL.Query()["allowTransfer"]; exists {
		allowTransfer = parseBooleanQuery(r.URL.Query().Get("allowTransfer"))
	}

	date, dateFrom, dateTo, err := parseTicketSearchDates(dateRaw, dateFromRaw, dateToRaw)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	trips, err := h.ticketService.SearchTickets(r.Context(), service.SearchTicketsInput{
		StartCity:     startCity,
		EndCity:       endCity,
		Date:          date,
		DateFrom:      dateFrom,
		DateTo:        dateTo,
		AllowTransfer: allowTransfer,
		MinSeat:       parseIntQuery(r.URL.Query().Get("minSeat")),
		MaxPriceCent:  parseIntQuery(r.URL.Query().Get("maxPriceCent")),
		VehicleType:   strings.TrimSpace(r.URL.Query().Get("vehicleType")),
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, trips)
}

func (h *TicketHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ticketID, err := parsePathUint(r.PathValue("ticketId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ticketId")
		return
	}

	trip, err := h.ticketService.GetTicketDetail(r.Context(), ticketID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, trip)
}

func parseTicketSearchDates(dateRaw, dateFromRaw, dateToRaw string) (time.Time, time.Time, time.Time, error) {
	if dateFromRaw != "" || dateToRaw != "" {
		var dateFrom time.Time
		var dateTo time.Time
		var err error
		if dateFromRaw != "" {
			dateFrom, err = time.ParseInLocation("2006-01-02", dateFromRaw, time.Local)
			if err != nil {
				return time.Time{}, time.Time{}, time.Time{}, errors.New("dateFrom must be format 2006-01-02")
			}
		}
		if dateToRaw != "" {
			dateTo, err = time.ParseInLocation("2006-01-02", dateToRaw, time.Local)
			if err != nil {
				return time.Time{}, time.Time{}, time.Time{}, errors.New("dateTo must be format 2006-01-02")
			}
		}
		return time.Time{}, dateFrom, dateTo, nil
	}
	date, err := time.ParseInLocation("2006-01-02", dateRaw, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, errors.New("date must be format 2006-01-02")
	}
	return date, time.Time{}, time.Time{}, nil
}

func parseIntQuery(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func parseBooleanQuery(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
