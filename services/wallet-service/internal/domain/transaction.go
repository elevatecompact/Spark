package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TxnDeposit    TransactionType = "deposit"
	TxnWithdraw   TransactionType = "withdraw"
	TxnTransfer   TransactionType = "transfer"
	TxnTip        TransactionType = "tip"
	TxnPurchase   TransactionType = "purchase"
	TxnPayout     TransactionType = "payout"
	TxnRefund     TransactionType = "refund"
	TxnFee        TransactionType = "fee"
)

type TransactionStatus string

const (
	TxnPending  TransactionStatus = "pending"
	TxnSettled  TransactionStatus = "settled"
	TxnFailed   TransactionStatus = "failed"
)

type Transaction struct {
	ID             uuid.UUID        `json:"id"`
	IdempotencyKey string           `json:"idempotency_key"`
	FromWalletID   *uuid.UUID       `json:"from_wallet_id,omitempty"`
	ToWalletID     *uuid.UUID       `json:"to_wallet_id,omitempty"`
	AmountCents    int64            `json:"amount_cents"`
	Currency       Currency         `json:"currency"`
	Type           TransactionType  `json:"type"`
	Status         TransactionStatus `json:"status"`
	FailureReason  *string          `json:"failure_reason,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	SettledAt      *time.Time       `json:"settled_at,omitempty"`
}

type CreateTransactionRequest struct {
	FromWalletID   *uuid.UUID      `json:"from_wallet_id,omitempty"`
	ToWalletID     *uuid.UUID      `json:"to_wallet_id,omitempty"`
	AmountCents    int64           `json:"amount_cents"`
	Currency       Currency        `json:"currency"`
	Type           TransactionType `json:"type"`
	IdempotencyKey string          `json:"idempotency_key"`
}
