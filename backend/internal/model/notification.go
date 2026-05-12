package model

import "time"

const (
	NotificationTypeRefundApproved = "refund_approved"
	NotificationTypeRefundRejected = "refund_rejected"
	NotificationTypeOrderExpired   = "order_expired"
)

type Notification struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"column:user_id;not null;index:idx_notifications_user_read_created,priority:1" json:"userId"`
	Type           string     `gorm:"size:50;not null" json:"type"`
	Title          string     `gorm:"size:100;not null" json:"title"`
	Content        string     `gorm:"type:text;not null" json:"content"`
	RelatedOrderID *uint      `gorm:"column:related_order_id;index" json:"relatedOrderId,omitempty"`
	IsRead         bool       `gorm:"column:is_read;not null;default:false;index:idx_notifications_user_read_created,priority:2" json:"isRead"`
	ReadAt         *time.Time `gorm:"column:read_at" json:"readAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (Notification) TableName() string {
	return "notifications"
}
