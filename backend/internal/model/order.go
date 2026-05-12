package model

import "time"

const (
	SeatTypeStandard = "standard"

	PayStatusUnpaid = "unpaid"
	PayStatusPaid   = "paid"

	OrderStatusPendingPayment      = "pending_payment"
	OrderStatusPendingVerification = "pending_verification"
	OrderStatusCompleted           = "completed"
	OrderStatusCancelled           = "cancelled"

	RefundStatusNone      = "none"
	RefundStatusRequested = "requested"
	RefundStatusRefunded  = "refunded"
	RefundStatusRejected  = "rejected"
)

type Order struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	OrderNo          string     `gorm:"column:order_no;size:32;uniqueIndex;not null" json:"orderNo"`
	UserID           uint       `gorm:"column:user_id;not null;index" json:"userId"`
	TripID           uint       `gorm:"column:trip_id;not null;index" json:"tripId"`
	TicketCount      int        `gorm:"column:ticket_count;not null" json:"ticketCount"`
	SeatType         string     `gorm:"column:seat_type;size:30;not null;default:standard" json:"seatType"`
	Amount           int        `gorm:"not null" json:"amount"`
	PayStatus        string     `gorm:"column:pay_status;size:20;not null;default:unpaid" json:"payStatus"`
	OrderStatus      string     `gorm:"column:order_status;size:30;not null;default:pending_payment" json:"orderStatus"`
	RefundStatus     string     `gorm:"column:refund_status;size:20;not null;default:none" json:"refundStatus"`
	RefundReviewNote string     `gorm:"column:refund_review_note;size:255" json:"refundReviewNote"`
	RefundReviewedAt *time.Time `gorm:"column:refund_reviewed_at" json:"refundReviewedAt,omitempty"`
	PaymentExpireAt  *time.Time `gorm:"column:payment_expire_at;index" json:"paymentExpireAt,omitempty"`
	User             User       `gorm:"foreignKey:UserID" json:"user"`
	Trip             Trip       `gorm:"foreignKey:TripID" json:"trip"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func (Order) TableName() string {
	return "orders"
}
