package processor

import (
	"context"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
)

type PaymentProcessor interface {
	Name() domain.PaymentProcessor
	CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error
	ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error
	CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error
	Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error)
	CreatePayout(ctx context.Context, payout *domain.Payout) error
}

type noopProcessor struct {
	name domain.PaymentProcessor
}

func NewNoopProcessor(name domain.PaymentProcessor) PaymentProcessor {
	return &noopProcessor{name: name}
}

func (p *noopProcessor) Name() domain.PaymentProcessor { return p.name }

func (p *noopProcessor) CreateIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	intent.ExternalID = string(p.name) + "_pi_noop_" + intent.ID.String()
	if intent.Status == "" {
		intent.Status = domain.IntentRequiresPaymentMethod
	}
	return nil
}

func (p *noopProcessor) ConfirmIntent(ctx context.Context, intent *domain.PaymentIntent, paymentMethodID string) error {
	intent.Status = domain.IntentSucceeded
	return nil
}

func (p *noopProcessor) CancelIntent(ctx context.Context, intent *domain.PaymentIntent) error {
	intent.Status = domain.IntentCanceled
	return nil
}

func (p *noopProcessor) Refund(ctx context.Context, intent *domain.PaymentIntent, amountCents *int64) (string, error) {
	return string(p.name) + "_rf_noop_" + intent.ID.String(), nil
}

func (p *noopProcessor) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	payout.ExternalID = string(p.name) + "_po_noop_" + payout.ID.String()
	payout.Status = "completed"
	return nil
}
