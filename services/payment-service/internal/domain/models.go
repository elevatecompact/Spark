package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentProcessor string

const (
	ProcessorStripe PaymentProcessor = "stripe"
	ProcessorPayPal PaymentProcessor = "paypal"
)

type IntentStatus string

const (
	IntentRequiresPaymentMethod IntentStatus = "requires_payment_method"
	IntentProcessing            IntentStatus = "processing"
	IntentSucceeded             IntentStatus = "succeeded"
	IntentFailed                IntentStatus = "failed"
	IntentCanceled              IntentStatus = "canceled"
)

type PaymentMethodType string

const (
	MethodCard   PaymentMethodType = "card"
	MethodPayPal PaymentMethodType = "paypal"
	MethodBank   PaymentMethodType = "bank"
)

type WebhookStatus string

const (
	WebhookReceived  WebhookStatus = "received"
	WebhookProcessed WebhookStatus = "processed"
	WebhookFailed    WebhookStatus = "failed"
)

type PaymentIntent struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	ExternalID      string          `json:"external_id"`
	Processor       PaymentProcessor `json:"processor"`
	AmountCents     int64           `json:"amount_cents"`
	Currency        string          `json:"currency"`
	Status          IntentStatus    `json:"status"`
	IdempotencyKey  string          `json:"idempotency_key"`
	Metadata        json.RawMessage `json:"metadata"`
	PaymentMethodID *uuid.UUID      `json:"payment_method_id,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type PaymentMethod struct {
	ID         uuid.UUID         `json:"id"`
	UserID     uuid.UUID         `json:"user_id"`
	ExternalID string            `json:"external_id"`
	Processor  PaymentProcessor  `json:"processor"`
	Type       PaymentMethodType `json:"type"`
	Fingerprint string           `json:"fingerprint"`
	Last4      string            `json:"last4"`
	ExpMonth   int               `json:"exp_month"`
	ExpYear    int               `json:"exp_year"`
	IsDefault  bool              `json:"is_default"`
	CreatedAt  time.Time         `json:"created_at"`
}

type WebhookEvent struct {
	ID             uuid.UUID       `json:"id"`
	Processor      PaymentProcessor `json:"processor"`
	ExternalEventID string          `json:"external_event_id"`
	Type           string          `json:"type"`
	Body           json.RawMessage `json:"body"`
	Status         WebhookStatus   `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	ProcessedAt    *time.Time      `json:"processed_at,omitempty"`
}

type Payout struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	ExternalID string          `json:"external_id"`
	Processor  PaymentProcessor `json:"processor"`
	AmountCents int64          `json:"amount_cents"`
	Currency   string          `json:"currency"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type CreateIntentRequest struct {
	AmountCents    int64           `json:"amount_cents"`
	Currency       string          `json:"currency"`
	IdempotencyKey string          `json:"idempotency_key"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	PaymentMethodID *uuid.UUID    `json:"payment_method_id,omitempty"`
}

type ConfirmIntentRequest struct {
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
}

type RefundRequest struct {
	AmountCents *int64 `json:"amount_cents,omitempty"`
}

type CreatePaymentMethodRequest struct {
	Processor      PaymentProcessor  `json:"processor"`
	Type           PaymentMethodType `json:"type"`
	Token          string            `json:"token"`
	SetAsDefault   bool              `json:"set_as_default"`
}

type CreatePayoutRequest struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}
