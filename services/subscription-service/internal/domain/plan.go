package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BillingPeriod string

const (
	BillingMonthly BillingPeriod = "monthly"
	BillingYearly  BillingPeriod = "yearly"
)

type SubscriptionPlan struct {
	ID           uuid.UUID       `json:"id"`
	CreatorID    *uuid.UUID      `json:"creator_id,omitempty"`
	Name         string          `json:"name"`
	PriceCents   int64           `json:"price_cents"`
	Currency     string          `json:"currency"`
	BillingPeriod BillingPeriod  `json:"billing_period"`
	Benefits     json.RawMessage `json:"benefits"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CreatePlanRequest struct {
	CreatorID     *uuid.UUID      `json:"creator_id,omitempty"`
	Name          string          `json:"name"`
	PriceCents    int64           `json:"price_cents"`
	Currency      string          `json:"currency"`
	BillingPeriod BillingPeriod   `json:"billing_period"`
	Benefits      json.RawMessage `json:"benefits"`
}

type UpdatePlanRequest struct {
	Name          *string          `json:"name,omitempty"`
	PriceCents    *int64           `json:"price_cents,omitempty"`
	BillingPeriod *BillingPeriod   `json:"billing_period,omitempty"`
	Benefits      *json.RawMessage `json:"benefits,omitempty"`
	IsActive      *bool            `json:"is_active,omitempty"`
}
