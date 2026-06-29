package domain

import (
	"time"

	"github.com/google/uuid"
)

type PayoutMethod string

const (
	PayoutPayPal PayoutMethod = "paypal"
	PayoutBank   PayoutMethod = "bank"
	PayoutCrypto PayoutMethod = "crypto"
)

type PayoutStatus string

const (
	PayoutRequested  PayoutStatus = "requested"
	PayoutProcessing PayoutStatus = "processing"
	PayoutCompleted  PayoutStatus = "completed"
	PayoutFailed     PayoutStatus = "failed"
)

type Payout struct {
	ID           uuid.UUID    `json:"id"`
	WalletID     uuid.UUID    `json:"wallet_id"`
	AmountCents  int64        `json:"amount_cents"`
	Currency     Currency     `json:"currency"`
	Method       PayoutMethod `json:"method"`
	Status       PayoutStatus `json:"status"`
	ExternalRef  *string      `json:"external_ref,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
}

type CreatePayoutRequest struct {
	AmountCents int64        `json:"amount_cents"`
	Currency    Currency     `json:"currency"`
	Method      PayoutMethod `json:"method"`
}
