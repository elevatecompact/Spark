package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
	"github.com/elevatecompact/spark/services/wallet-service/internal/events"
	"github.com/elevatecompact/spark/services/wallet-service/internal/repository"
)

type PayoutService interface {
	Request(ctx context.Context, userID uuid.UUID, req domain.CreatePayoutRequest) (*domain.Payout, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Payout, error)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Payout, error)
}

type payoutService struct {
	walletRepo  repository.WalletRepository
	payoutRepo  repository.PayoutRepository
	eventPub    events.EventProducer
	paymentProc PaymentProcessor
	minPayout   int64
}

func NewPayoutService(
	walletRepo repository.WalletRepository,
	payoutRepo repository.PayoutRepository,
	eventPub events.EventProducer,
	paymentProc PaymentProcessor,
	minPayout int64,
) PayoutService {
	return &payoutService{
		walletRepo:  walletRepo,
		payoutRepo:  payoutRepo,
		eventPub:    eventPub,
		paymentProc: paymentProc,
		minPayout:   minPayout,
	}
}

func (s *payoutService) Request(ctx context.Context, userID uuid.UUID, req domain.CreatePayoutRequest) (*domain.Payout, error) {
	if req.AmountCents < s.minPayout {
		return nil, domain.NewDomainErrorMsg(domain.ErrPayoutMinimum, fmt.Sprintf("minimum payout is %d cents", s.minPayout), 400)
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet.Status != domain.WalletActive {
		return nil, domain.NewDomainError(domain.ErrWalletFrozen, 403)
	}
	if wallet.BalanceCents < req.AmountCents {
		return nil, domain.NewDomainError(domain.ErrInsufficientFunds, 402)
	}

	payout := &domain.Payout{
		ID:          uuid.New(),
		WalletID:    wallet.ID,
		AmountCents: req.AmountCents,
		Currency:    req.Currency,
		Method:      req.Method,
		Status:      domain.PayoutRequested,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.payoutRepo.Create(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, wallet, -req.AmountCents); err != nil {
		s.payoutRepo.UpdateStatus(ctx, payout.ID, domain.PayoutFailed, nil)
		return nil, fmt.Errorf("failed to hold funds for payout: %w", err)
	}

	payout.Status = domain.PayoutProcessing
	s.payoutRepo.UpdateStatus(ctx, payout.ID, domain.PayoutProcessing, nil)

	result, err := s.paymentProc.Payout(ctx, userID.String(), req.AmountCents, string(req.Currency), string(req.Method))
	if err != nil || !result.Success {
		s.walletRepo.UpdateBalance(ctx, wallet, req.AmountCents)
		s.payoutRepo.UpdateStatus(ctx, payout.ID, domain.PayoutFailed, &result.Error)
		return nil, fmt.Errorf("payout processing failed: %s", result.Error)
	}

	completedAt := time.Now().UTC()
	payout.Status = domain.PayoutCompleted
	payout.CompletedAt = &completedAt
	payout.ExternalRef = &result.ExternalRef
	s.payoutRepo.UpdateStatus(ctx, payout.ID, domain.PayoutCompleted, &result.ExternalRef)

	if err := s.eventPub.PublishPayoutCompleted(ctx, payout); err != nil {
		log.Warn().Err(err).Msg("failed to publish payout event")
	}

	return payout, nil
}

func (s *payoutService) Get(ctx context.Context, id uuid.UUID) (*domain.Payout, error) {
	return s.payoutRepo.GetByID(ctx, id)
}

func (s *payoutService) ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Payout, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.payoutRepo.ListByWallet(ctx, wallet.ID, cursor, limit)
}
