package domain

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoicePending  InvoiceStatus = "pending"
	InvoicePaid     InvoiceStatus = "paid"
	InvoiceFailed   InvoiceStatus = "failed"
	InvoiceRefunded InvoiceStatus = "refunded"
)

type Invoice struct {
	ID             uuid.UUID      `json:"id"`
	SubscriptionID uuid.UUID      `json:"subscription_id"`
	AmountCents    int64          `json:"amount_cents"`
	Currency       string         `json:"currency"`
	Status         InvoiceStatus  `json:"status"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	PeriodStart    time.Time      `json:"period_start"`
	PeriodEnd      time.Time      `json:"period_end"`
	CreatedAt      time.Time      `json:"created_at"`
}
