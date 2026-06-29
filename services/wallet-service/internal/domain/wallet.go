package domain

import (
	"time"

	"github.com/google/uuid"
)

type WalletStatus string

const (
	WalletActive WalletStatus = "active"
	WalletFrozen WalletStatus = "frozen"
	WalletClosed WalletStatus = "closed"
)

type Currency string

const (
	CurrencyUSD   Currency = "USD"
	CurrencyToken Currency = "TITAN"
)

type Wallet struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	BalanceCents int64        `json:"balance_cents"`
	Currency     Currency     `json:"currency"`
	Status       WalletStatus `json:"status"`
	Version      int          `json:"version"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
