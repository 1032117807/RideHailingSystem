package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

const (
	minTransferWait   = 20 * time.Minute
	maxTransferWait   = 6 * time.Hour
	maxTransferPlan   = 8
	maxSuggestedRoute = 6
	maxCityMatchTrip  = 20
)

type SearchTicketsInput struct {
	StartCity     string
	EndCity       string
	Date          time.Time
	AllowTransfer bool
}

type SearchCityTicketsInput struct {
	City string
	Date time.Time
	Role string
}

type TicketSearchLeg struct {
	TripID        uint      `json:"tripId"`
	StartCity     string    `json:"startCity"`
	EndCity       string    `json:"endCity"`
	DepartureTime time.Time `json:"departureTime"`
	ArrivalTime   time.Time `json:"arrivalTime"`
	SeatAvailable int       `json:"seatAvailable"`
	PriceCent     int       `json:"priceCent"`
	VehicleType   string    `json:"vehicleType"`
}

type TicketRouteSuggestion struct {
	Route            string `json:"route"`
	TransferCity     string `json:"transferCity"`
	FirstLegCount    int    `json:"firstLegCount"`
	SecondLegCount   int    `json:"secondLegCount"`
	TotalOptionCount int    `json:"totalOptionCount"`
	Reason           string `json:"reason"`
}

type TicketSearchResult struct {
	Kind               string                   `json:"kind"`
	ID                 uint                     `json:"id,omitempty"`
	StartCity          string                   `json:"startCity"`
	EndCity            string                   `json:"endCity"`
	DepartureTime      time.Time                `json:"departureTime"`
	ArrivalTime        time.Time                `json:"arrivalTime"`
	SeatAvailable      int                      `json:"seatAvailable"`
	PriceCent          int                      `json:"priceCent"`
	VehicleType        string                   `json:"vehicleType"`
	Status             string                   `json:"status"`
	TransferCity       string                   `json:"transferCity,omitempty"`
	TransferWaitMinute int                      `json:"transferWaitMinute,omitempty"`
	MatchedCity        string                   `json:"matchedCity,omitempty"`
	MatchRoles         []string                 `json:"matchRoles,omitempty"`
	Stops              []string                 `json:"stops,omitempty"`
	Legs               []*TicketSearchLeg       `json:"legs,omitempty"`
	Suggestions        []*TicketRouteSuggestion `json:"suggestions,omitempty"`
}

type TicketService struct {
	tripRepo repository.TripRepository
}

type tripPoint struct {
	City      string
	Order     int
	ArriveAt  time.Time
	DepartAt  time.Time
	PointRole string
}

type tripSegment struct {
	Trip  *model.Trip
	From  tripPoint
	To    tripPoint
	Route []tripPoint
}

func NewTicketService(tripRepo repository.TripRepository) *TicketService {
	return &TicketService{tripRepo: tripRepo}
}

func (s *TicketService) SearchTickets(ctx context.Context, input SearchTicketsInput) ([]*TicketSearchResult, error) {
	input.StartCity = strings.TrimSpace(input.StartCity)
	input.EndCity = strings.TrimSpace(input.EndCity)

	if input.StartCity == "" || input.EndCity == "" {
		return nil, errors.New("startCity and endCity are required")
	}
	if input.Date.IsZero() {
		return nil, errors.New("date is required")
	}
	if input.StartCity == input.EndCity {
		return nil, errors.New("startCity and endCity must be different")
	}

	dayStart, dayEnd := dayRange(input.Date)
	allTrips, err := s.tripRepo.ListPublishedTripsByDate(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	results := buildDirectTicketSearchResults(input.StartCity, input.EndCity, allTrips)
	if !input.AllowTransfer {
		sortTicketSearchResults(results)
		return results, nil
	}

	results = append(results, buildTransferTicketSearchResults(input.StartCity, input.EndCity, allTrips)...)
	if suggestions := buildFuzzyRouteSuggestions(input.StartCity, input.EndCity, allTrips); len(suggestions) > 0 {
		results = append(results, &TicketSearchResult{
			Kind:          "suggestion",
			StartCity:     input.StartCity,
			EndCity:       input.EndCity,
			DepartureTime: dayStart,
			ArrivalTime:   dayStart,
			Suggestions:   suggestions,
			Status:        model.TripStatusPublished,
		})
	}

	sortTicketSearchResults(results)
	return results, nil
}

func (s *TicketService) SearchCityTickets(ctx context.Context, input SearchCityTicketsInput) ([]*TicketSearchResult, error) {
	input.City = strings.TrimSpace(input.City)
	input.Role = normalizeCitySearchRole(input.Role)
	if input.City == "" {
		return nil, errors.New("city is required")
	}
	if input.Date.IsZero() {
		return nil, errors.New("date is required")
	}

	dayStart, dayEnd := dayRange(input.Date)
	trips, err := s.tripRepo.ListPublishedTripsByDate(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	results := make([]*TicketSearchResult, 0)
	for _, trip := range trips {
		roles := matchTripCityRoles(trip, input.City)
		if len(roles) == 0 || !cityRolesMatch(roles, input.Role) {
			continue
		}
		result := buildWholeTripSearchResult(trip, "city_match")
		result.MatchedCity = input.City
		result.MatchRoles = roles
		results = append(results, result)
	}

	sortTicketSearchResults(results)
	if len(results) > maxCityMatchTrip {
		results = results[:maxCityMatchTrip]
	}
	return results, nil
}

func (s *TicketService) GetTicketDetail(ctx context.Context, tripID uint) (*model.Trip, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	if trip == nil || trip.Status != model.TripStatusPublished {
		return nil, errors.New("ticket not found")
	}
	return trip, nil
}

func buildDirectTicketSearchResults(startCity, endCity string, trips []*model.Trip) []*TicketSearchResult {
	results := make([]*TicketSearchResult, 0)
	seen := make(map[string]struct{})
	for _, trip := range trips {
		for _, segment := range findTripSegments(trip, startCity, endCity) {
			key := buildSegmentResultKey(segment)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			results = append(results, buildSegmentTicketSearchResult("direct", segment))
		}
	}
	sortTicketSearchResults(results)
	return results
}

func buildWholeTripSearchResult(trip *model.Trip, kind string) *TicketSearchResult {
	result := &TicketSearchResult{
		Kind:          kind,
		ID:            trip.ID,
		StartCity:     trip.StartCity,
		EndCity:       trip.EndCity,
		DepartureTime: trip.DepartureTime,
		ArrivalTime:   trip.ArrivalTime,
		SeatAvailable: trip.SeatAvailable,
		PriceCent:     trip.PriceCent,
		VehicleType:   trip.VehicleType,
		Status:        trip.Status,
		Stops:         tripStopNames(trip),
		Legs: []*TicketSearchLeg{
			buildTicketSearchLeg(trip, trip.StartCity, trip.EndCity, trip.DepartureTime, trip.ArrivalTime),
		},
	}
	return result
}

func buildSegmentTicketSearchResult(kind string, segment tripSegment) *TicketSearchResult {
	trip := segment.Trip
	return &TicketSearchResult{
		Kind:          kind,
		ID:            trip.ID,
		StartCity:     segment.From.City,
		EndCity:       segment.To.City,
		DepartureTime: segment.From.DepartAt,
		ArrivalTime:   segment.To.ArriveAt,
		SeatAvailable: trip.SeatAvailable,
		PriceCent:     trip.PriceCent,
		VehicleType:   trip.VehicleType,
		Status:        trip.Status,
		Stops:         routePointNames(segment.Route),
		Legs: []*TicketSearchLeg{
			buildTicketSearchLeg(trip, segment.From.City, segment.To.City, segment.From.DepartAt, segment.To.ArriveAt),
		},
	}
}

func buildTicketSearchLeg(trip *model.Trip, startCity, endCity string, departureTime, arrivalTime time.Time) *TicketSearchLeg {
	return &TicketSearchLeg{
		TripID:        trip.ID,
		StartCity:     startCity,
		EndCity:       endCity,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		SeatAvailable: trip.SeatAvailable,
		PriceCent:     trip.PriceCent,
		VehicleType:   trip.VehicleType,
	}
}

func buildTransferTicketSearchResults(startCity, endCity string, trips []*model.Trip) []*TicketSearchResult {
	results := make([]*TicketSearchResult, 0)
	seen := make(map[string]struct{})

	for _, firstTrip := range trips {
		for _, transferCity := range tripReachableCitiesAfter(firstTrip, startCity) {
			if transferCity == startCity || transferCity == endCity {
				continue
			}
			firstSegments := findTripSegments(firstTrip, startCity, transferCity)
			if len(firstSegments) == 0 {
				continue
			}
			for _, secondTrip := range trips {
				secondSegments := findTripSegments(secondTrip, transferCity, endCity)
				for _, first := range firstSegments {
					for _, second := range secondSegments {
						if first.Trip.ID == second.Trip.ID {
							continue
						}
						wait := second.From.DepartAt.Sub(first.To.ArriveAt)
						if wait < minTransferWait || wait > maxTransferWait {
							continue
						}
						key := buildTransferSegmentResultKey(first, second)
						if _, ok := seen[key]; ok {
							continue
						}
						seen[key] = struct{}{}
						results = append(results, buildTransferTicketSearchResult(startCity, endCity, transferCity, wait, first, second))
					}
				}
			}
		}
	}

	sortTicketSearchResults(results)
	if len(results) > maxTransferPlan {
		results = results[:maxTransferPlan]
	}
	return results
}

func buildTransferTicketSearchResult(startCity, endCity, transferCity string, wait time.Duration, first, second tripSegment) *TicketSearchResult {
	return &TicketSearchResult{
		Kind:               "transfer",
		StartCity:          startCity,
		EndCity:            endCity,
		DepartureTime:      first.From.DepartAt,
		ArrivalTime:        second.To.ArriveAt,
		SeatAvailable:      minInt(first.Trip.SeatAvailable, second.Trip.SeatAvailable),
		PriceCent:          first.Trip.PriceCent + second.Trip.PriceCent,
		VehicleType:        first.Trip.VehicleType + " + " + second.Trip.VehicleType,
		Status:             model.TripStatusPublished,
		TransferCity:       transferCity,
		TransferWaitMinute: int(wait / time.Minute),
		Legs: []*TicketSearchLeg{
			buildTicketSearchLeg(first.Trip, first.From.City, first.To.City, first.From.DepartAt, first.To.ArriveAt),
			buildTicketSearchLeg(second.Trip, second.From.City, second.To.City, second.From.DepartAt, second.To.ArriveAt),
		},
	}
}

func buildFuzzyRouteSuggestions(startCity, endCity string, trips []*model.Trip) []*TicketRouteSuggestion {
	outgoingCounts := make(map[string]int)
	incomingCounts := make(map[string]int)

	for _, trip := range trips {
		for _, city := range tripReachableCitiesAfter(trip, startCity) {
			if city != "" && city != startCity && city != endCity {
				outgoingCounts[city]++
			}
		}
		for _, city := range tripReachableCitiesBefore(trip, endCity) {
			if city != "" && city != startCity && city != endCity {
				incomingCounts[city]++
			}
		}
	}

	suggestions := make([]*TicketRouteSuggestion, 0)
	for city, firstLegCount := range outgoingCounts {
		secondLegCount := incomingCounts[city]
		if secondLegCount == 0 {
			continue
		}
		suggestions = append(suggestions, &TicketRouteSuggestion{
			Route:            startCity + " -> " + city + " -> " + endCity,
			TransferCity:     city,
			FirstLegCount:    firstLegCount,
			SecondLegCount:   secondLegCount,
			TotalOptionCount: firstLegCount + secondLegCount,
			Reason:           fmt.Sprintf("当天可尝试经%s中转：%d班可从%s到%s，%d班可从%s到%s。", city, firstLegCount, startCity, city, secondLegCount, city, endCity),
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		left := suggestions[i]
		right := suggestions[j]
		if left.TotalOptionCount == right.TotalOptionCount {
			if left.FirstLegCount == right.FirstLegCount {
				return left.TransferCity < right.TransferCity
			}
			return left.FirstLegCount > right.FirstLegCount
		}
		return left.TotalOptionCount > right.TotalOptionCount
	})

	if len(suggestions) > maxSuggestedRoute {
		suggestions = suggestions[:maxSuggestedRoute]
	}
	return suggestions
}

func findTripSegments(trip *model.Trip, startCity, endCity string) []tripSegment {
	points := tripRoutePoints(trip)
	results := make([]tripSegment, 0)
	for i, from := range points {
		if from.City != startCity {
			continue
		}
		for j := i + 1; j < len(points); j++ {
			to := points[j]
			if to.City == endCity && to.ArriveAt.After(from.DepartAt) {
				results = append(results, tripSegment{
					Trip:  trip,
					From:  from,
					To:    to,
					Route: points[i : j+1],
				})
			}
		}
	}
	return results
}

func tripRoutePoints(trip *model.Trip) []tripPoint {
	if trip == nil {
		return nil
	}
	points := []tripPoint{
		{
			City:      strings.TrimSpace(trip.StartCity),
			Order:     0,
			ArriveAt:  trip.DepartureTime,
			DepartAt:  trip.DepartureTime,
			PointRole: "start",
		},
	}
	stops := append([]model.TripStop(nil), trip.Stops...)
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].StopOrder < stops[j].StopOrder
	})
	for _, stop := range stops {
		city := strings.TrimSpace(stop.StopName)
		if city == "" || city == trip.StartCity || city == trip.EndCity {
			continue
		}
		arriveAt := trip.DepartureTime
		if stop.PlanArrivalTime != nil {
			arriveAt = *stop.PlanArrivalTime
		} else if stop.PlanDepartureTime != nil {
			arriveAt = *stop.PlanDepartureTime
		}
		departAt := arriveAt
		if stop.PlanDepartureTime != nil {
			departAt = *stop.PlanDepartureTime
		}
		points = append(points, tripPoint{
			City:      city,
			Order:     stop.StopOrder,
			ArriveAt:  arriveAt,
			DepartAt:  departAt,
			PointRole: "stop",
		})
	}
	points = append(points, tripPoint{
		City:      strings.TrimSpace(trip.EndCity),
		Order:     maxStopOrder(points) + 1,
		ArriveAt:  trip.ArrivalTime,
		DepartAt:  trip.ArrivalTime,
		PointRole: "end",
	})
	sort.Slice(points, func(i, j int) bool {
		return points[i].Order < points[j].Order
	})
	return points
}

func tripReachableCitiesAfter(trip *model.Trip, city string) []string {
	points := tripRoutePoints(trip)
	result := make([]string, 0)
	seen := make(map[string]struct{})
	for i, point := range points {
		if point.City != city {
			continue
		}
		for j := i + 1; j < len(points); j++ {
			next := points[j].City
			if next == "" {
				continue
			}
			if _, ok := seen[next]; ok {
				continue
			}
			seen[next] = struct{}{}
			result = append(result, next)
		}
	}
	return result
}

func tripReachableCitiesBefore(trip *model.Trip, city string) []string {
	points := tripRoutePoints(trip)
	result := make([]string, 0)
	seen := make(map[string]struct{})
	for i, point := range points {
		if point.City != city {
			continue
		}
		for j := 0; j < i; j++ {
			prev := points[j].City
			if prev == "" {
				continue
			}
			if _, ok := seen[prev]; ok {
				continue
			}
			seen[prev] = struct{}{}
			result = append(result, prev)
		}
	}
	return result
}

func matchTripCityRoles(trip *model.Trip, city string) []string {
	roles := make([]string, 0)
	seen := make(map[string]struct{})
	for _, point := range tripRoutePoints(trip) {
		if point.City != city {
			continue
		}
		role := point.PointRole
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		roles = append(roles, role)
	}
	return roles
}

func cityRolesMatch(roles []string, expected string) bool {
	if expected == "" || expected == "any" {
		return true
	}
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

func normalizeCitySearchRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "start", "origin", "起点", "出发":
		return "start"
	case "end", "destination", "终点", "到达":
		return "end"
	case "stop", "pass", "via", "经过", "途经", "中间站":
		return "stop"
	default:
		return "any"
	}
}

func routePointNames(points []tripPoint) []string {
	names := make([]string, 0, len(points))
	for _, point := range points {
		if point.City != "" {
			names = append(names, point.City)
		}
	}
	return cleanStringList(names)
}

func tripStopNames(trip *model.Trip) []string {
	return routePointNames(tripRoutePoints(trip))
}

func maxStopOrder(points []tripPoint) int {
	maxOrder := 0
	for _, point := range points {
		if point.Order > maxOrder {
			maxOrder = point.Order
		}
	}
	return maxOrder
}

func sortTicketSearchResults(results []*TicketSearchResult) {
	sort.Slice(results, func(i, j int) bool {
		left := results[i]
		right := results[j]
		if left.Kind != right.Kind {
			return ticketResultKindRank(left.Kind) < ticketResultKindRank(right.Kind)
		}
		if left.DepartureTime.Equal(right.DepartureTime) {
			if left.ArrivalTime.Equal(right.ArrivalTime) {
				if left.PriceCent == right.PriceCent {
					return left.ID < right.ID
				}
				return left.PriceCent < right.PriceCent
			}
			return left.ArrivalTime.Before(right.ArrivalTime)
		}
		return left.DepartureTime.Before(right.DepartureTime)
	})
}

func ticketResultKindRank(kind string) int {
	switch kind {
	case "direct":
		return 0
	case "transfer":
		return 1
	case "suggestion":
		return 2
	case "city_match":
		return 3
	default:
		return 4
	}
}

func buildSegmentResultKey(segment tripSegment) string {
	return strings.Join([]string{
		fmt.Sprint(segment.Trip.ID),
		segment.From.City,
		segment.To.City,
		segment.From.DepartAt.Format(time.RFC3339),
		segment.To.ArriveAt.Format(time.RFC3339),
	}, "|")
}

func buildTransferSegmentResultKey(first, second tripSegment) string {
	return strings.Join([]string{
		buildSegmentResultKey(first),
		buildSegmentResultKey(second),
	}, "||")
}

func dayRange(date time.Time) (time.Time, time.Time) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return dayStart, dayStart.Add(24 * time.Hour)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
