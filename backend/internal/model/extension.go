package model

import "time"

const (
	PriceAlertStatusActive    = "active"
	PriceAlertStatusTriggered = "triggered"
	PriceAlertStatusDisabled  = "disabled"

	DriverProfileStatusPending  = "pending"
	DriverProfileStatusApproved = "approved"
	DriverProfileStatusRejected = "rejected"
)

type Passenger struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"userId"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	IDCard    string    `gorm:"column:id_card;size:32;not null" json:"idCard"`
	Phone     string    `gorm:"size:20;not null" json:"phone"`
	IsDefault bool      `gorm:"column:is_default;not null;default:false" json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Passenger) TableName() string {
	return "passengers"
}

type PriceAlert struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"column:user_id;not null;index" json:"userId"`
	StartCity       string     `gorm:"column:start_city;size:50;not null;index" json:"startCity"`
	EndCity         string     `gorm:"column:end_city;size:50;not null;index" json:"endCity"`
	TargetPriceCent int        `gorm:"column:target_price_cent;not null" json:"targetPriceCent"`
	StartDate       time.Time  `gorm:"column:start_date;type:date;not null" json:"startDate"`
	EndDate         time.Time  `gorm:"column:end_date;type:date;not null" json:"endDate"`
	Status          string     `gorm:"size:20;not null;default:active;index" json:"status"`
	TriggeredTripID *uint      `gorm:"column:triggered_trip_id" json:"triggeredTripId,omitempty"`
	TriggeredAt     *time.Time `gorm:"column:triggered_at" json:"triggeredAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (PriceAlert) TableName() string {
	return "price_alerts"
}

type DriverProfile struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"column:user_id;not null;uniqueIndex" json:"userId"`
	RealName        string     `gorm:"column:real_name;size:50;not null" json:"realName"`
	IDCard          string     `gorm:"column:id_card;size:32;not null" json:"idCard"`
	LicenseNo       string     `gorm:"column:license_no;size:64;not null" json:"licenseNo"`
	IDCardImageURL  string     `gorm:"column:id_card_image_url;size:255" json:"idCardImageUrl"`
	LicenseImageURL string     `gorm:"column:license_image_url;size:255" json:"licenseImageUrl"`
	Status          string     `gorm:"size:20;not null;default:pending;index" json:"status"`
	ReviewNote      string     `gorm:"column:review_note;size:255" json:"reviewNote"`
	ReviewedAt      *time.Time `gorm:"column:reviewed_at" json:"reviewedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (DriverProfile) TableName() string {
	return "driver_profiles"
}

type Vehicle struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DriverID  uint      `gorm:"column:driver_id;not null;index" json:"driverId"`
	PlateNo   string    `gorm:"column:plate_no;size:32;not null" json:"plateNo"`
	Brand     string    `gorm:"size:64" json:"brand"`
	ModelName string    `gorm:"column:model_name;size:64" json:"modelName"`
	SeatCount int       `gorm:"column:seat_count;not null" json:"seatCount"`
	Status    string    `gorm:"size:20;not null;default:active;index" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

type DriverSettlement struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	DriverID         uint      `gorm:"column:driver_id;not null;index" json:"driverId"`
	SettlementDate   time.Time `gorm:"column:settlement_date;type:date;not null;index" json:"settlementDate"`
	GrossAmountCent  int       `gorm:"column:gross_amount_cent;not null" json:"grossAmountCent"`
	RefundAmountCent int       `gorm:"column:refund_amount_cent;not null" json:"refundAmountCent"`
	ServiceFeeCent   int       `gorm:"column:service_fee_cent;not null" json:"serviceFeeCent"`
	NetAmountCent    int       `gorm:"column:net_amount_cent;not null" json:"netAmountCent"`
	Status           string    `gorm:"size:20;not null;default:pending;index" json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (DriverSettlement) TableName() string {
	return "driver_settlements"
}
