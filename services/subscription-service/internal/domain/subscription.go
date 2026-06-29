package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubActive      SubscriptionStatus = "active"
	SubCancelled   SubscriptionStatus = "cancelled"
	SubExpired     SubscriptionStatus = "expired"
	SubGracePeriod SubscriptionStatus = "grace_period"
)

type Subscription struct {
	ID                uuid.UUID          `json:"id"`
	UserID            uuid.UUID          `json:"user_id"`
	PlanID            uuid.UUID          `json:"plan_id"`
	Status            SubscriptionStatus `json:"status"`
	CurrentPeriodStart time.Time         `json:"current_period_start"`
	CurrentPeriodEnd  time.Time          `json:"current_period_end"`
	CancelledAt       *time.Time         `json:"cancelled_at,omitempty"`
	GracePeriodEnd    *time.Time         `json:"grace_period_end,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	PlanID uuid.UUID `json:"plan_id"`
}
