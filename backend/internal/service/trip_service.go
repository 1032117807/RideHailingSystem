package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type TripStopInput struct {
	StopOrder         int
	StopName          string
	PlanArrivalTime   *time.Time
	PlanDepartureTime *time.Time
}

type CreateTripInput struct {
	VehicleType   string
	StartCity     string
	EndCity       string
	DepartureTime time.Time
	ArrivalTime   time.Time
	SeatTotal     int
	PriceCent     int
	Stops         []TripStopInput
}

type DriverTripOrderSummary struct {
	PendingVerificationCount int    `json:"pendingVerificationCount"`
	RefundRequestCount       int    `json:"refundRequestCount"`
	PendingVerificationNote  string `json:"pendingVerificationNote"`
	RefundRequestNote        string `json:"refundRequestNote"`
}

type DriverTripDetail struct {
	*model.Trip
	OrderSummary DriverTripOrderSummary `json:"orderSummary"`
	Orders       []*model.Order         `json:"orders"`
}

type DriverDashboardTripItem struct {
	ID              uint      `json:"id"`
	Route           string    `json:"route"`
	DepartureTime   time.Time `json:"departureTime"`
	SoldTickets     int       `json:"soldTickets"`
	SeatTotal       int       `json:"seatTotal"`
	EstimatedIncome int       `json:"estimatedIncome"`
}

type DriverDashboardSummary struct {
	TodayTripCount           int                       `json:"todayTripCount"`
	CompletedTodayTripCount  int                       `json:"completedTodayTripCount"`
	SoldTicketCount          int                       `json:"soldTicketCount"`
	SeatOccupancyRate        float64                   `json:"seatOccupancyRate"`
	TodayIncome              int                       `json:"todayIncome"`
	PendingVerificationCount int                       `json:"pendingVerificationCount"`
	RefundRequestCount       int                       `json:"refundRequestCount"`
	UpcomingTrips            []DriverDashboardTripItem `json:"upcomingTrips"`
	Alerts                   []string                  `json:"alerts"`
}

type DriverIncomeRouteItem struct {
	Route         string  `json:"route"`
	Income        int     `json:"income"`
	TicketCount   int     `json:"ticketCount"`
	OccupancyRate float64 `json:"occupancyRate"`
}

type DriverIncomeSummary struct {
	TodayIncome         int                     `json:"todayIncome"`
	WeeklyIncome        int                     `json:"weeklyIncome"`
	AvgOrderAmount      int                     `json:"avgOrderAmount"`
	RefundRate          float64                 `json:"refundRate"`
	PendingSettleAmount int                     `json:"pendingSettleAmount"`
	TopRoutes           []DriverIncomeRouteItem `json:"topRoutes"`
	Suggestions         []string                `json:"suggestions"`
}

type TripService struct {
	tripRepo  repository.TripRepository
	orderRepo repository.OrderRepository
}

func NewTripService(tripRepo repository.TripRepository, orderRepo repository.OrderRepository) *TripService {
	return &TripService{
		tripRepo:  tripRepo,
		orderRepo: orderRepo,
	}
}

func (s *TripService) CreateTrip(ctx context.Context, currentUserID uint, currentUserRole string, input CreateTripInput) (*model.Trip, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can create trip")
	}

	input.VehicleType = strings.TrimSpace(input.VehicleType)
	input.StartCity = strings.TrimSpace(input.StartCity)
	input.EndCity = strings.TrimSpace(input.EndCity)

	if input.VehicleType == "" {
		input.VehicleType = model.VehicleTypeCar
	}
	if input.StartCity == "" || input.EndCity == "" {
		return nil, errors.New("startCity and endCity are required")
	}
	if input.StartCity == input.EndCity {
		return nil, errors.New("startCity and endCity cannot be the same")
	}
	if !input.ArrivalTime.After(input.DepartureTime) {
		return nil, errors.New("arrivalTime must be later than departureTime")
	}
	if input.SeatTotal <= 0 {
		return nil, errors.New("seatTotal must be greater than 0")
	}
	if input.PriceCent <= 0 {
		return nil, errors.New("priceCent must be greater than 0")
	}

	stops := make([]model.TripStop, 0, len(input.Stops))
	seenOrders := make(map[int]struct{}, len(input.Stops))
	seenStopNames := make(map[string]struct{}, len(input.Stops))
	for _, stop := range input.Stops {
		stop.StopName = strings.TrimSpace(stop.StopName)
		if stop.StopOrder <= 0 {
			return nil, errors.New("stopOrder must be greater than 0")
		}
		if stop.StopName == "" {
			return nil, errors.New("stopName is required")
		}
		if stop.StopName == input.StartCity || stop.StopName == input.EndCity {
			return nil, errors.New("stopName cannot be the same as startCity or endCity")
		}
		if _, exists := seenOrders[stop.StopOrder]; exists {
			return nil, errors.New("stopOrder cannot repeat")
		}
		if _, exists := seenStopNames[stop.StopName]; exists {
			return nil, errors.New("stopName cannot repeat")
		}
		seenOrders[stop.StopOrder] = struct{}{}
		seenStopNames[stop.StopName] = struct{}{}

		stops = append(stops, model.TripStop{
			StopOrder:         stop.StopOrder,
			StopName:          stop.StopName,
			PlanArrivalTime:   stop.PlanArrivalTime,
			PlanDepartureTime: stop.PlanDepartureTime,
		})
	}

	sort.Slice(stops, func(i, j int) bool {
		return stops[i].StopOrder < stops[j].StopOrder
	})

	trip := &model.Trip{
		DriverID:      currentUserID,
		VehicleType:   input.VehicleType,
		StartCity:     input.StartCity,
		EndCity:       input.EndCity,
		DepartureTime: input.DepartureTime,
		ArrivalTime:   input.ArrivalTime,
		SeatTotal:     input.SeatTotal,
		SeatAvailable: input.SeatTotal,
		PriceCent:     input.PriceCent,
		Status:        model.TripStatusPublished,
	}

	if err := s.tripRepo.CreateWithStops(ctx, trip, stops); err != nil {
		return nil, err
	}
	return trip, nil
}

func (s *TripService) ListDriverTrips(ctx context.Context, currentUserID uint, currentUserRole string) ([]*model.Trip, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can view driver trips")
	}
	return s.tripRepo.ListByDriverID(ctx, currentUserID)
}

func (s *TripService) GetDriverTripDetail(ctx context.Context, currentUserID uint, currentUserRole string, tripID uint) (*DriverTripDetail, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can view driver trip detail")
	}

	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, errors.New("trip not found")
	}
	if trip.DriverID != currentUserID {
		return nil, errors.New("trip does not belong to current driver")
	}

	pendingVerificationCount := 0
	refundRequestCount := 0
	var orders []*model.Order
	if s.orderRepo != nil {
		pendingVerificationCount, refundRequestCount, err = s.orderRepo.CountSummaryByTripID(ctx, tripID)
		if err != nil {
			return nil, err
		}
		orders, err = s.orderRepo.ListByTripID(ctx, tripID)
		if err != nil {
			return nil, err
		}
	}

	return &DriverTripDetail{
		Trip:   trip,
		Orders: orders,
		OrderSummary: DriverTripOrderSummary{
			PendingVerificationCount: pendingVerificationCount,
			RefundRequestCount:       refundRequestCount,
			PendingVerificationNote:  "建议发车前 20 分钟开启检票页面，提高现场效率。",
			RefundRequestNote:        "当前暂无退款申请。",
		},
	}, nil
}

func (s *TripService) GetDriverDashboard(ctx context.Context, currentUserID uint, currentUserRole string) (*DriverDashboardSummary, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can view driver dashboard")
	}

	trips, err := s.tripRepo.ListByDriverID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	upcomingEnd := now.Add(4 * time.Hour)

	var todayTripCount int
	var completedTodayTripCount int
	var soldTicketCount int
	var totalSeats int
	var todayIncome int
	var pendingVerificationCount int
	var refundRequestCount int

	upcomingTrips := make([]DriverDashboardTripItem, 0)

	for _, trip := range trips {
		if trip == nil {
			continue
		}

		isTodayTrip := !trip.DepartureTime.Before(dayStart) && trip.DepartureTime.Before(dayEnd)
		if isTodayTrip {
			todayTripCount++
			if trip.ArrivalTime.Before(now) {
				completedTodayTripCount++
			}
			totalSeats += trip.SeatTotal
		}

		orders, err := s.orderRepo.ListByTripID(ctx, trip.ID)
		if err != nil {
			return nil, err
		}

		var tripSoldTickets int
		var tripEstimatedIncome int

		for _, order := range orders {
			if order == nil {
				continue
			}

			if order.OrderStatus != model.OrderStatusCancelled {
				tripSoldTickets += order.TicketCount
			}
			if order.PayStatus == model.PayStatusPaid {
				tripEstimatedIncome += order.Amount
				if isTodayTrip {
					todayIncome += order.Amount
				}
			}
			if order.OrderStatus == model.OrderStatusPendingVerification {
				pendingVerificationCount += order.TicketCount
			}
			if order.RefundStatus == model.RefundStatusRequested {
				refundRequestCount++
			}
		}

		if isTodayTrip {
			soldTicketCount += tripSoldTickets
		}

		if trip.DepartureTime.After(now) && trip.DepartureTime.Before(upcomingEnd) {
			upcomingTrips = append(upcomingTrips, DriverDashboardTripItem{
				ID:              trip.ID,
				Route:           trip.StartCity + " -> " + trip.EndCity,
				DepartureTime:   trip.DepartureTime,
				SoldTickets:     tripSoldTickets,
				SeatTotal:       trip.SeatTotal,
				EstimatedIncome: tripEstimatedIncome,
			})
		}
	}

	sort.Slice(upcomingTrips, func(i, j int) bool {
		return upcomingTrips[i].DepartureTime.Before(upcomingTrips[j].DepartureTime)
	})
	if len(upcomingTrips) > 3 {
		upcomingTrips = upcomingTrips[:3]
	}

	var seatOccupancyRate float64
	if totalSeats > 0 {
		seatOccupancyRate = float64(soldTicketCount) / float64(totalSeats)
	}

	alerts := make([]string, 0, 3)
	if pendingVerificationCount > 0 {
		alerts = append(alerts, "有待核销订单，请在发车前及时处理。")
	}
	if refundRequestCount > 0 {
		alerts = append(alerts, "有退款申请待关注，请尽快查看对应行程订单。")
	}
	if len(alerts) == 0 {
		alerts = append(alerts, "当前没有待处理异常，司机端运营状态正常。")
	}

	return &DriverDashboardSummary{
		TodayTripCount:           todayTripCount,
		CompletedTodayTripCount:  completedTodayTripCount,
		SoldTicketCount:          soldTicketCount,
		SeatOccupancyRate:        seatOccupancyRate,
		TodayIncome:              todayIncome,
		PendingVerificationCount: pendingVerificationCount,
		RefundRequestCount:       refundRequestCount,
		UpcomingTrips:            upcomingTrips,
		Alerts:                   alerts,
	}, nil
}

func (s *TripService) GetDriverIncome(ctx context.Context, currentUserID uint, currentUserRole string) (*DriverIncomeSummary, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can view driver income")
	}

	trips, err := s.tripRepo.ListByDriverID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := dayStart.AddDate(0, 0, -6)

	var todayIncome int
	var weeklyIncome int
	var paidOrderCount int
	var paidOrderAmount int
	var refundOrderCount int
	var pendingSettleAmount int

	type routeAgg struct {
		Income      int
		TicketCount int
		SeatTotal   int
	}
	routeMap := make(map[string]*routeAgg)

	for _, trip := range trips {
		if trip == nil {
			continue
		}

		route := trip.StartCity + " -> " + trip.EndCity
		if _, ok := routeMap[route]; !ok {
			routeMap[route] = &routeAgg{}
		}
		routeMap[route].SeatTotal += trip.SeatTotal

		orders, err := s.orderRepo.ListByTripID(ctx, trip.ID)
		if err != nil {
			return nil, err
		}

		for _, order := range orders {
			if order == nil {
				continue
			}

			if order.PayStatus == model.PayStatusPaid {
				paidOrderCount++
				paidOrderAmount += order.Amount
				routeMap[route].Income += order.Amount
				routeMap[route].TicketCount += order.TicketCount

				if !trip.DepartureTime.Before(dayStart) {
					todayIncome += order.Amount
				}
				if !trip.DepartureTime.Before(weekStart) {
					weeklyIncome += order.Amount
				}
				if trip.ArrivalTime.After(now) {
					pendingSettleAmount += order.Amount
				}
			}

			if order.RefundStatus == model.RefundStatusRequested ||
				order.RefundStatus == model.RefundStatusRefunded ||
				order.RefundStatus == model.RefundStatusRejected {
				refundOrderCount++
			}
		}
	}

	topRoutes := make([]DriverIncomeRouteItem, 0, len(routeMap))
	for route, agg := range routeMap {
		var occupancyRate float64
		if agg.SeatTotal > 0 {
			occupancyRate = float64(agg.TicketCount) / float64(agg.SeatTotal)
		}

		topRoutes = append(topRoutes, DriverIncomeRouteItem{
			Route:         route,
			Income:        agg.Income,
			TicketCount:   agg.TicketCount,
			OccupancyRate: occupancyRate,
		})
	}

	sort.Slice(topRoutes, func(i, j int) bool {
		return topRoutes[i].Income > topRoutes[j].Income
	})
	if len(topRoutes) > 3 {
		topRoutes = topRoutes[:3]
	}

	var avgOrderAmount int
	if paidOrderCount > 0 {
		avgOrderAmount = paidOrderAmount / paidOrderCount
	}

	var refundRate float64
	if paidOrderCount > 0 {
		refundRate = float64(refundOrderCount) / float64(paidOrderCount)
	}

	suggestions := []string{
		"优先优化高收入路线的发车时段，尽量提升满座率。",
		"发车前集中处理待核销订单，可减少现场排队时间。",
	}
	if refundRate > 0.1 {
		suggestions = append(suggestions, "当前退款率偏高，建议检查票价和出发时间设置。")
	}

	return &DriverIncomeSummary{
		TodayIncome:         todayIncome,
		WeeklyIncome:        weeklyIncome,
		AvgOrderAmount:      avgOrderAmount,
		RefundRate:          refundRate,
		PendingSettleAmount: pendingSettleAmount,
		TopRoutes:           topRoutes,
		Suggestions:         suggestions,
	}, nil
}
