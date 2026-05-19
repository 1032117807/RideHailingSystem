package model

import "time"

type RefundAuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"column:order_id;not null;index" json:"orderId"`
	RefundStatus string    `gorm:"column:refund_status;size:20;not null;index" json:"refundStatus"`
	ReviewNote   string    `gorm:"column:review_note;size:255" json:"reviewNote"`
	ReviewerID   uint      `gorm:"column:reviewer_id;not null;index" json:"reviewerId"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (RefundAuditLog) TableName() string {
	return "refund_audit_logs"
}

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ActorUserID  uint      `gorm:"column:actor_user_id;not null;index" json:"actorUserId"`
	ActorRole    string    `gorm:"column:actor_role;size:20;not null" json:"actorRole"`
	Action       string    `gorm:"size:80;not null;index" json:"action"`
	ResourceType string    `gorm:"column:resource_type;size:80;not null;index" json:"resourceType"`
	ResourceID   string    `gorm:"column:resource_id;size:80;not null;index" json:"resourceId"`
	Detail       string    `gorm:"type:text" json:"detail"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
