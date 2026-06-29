package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
	"github.com/elevatecompact/spark/services/wallet-service/internal/events"
	"github.com/elevatecompact/spark/services/wallet-service/internal/repository"
)

type WalletService interface {
	GetOrCreate(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	Freeze(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context, id uuid.UUID) error
}

type walletService struct {
	walletRepo repository.WalletRepository
	eventPub   events.EventProducer
	maxBalance int64
}

func NewWalletService(walletRepo repository.WalletRepository, eventPub events.EventProducer, maxBalance int64) WalletService {
	return &walletService{walletRepo: walletRepo, eventPub: eventPub, maxBalance: maxBalance}
}

func (s *walletService) GetOrCreate(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err == domain.ErrWalletNotFound {
		now := time.Now().UTC()
		wallet = &domain.Wallet{
			ID:           uuid.New(),
			UserID:       userID,
			BalanceCents: 0,
			Currency:     domain.CurrencyUSD,
			Status:       domain.WalletActive,
			Version:      1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.walletRepo.Create(ctx, wallet); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
		return wallet, nil
	}
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	return s.walletRepo.GetByID(ctx, id)
}

func (s *walletService) GetByUser(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	return s.walletRepo.GetByUserID(ctx, userID)
}

func (s *walletService) Freeze(ctx context.Context, id uuid.UUID) error {
	return s.walletRepo.UpdateStatus(ctx, id, domain.WalletFrozen)
}

func (s *walletService) Close(ctx context.Context, id uuid.UUID) error {
	return s.walletRepo.UpdateStatus(ctx, id, domain.WalletClosed)
}
