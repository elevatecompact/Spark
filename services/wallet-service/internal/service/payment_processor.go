package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

type PaymentResult struct {
	Success     bool
	ExternalRef string
	Error       string
}

type PaymentProcessor interface {
	Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error)
	Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error)
	Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error)
	Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error)
}

type noopPaymentProcessor struct{}

func NewNoopPaymentProcessor() PaymentProcessor {
	return &noopPaymentProcessor{}
}

func (p *noopPaymentProcessor) Deposit(ctx context.Context, userID string, amountCents int64, currency string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Msg("noop: deposit")
	return &PaymentResult{Success: true, ExternalRef: "noop-" + userID}, nil
}

func (p *noopPaymentProcessor) Withdraw(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Str("method", method).Msg("noop: withdraw")
	return &PaymentResult{Success: true, ExternalRef: "noop-" + userID}, nil
}

func (p *noopPaymentProcessor) Payout(ctx context.Context, userID string, amountCents int64, currency string, method string) (*PaymentResult, error) {
	log.Debug().Str("user_id", userID).Int64("amount_cents", amountCents).Str("method", method).Msg("noop: payout")
	return &PaymentResult{Success: true, ExternalRef: "noop-" + userID}, nil
}

func (p *noopPaymentProcessor) Refund(ctx context.Context, externalRef string, amountCents int64) (*PaymentResult, error) {
	log.Debug().Str("external_ref", externalRef).Int64("amount_cents", amountCents).Msg("noop: refund")
	return &PaymentResult{Success: true, ExternalRef: "refund-" + externalRef}, nil
}
