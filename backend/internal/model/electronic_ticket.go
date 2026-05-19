package model

import "time"

const (
	ElectronicTicketStatusIssued   = "issued"
	ElectronicTicketStatusVerified = "verified"
	ElectronicTicketStatusVoided   = "voided"

	TicketVerificationResultSuccess = "success"
	TicketVerificationResultFailed  = "failed"
)

type ElectronicTicket struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	OrderID            uint       `gorm:"column:order_id;not null;uniqueIndex" json:"orderId"`
	UserID             uint       `gorm:"column:user_id;not null;index" json:"userId"`
	TripID             uint       `gorm:"column:trip_id;not null;index" json:"tripId"`
	Token              string     `gorm:"type:text;not null" json:"token"`
	TokenHash          string     `gorm:"column:token_hash;size:64;uniqueIndex;not null" json:"-"`
	Status             string     `gorm:"size:20;not null;default:issued;index" json:"status"`
	ExpiresAt          time.Time  `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	VerifiedAt         *time.Time `gorm:"column:verified_at" json:"verifiedAt,omitempty"`
	VerifiedByDriverID *uint      `gorm:"column:verified_by_driver_id;index" json:"verifiedByDriverId,omitempty"`
	Order              Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Trip               Trip       `gorm:"foreignKey:TripID" json:"trip,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

func (ElectronicTicket) TableName() string {
	return "electronic_tickets"
}

type TicketVerification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"column:ticket_id;not null;index" json:"ticketId"`
	OrderID   uint      `gorm:"column:order_id;not null;index" json:"orderId"`
	DriverID  uint      `gorm:"column:driver_id;not null;index" json:"driverId"`
	Result    string    `gorm:"size:20;not null;index" json:"result"`
	Message   string    `gorm:"size:255" json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

func (TicketVerification) TableName() string {
	return "ticket_verifications"
}
