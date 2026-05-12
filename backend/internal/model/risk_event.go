package model

import "time"

const (
	RiskSeverityHigh   = "high"
	RiskSeverityMedium = "medium"
	RiskSeverityLow    = "low"

	RiskStatusOpen         = "open"
	RiskStatusAcknowledged = "acknowledged"
	RiskStatusResolved     = "resolved"

	RiskEventTypeAIRateLimit = "ai_rate_limit"
	RiskEventTypeTokenSpike  = "token_spike"
)

type RiskEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Severity    string    `gorm:"size:20;not null;index" json:"severity"`
	EventType   string    `gorm:"size:50;not null;index" json:"eventType"`
	SubjectType string    `gorm:"size:30;not null;index" json:"subjectType"`
	SubjectID   string    `gorm:"size:100;not null;index" json:"subjectId"`
	Fingerprint string    `gorm:"size:150;not null;index" json:"fingerprint"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Detail      string    `gorm:"type:text;not null" json:"detail"`
	Status      string    `gorm:"size:20;not null;default:'open';index" json:"status"`
	MetricsJSON string    `gorm:"type:longtext" json:"metricsJson"`
	CreatedAt   time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (RiskEvent) TableName() string {
	return "risk_events"
}
