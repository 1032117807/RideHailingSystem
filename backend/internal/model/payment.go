package model

import "time"

const (
	PaymentChannelMock = "mock"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
	PaymentStatusClosed  = "closed"
)

type Payment struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	PaymentNo string     `gorm:"column:payment_no;size:32;uniqueIndex;not null" json:"paymentNo"`
	OrderID   uint       `gorm:"column:order_id;not null;index" json:"orderId"`
	UserID    uint       `gorm:"column:user_id;not null;index" json:"userId"`
	Amount    int        `gorm:"not null" json:"amount"`
	Channel   string     `gorm:"size:20;not null;default:mock" json:"channel"`
	Status    string     `gorm:"size:20;not null;default:pending" json:"status"`
	PaidAt    *time.Time `json:"paidAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (Payment) TableName() string {
	return "payments"
}
