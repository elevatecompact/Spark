package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
	"github.com/elevatecompact/spark/services/subscription-service/internal/events"
	"github.com/elevatecompact/spark/services/subscription-service/internal/repository"
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, userID uuid.UUID, req domain.CreateSubscriptionRequest) (*domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	GetMy(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error)
	ChangePlan(ctx context.Context, subID, userID, newPlanID uuid.UUID) (*domain.Subscription, error)
	Cancel(ctx context.Context, subID, userID uuid.UUID) error
	Reactivate(ctx context.Context, subID, userID uuid.UUID) error
}

type subscriptionService struct {
	subRepo   repository.SubscriptionRepository
	planRepo  repository.PlanRepository
	invRepo   repository.InvoiceRepository
	eventPub  events.EventProducer
	graceDays int
	maxActive int
	trialDays int
}

func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
	invRepo repository.InvoiceRepository,
	eventPub events.EventProducer,
	graceDays, maxActive, trialDays int,
) SubscriptionService {
	return &subscriptionService{
		subRepo:   subRepo,
		planRepo:  planRepo,
		invRepo:   invRepo,
		eventPub:  eventPub,
		graceDays: graceDays,
		maxActive: maxActive,
		trialDays: trialDays,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, userID uuid.UUID, req domain.CreateSubscriptionRequest) (*domain.Subscription, error) {
	plan, err := s.planRepo.GetByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}
	if !plan.IsActive {
		return nil, domain.ErrPlanInactive
	}

	existing, err := s.subRepo.GetByUserAndPlan(ctx, userID, req.PlanID)
	if err == nil && existing != nil {
		if existing.Status == domain.SubActive || existing.Status == domain.SubGracePeriod {
			return nil, domain.ErrAlreadySubscribed
		}
	}

	count, err := s.subRepo.CountActiveByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= s.maxActive {
		return nil, domain.ErrMaxSubscriptions
	}

	now := time.Now().UTC()
	var periodEnd time.Time
	switch plan.BillingPeriod {
	case domain.BillingMonthly:
		periodEnd = now.AddDate(0, 1, 0)
	case domain.BillingYearly:
		periodEnd = now.AddDate(1, 0, 0)
	}

	sub := &domain.Subscription{
		ID:                 uuid.New(),
		UserID:             userID,
		PlanID:             req.PlanID,
		Status:             domain.SubActive,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   periodEnd,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	inv := &domain.Invoice{
		ID:             uuid.New(),
		SubscriptionID: sub.ID,
		AmountCents:    plan.PriceCents,
		Currency:       plan.Currency,
		Status:         domain.InvoicePending,
		PeriodStart:    now,
		PeriodEnd:      periodEnd,
		CreatedAt:      now,
	}
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	if err := s.eventPub.PublishActivated(ctx, sub); err != nil {
		log.Warn().Err(err).Msg("failed to publish subscription.activated")
	}

	return sub, nil
}

func (s *subscriptionService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.subRepo.GetByID(ctx, id)
}

func (s *subscriptionService) GetMy(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	return s.subRepo.ListByUser(ctx, userID)
}

func (s *subscriptionService) ChangePlan(ctx context.Context, subID, userID, newPlanID uuid.UUID) (*domain.Subscription, error) {
	sub, err := s.subRepo.GetByID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID {
		return nil, domain.ErrNotOwner
	}
	if sub.Status != domain.SubActive {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "subscription is not active", 400)
	}

	oldPlan, err := s.planRepo.GetByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	newPlan, err := s.planRepo.GetByID(ctx, newPlanID)
	if err != nil {
		return nil, err
	}
	if !newPlan.IsActive {
		return nil, domain.ErrPlanInactive
	}

	// Proration: charge or credit the user the prorated difference between the
	// old and new plans for the remainder of the current billing period.
	prorated := computeProration(oldPlan.PriceCents, newPlan.PriceCents, sub.CurrentPeriodStart, sub.CurrentPeriodEnd)
	if err := s.recordProrationInvoice(ctx, sub, newPlan, prorated); err != nil {
		return nil, err
	}

	sub.PlanID = newPlanID
	sub.UpdatedAt = time.Now().UTC()
	if err := s.subRepo.Update(ctx, sub); err != nil {
		return nil, err
	}

	if err := s.eventPub.PublishUpgraded(ctx, sub); err != nil {
		log.Warn().Err(err).Msg("failed to publish subscription.upgraded")
	}

	return sub, nil
}

func (s *subscriptionService) recordProrationInvoice(ctx context.Context, sub *domain.Subscription, newPlan *domain.SubscriptionPlan, proratedCents int64) error {
	if proratedCents == 0 {
		return nil
	}
	status := domain.InvoicePending
	if proratedCents < 0 {
		// Downgrade / refund: issue a credit invoice for the unused amount.
		status = domain.InvoiceRefunded
	}
	inv := &domain.Invoice{
		ID:             uuid.New(),
		SubscriptionID: sub.ID,
		AmountCents:    absInt64(proratedCents),
		Currency:       newPlan.Currency,
		Status:         status,
		PeriodStart:    time.Now().UTC(),
		PeriodEnd:      sub.CurrentPeriodEnd,
		CreatedAt:      time.Now().UTC(),
	}
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return fmt.Errorf("failed to create proration invoice: %w", err)
	}
	log.Info().
		Str("subscription_id", sub.ID.String()).
		Int64("prorated_cents", proratedCents).
		Str("status", string(status)).
		Msg("proration invoice recorded")
	return nil
}

// computeProration returns the prorated difference (newPlan - oldPlan) for the
// remainder of the current period. A positive value means the user owes money
// (upgrade), a negative value means a credit (downgrade).
func computeProration(oldPrice, newPrice int64, periodStart, periodEnd time.Time) int64 {
	total := periodEnd.Sub(periodStart)
	if total <= 0 {
		return newPrice - oldPrice
	}
	remaining := periodEnd.Sub(time.Now().UTC())
	if remaining <= 0 {
		return 0
	}
	// Use millisecond resolution so very short windows still produce sensible
	// numbers and the multiplication cannot overflow int64 for any realistic
	// subscription price (< 1e15 cents) or window (< 1e12 ms).
	fractionMs := float64(remaining.Milliseconds()) / float64(total.Milliseconds())
	diff := float64(newPrice - oldPrice)
	prorated := int64(diff * fractionMs)
	return prorated
}

func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func (s *subscriptionService) Cancel(ctx context.Context, subID, userID uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, subID)
	if err != nil {
		return err
	}
	if sub.UserID != userID {
		return domain.ErrNotOwner
	}
	if sub.Status != domain.SubActive {
		return domain.NewDomainErrorMsg(domain.ErrValidation, "subscription is not active", 400)
	}

	now := time.Now().UTC()
	sub.Status = domain.SubCancelled
	sub.CancelledAt = &now
	sub.UpdatedAt = now
	if err := s.subRepo.Update(ctx, sub); err != nil {
		return err
	}

	return s.eventPub.PublishCancelled(ctx, sub)
}

func (s *subscriptionService) Reactivate(ctx context.Context, subID, userID uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, subID)
	if err != nil {
		return err
	}
	if sub.UserID != userID {
		return domain.ErrNotOwner
	}
	if sub.Status != domain.SubCancelled {
		return domain.NewDomainErrorMsg(domain.ErrValidation, "only cancelled subscriptions can be reactivated", 400)
	}

	plan, err := s.planRepo.GetByID(ctx, sub.PlanID)
	if err != nil {
		return err
	}
	if !plan.IsActive {
		return domain.ErrPlanInactive
	}

	now := time.Now().UTC()
	var periodEnd time.Time
	switch plan.BillingPeriod {
	case domain.BillingMonthly:
		periodEnd = now.AddDate(0, 1, 0)
	case domain.BillingYearly:
		periodEnd = now.AddDate(1, 0, 0)
	}

	sub.Status = domain.SubActive
	sub.CancelledAt = nil
	sub.CurrentPeriodStart = now
	sub.CurrentPeriodEnd = periodEnd
	sub.UpdatedAt = now
	if err := s.subRepo.Update(ctx, sub); err != nil {
		return err
	}

	inv := &domain.Invoice{
		ID:             uuid.New(),
		SubscriptionID: sub.ID,
		AmountCents:    plan.PriceCents,
		Currency:       plan.Currency,
		Status:         domain.InvoicePending,
		PeriodStart:    now,
		PeriodEnd:      periodEnd,
		CreatedAt:      now,
	}
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return err
	}

	return s.eventPub.PublishActivated(ctx, sub)
}
