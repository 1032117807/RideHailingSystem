package model

import "time"

type TripStop struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	TripID            uint       `gorm:"not null;uniqueIndex:idx_trip_stop_order,priority:1" json:"tripId"`
	StopOrder         int        `gorm:"not null;uniqueIndex:idx_trip_stop_order,priority:2" json:"stopOrder"`
	StopName          string     `gorm:"size:50;not null" json:"stopName"`
	PlanArrivalTime   *time.Time `json:"planArrivalTime,omitempty"`
	PlanDepartureTime *time.Time `json:"planDepartureTime,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func (TripStop) TableName() string {
	return "trip_stops"
}
