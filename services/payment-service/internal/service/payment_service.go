package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
	"github.com/elevatecompact/spark/services/payment-service/internal/events"
	"github.com/elevatecompact/spark/services/payment-service/internal/processor"
	"github.com/elevatecompact/spark/services/payment-service/internal/repository"
)

type PaymentService interface {
	// Intents
	CreateIntent(ctx context.Context, userID uuid.UUID, req domain.CreateIntentRequest) (*domain.PaymentIntent, error)
	GetIntent(ctx context.Context, id uuid.UUID) (*domain.PaymentIntent, error)
	ListIntents(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.PaymentIntent, error)
	ConfirmIntent(ctx context.Context, id, userID uuid.UUID, req domain.ConfirmIntentRequest) (*domain.PaymentIntent, error)
	CancelIntent(ctx context.Context, id, userID uuid.UUID) error

	// Methods
	CreatePaymentMethod(ctx context.Context, userID uuid.UUID, req domain.CreatePaymentMethodRequest) (*domain.PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error)
	ListPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, id, userID uuid.UUID) error
	DeletePaymentMethod(ctx context.Context, id, userID uuid.UUID) error

	// Refunds
	RefundIntent(ctx context.Context, id, userID uuid.UUID, req domain.RefundRequest) error

	// Payouts
	CreatePayout(ctx context.Context, userID uuid.UUID, req domain.CreatePayoutRequest) (*domain.Payout, error)
	GetPayout(ctx context.Context, id uuid.UUID) (*domain.Payout, error)

	// Webhooks
	ProcessWebhook(ctx context.Context, proc domain.PaymentProcessor, externalID, eventType string, body []byte) error
	RetryWebhook(ctx context.Context, id uuid.UUID) error

	// Admin
	GetProcessorStatus(ctx context.Context) map[string]bool
}

type paymentService struct {
	intentRepo  repository.PaymentIntentRepository
	methodRepo  repository.PaymentMethodRepository
	webhookRepo repository.WebhookRepository
	processors  map[domain.PaymentProcessor]processor.PaymentProcessor
	eventPub    events.EventProducer

	stripeEnabled bool
	paypalEnabled bool
	saveMethods   bool
	refundsEnabled bool
}

func NewPaymentService(
	intentRepo repository.PaymentIntentRepository,
	methodRepo repository.PaymentMethodRepository,
	webhookRepo repository.WebhookRepository,
	procMap map[domain.PaymentProcessor]processor.PaymentProcessor,
	eventPub events.EventProducer,
	stripeEnabled, paypalEnabled, saveMethods, refundsEnabled bool,
) PaymentService {
	return &paymentService{
		intentRepo:     intentRepo,
		methodRepo:     methodRepo,
		webhookRepo:    webhookRepo,
		processors:     procMap,
		eventPub:       eventPub,
		stripeEnabled:  stripeEnabled,
		paypalEnabled:  paypalEnabled,
		saveMethods:    saveMethods,
		refundsEnabled: refundsEnabled,
	}
}

func (s *paymentService) getProcessor(proc domain.PaymentProcessor) (processor.PaymentProcessor, error) {
	p, ok := s.processors[proc]
	if !ok {
		return nil, domain.NewDomainErrorMsg(domain.ErrProcessorDisabled, fmt.Sprintf("processor %s not configured", proc), 400)
	}
	if proc == domain.ProcessorStripe && !s.stripeEnabled {
		return nil, domain.ErrProcessorDisabled
	}
	if proc == domain.ProcessorPayPal && !s.paypalEnabled {
		return nil, domain.ErrProcessorDisabled
	}
	return p, nil
}

func (s *paymentService) CreateIntent(ctx context.Context, userID uuid.UUID, req domain.CreateIntentRequest) (*domain.PaymentIntent, error) {
	if req.AmountCents <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if req.IdempotencyKey != "" {
		existing, err := s.intentRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	proc := domain.ProcessorStripe
	procImpl, err := s.getProcessor(proc)
	if err != nil {
		if s.paypalEnabled {
			proc = domain.ProcessorPayPal
			procImpl, err = s.getProcessor(proc)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	now := time.Now().UTC()
	intent := &domain.PaymentIntent{
		ID:             uuid.New(),
		UserID:         userID,
		Processor:      proc,
		AmountCents:    req.AmountCents,
		Currency:       currency,
		Status:         domain.IntentRequiresPaymentMethod,
		IdempotencyKey: req.IdempotencyKey,
		Metadata:       req.Metadata,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := procImpl.CreateIntent(ctx, intent); err != nil {
		return nil, domain.NewDomainErrorMsg(domain.ErrProcessorError, fmt.Sprintf("processor error: %v", err), 500)
	}

	if err := s.intentRepo.Create(ctx, intent); err != nil {
		return nil, fmt.Errorf("failed to create intent: %w", err)
	}

	if err := s.eventPub.PublishIntentCreated(ctx, intent); err != nil {
		log.Warn().Err(err).Msg("failed to publish intent.created")
	}

	return intent, nil
}

func (s *paymentService) GetIntent(ctx context.Context, id uuid.UUID) (*domain.PaymentIntent, error) {
	return s.intentRepo.GetByID(ctx, id)
}

func (s *paymentService) ListIntents(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.PaymentIntent, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.intentRepo.ListByUser(ctx, userID, cursor, limit)
}

func (s *paymentService) ConfirmIntent(ctx context.Context, id, userID uuid.UUID, req domain.ConfirmIntentRequest) (*domain.PaymentIntent, error) {
	intent, err := s.intentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if intent.UserID != userID {
		return nil, domain.ErrForbidden
	}
	if intent.Status != domain.IntentRequiresPaymentMethod {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "intent cannot be confirmed in current status", 400)
	}

	method, err := s.methodRepo.GetByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, err
	}
	if method.UserID != userID {
		return nil, domain.ErrForbidden
	}

	procImpl, err := s.getProcessor(intent.Processor)
	if err != nil {
		return nil, err
	}

	if err := procImpl.ConfirmIntent(ctx, intent, method.ExternalID); err != nil {
		_ = s.intentRepo.UpdateStatus(ctx, id, domain.IntentFailed)
		_ = s.eventPub.PublishIntentFailed(ctx, intent)
		return nil, domain.NewDomainErrorMsg(domain.ErrConfirmFailed, fmt.Sprintf("confirmation failed: %v", err), 500)
	}

	intent.PaymentMethodID = &method.ID
	if err := s.intentRepo.Update(ctx, intent); err != nil {
		return nil, err
	}

	intent.Status = domain.IntentSucceeded
	if err := s.intentRepo.UpdateStatus(ctx, id, domain.IntentSucceeded); err != nil {
		return nil, err
	}

	if err := s.eventPub.PublishIntentSucceeded(ctx, intent); err != nil {
		log.Warn().Err(err).Msg("failed to publish intent.succeeded")
	}

	return intent, nil
}

func (s *paymentService) CancelIntent(ctx context.Context, id, userID uuid.UUID) error {
	intent, err := s.intentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if intent.UserID != userID {
		return domain.ErrForbidden
	}
	if intent.Status != domain.IntentRequiresPaymentMethod {
		return domain.NewDomainErrorMsg(domain.ErrValidation, "intent cannot be canceled in current status", 400)
	}

	procImpl, err := s.getProcessor(intent.Processor)
	if err != nil {
		return err
	}

	if err := procImpl.CancelIntent(ctx, intent); err != nil {
		return domain.NewDomainErrorMsg(domain.ErrCancelFailed, fmt.Sprintf("cancellation failed: %v", err), 500)
	}

	return s.intentRepo.UpdateStatus(ctx, id, domain.IntentCanceled)
}

func (s *paymentService) CreatePaymentMethod(ctx context.Context, userID uuid.UUID, req domain.CreatePaymentMethodRequest) (*domain.PaymentMethod, error) {
	if !s.saveMethods {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "saving payment methods is disabled", 400)
	}

	now := time.Now().UTC()
	method := &domain.PaymentMethod{
		ID:         uuid.New(),
		UserID:     userID,
		Processor:  req.Processor,
		Type:       req.Type,
		ExternalID: req.Token,
		Fingerprint: "fp_noop_" + uuid.New().String(),
		Last4:      "4242",
		ExpMonth:   12,
		ExpYear:    time.Now().UTC().Year() + 3,
		IsDefault:  req.SetAsDefault,
		CreatedAt:  now,
	}

	if req.SetAsDefault {
		existing, _ := s.methodRepo.ListByUser(ctx, userID)
		for _, m := range existing {
			if m.IsDefault {
				method.IsDefault = true
				break
			}
		}
	}

	if err := s.methodRepo.Create(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return method, nil
}

func (s *paymentService) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error) {
	return s.methodRepo.GetByID(ctx, id)
}

func (s *paymentService) ListPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentMethod, error) {
	return s.methodRepo.ListByUser(ctx, userID)
}

func (s *paymentService) SetDefaultPaymentMethod(ctx context.Context, id, userID uuid.UUID) error {
	return s.methodRepo.SetDefault(ctx, id, userID)
}

func (s *paymentService) DeletePaymentMethod(ctx context.Context, id, userID uuid.UUID) error {
	method, err := s.methodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if method.UserID != userID {
		return domain.ErrForbidden
	}
	return s.methodRepo.Delete(ctx, id)
}

func (s *paymentService) RefundIntent(ctx context.Context, id, userID uuid.UUID, req domain.RefundRequest) error {
	if !s.refundsEnabled {
		return domain.NewDomainErrorMsg(domain.ErrValidation, "refunds are disabled", 400)
	}

	intent, err := s.intentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if intent.UserID != userID {
		return domain.ErrForbidden
	}
	if intent.Status != domain.IntentSucceeded {
		return domain.NewDomainErrorMsg(domain.ErrValidation, "only succeeded intents can be refunded", 400)
	}

	if req.AmountCents != nil {
		if *req.AmountCents <= 0 || *req.AmountCents > intent.AmountCents {
			return domain.ErrInvalidAmount
		}
	}

	procImpl, err := s.getProcessor(intent.Processor)
	if err != nil {
		return err
	}

	refundID, err := procImpl.Refund(ctx, intent, req.AmountCents)
	if err != nil {
		return domain.NewDomainErrorMsg(domain.ErrRefundFailed, fmt.Sprintf("refund failed: %v", err), 500)
	}

	if err := s.eventPub.PublishRefundProcessed(ctx, intent, refundID); err != nil {
		log.Warn().Err(err).Msg("failed to publish refund.processed")
	}

	return nil
}

func (s *paymentService) CreatePayout(ctx context.Context, userID uuid.UUID, req domain.CreatePayoutRequest) (*domain.Payout, error) {
	if req.AmountCents <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	proc := domain.ProcessorStripe
	procImpl, err := s.getProcessor(proc)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	payout := &domain.Payout{
		ID:         uuid.New(),
		UserID:     userID,
		Processor:  proc,
		AmountCents: req.AmountCents,
		Currency:   currency,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := procImpl.CreatePayout(ctx, payout); err != nil {
		return nil, domain.NewDomainErrorMsg(domain.ErrProcessorError, fmt.Sprintf("payout failed: %v", err), 500)
	}

	if err := s.eventPub.PublishPayoutCompleted(ctx, payout); err != nil {
		log.Warn().Err(err).Msg("failed to publish payout.completed")
	}

	return payout, nil
}

func (s *paymentService) GetPayout(ctx context.Context, id uuid.UUID) (*domain.Payout, error) {
	return nil, domain.ErrNotFound
}

func (s *paymentService) ProcessWebhook(ctx context.Context, proc domain.PaymentProcessor, externalID, eventType string, body []byte) error {
	existing, err := s.webhookRepo.GetByExternalEventID(ctx, proc, externalID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	now := time.Now().UTC()
	event := &domain.WebhookEvent{
		ID:              uuid.New(),
		Processor:       proc,
		ExternalEventID: externalID,
		Type:            eventType,
		Body:            body,
		Status:          domain.WebhookReceived,
		CreatedAt:       now,
	}

	if err := s.webhookRepo.Create(ctx, event); err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	if err := s.webhookRepo.UpdateStatus(ctx, event.ID, domain.WebhookProcessed); err != nil {
		log.Warn().Err(err).Msg("failed to mark webhook as processed")
	}

	return nil
}

func (s *paymentService) RetryWebhook(ctx context.Context, id uuid.UUID) error {
	event, err := s.webhookRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if event.Status == domain.WebhookProcessed {
		return nil
	}
	return s.webhookRepo.UpdateStatus(ctx, id, domain.WebhookProcessed)
}

func (s *paymentService) GetProcessorStatus(ctx context.Context) map[string]bool {
	return map[string]bool{
		"stripe": s.stripeEnabled,
		"paypal": s.paypalEnabled,
	}
}
