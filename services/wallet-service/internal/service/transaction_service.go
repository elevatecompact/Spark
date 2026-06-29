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

type TransactionService interface {
	Deposit(ctx context.Context, userID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	Withdraw(ctx context.Context, userID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	Transfer(ctx context.Context, fromUserID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	Tip(ctx context.Context, fromUserID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Transaction, error)
}

type transactionService struct {
	walletRepo repository.WalletRepository
	txnRepo    repository.TransactionRepository
	eventPub   events.EventProducer
	maxBalance int64
}

func NewTransactionService(
	walletRepo repository.WalletRepository,
	txnRepo repository.TransactionRepository,
	eventPub events.EventProducer,
	maxBalance int64,
) TransactionService {
	return &transactionService{
		walletRepo: walletRepo,
		txnRepo:    txnRepo,
		eventPub:   eventPub,
		maxBalance: maxBalance,
	}
}

func (s *transactionService) Deposit(ctx context.Context, userID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	if req.AmountCents <= 0 {
		return nil, domain.NewDomainError(domain.ErrNegativeAmount, 400)
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet.Status != domain.WalletActive {
		return nil, domain.NewDomainError(domain.ErrWalletFrozen, 403)
	}

	if wallet.BalanceCents+req.AmountCents > s.maxBalance {
		return nil, domain.NewDomainError(domain.ErrBalanceExceeded, 400)
	}

	txn := &domain.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		ToWalletID:     &wallet.ID,
		AmountCents:    req.AmountCents,
		Currency:       wallet.Currency,
		Type:           domain.TxnDeposit,
		Status:         domain.TxnSettled,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return nil, fmt.Errorf("failed to create deposit transaction: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, wallet, req.AmountCents); err != nil {
		s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnFailed, strPtr("balance update failed"))
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	txn.Status = domain.TxnSettled
	settledAt := time.Now().UTC()
	txn.SettledAt = &settledAt
	s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnSettled, nil)

	if err := s.eventPub.PublishTransactionSettled(ctx, txn); err != nil {
		log.Warn().Err(err).Msg("failed to publish deposit event")
	}

	return txn, nil
}

func (s *transactionService) Withdraw(ctx context.Context, userID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	if req.AmountCents <= 0 {
		return nil, domain.NewDomainError(domain.ErrNegativeAmount, 400)
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

	txn := &domain.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		FromWalletID:   &wallet.ID,
		AmountCents:    req.AmountCents,
		Currency:       wallet.Currency,
		Type:           domain.TxnWithdraw,
		Status:         domain.TxnSettled,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return nil, fmt.Errorf("failed to create withdraw transaction: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, wallet, -req.AmountCents); err != nil {
		s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnFailed, strPtr("balance update failed"))
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	settledAt := time.Now().UTC()
	txn.SettledAt = &settledAt
	s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnSettled, nil)

	if err := s.eventPub.PublishTransactionSettled(ctx, txn); err != nil {
		log.Warn().Err(err).Msg("failed to publish withdraw event")
	}

	return txn, nil
}

func (s *transactionService) Transfer(ctx context.Context, fromUserID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	return s.transferBetween(ctx, fromUserID, req, domain.TxnTransfer)
}

func (s *transactionService) Tip(ctx context.Context, fromUserID uuid.UUID, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	return s.transferBetween(ctx, fromUserID, req, domain.TxnTip)
}

func (s *transactionService) transferBetween(ctx context.Context, fromUserID uuid.UUID, req domain.CreateTransactionRequest, txnType domain.TransactionType) (*domain.Transaction, error) {
	if req.AmountCents <= 0 {
		return nil, domain.NewDomainError(domain.ErrNegativeAmount, 400)
	}
	if req.ToWalletID == nil {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "destination wallet required", 400)
	}

	fromWallet, err := s.walletRepo.GetByUserID(ctx, fromUserID)
	if err != nil {
		return nil, err
	}
	if fromWallet.Status != domain.WalletActive {
		return nil, domain.NewDomainError(domain.ErrWalletFrozen, 403)
	}
	if fromWallet.BalanceCents < req.AmountCents {
		return nil, domain.NewDomainError(domain.ErrInsufficientFunds, 402)
	}

	toWallet, err := s.walletRepo.GetByID(ctx, *req.ToWalletID)
	if err != nil {
		return nil, err
	}
	if toWallet.Status != domain.WalletActive {
		return nil, domain.NewDomainError(domain.ErrWalletFrozen, 403)
	}

	txn := &domain.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		FromWalletID:   &fromWallet.ID,
		ToWalletID:     &toWallet.ID,
		AmountCents:    req.AmountCents,
		Currency:       fromWallet.Currency,
		Type:           txnType,
		Status:         domain.TxnPending,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return nil, fmt.Errorf("failed to create transfer transaction: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, fromWallet, -req.AmountCents); err != nil {
		s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnFailed, strPtr("sender debit failed"))
		return nil, fmt.Errorf("failed to debit sender: %w", err)
	}

	if err := s.walletRepo.UpdateBalance(ctx, toWallet, req.AmountCents); err != nil {
		s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnFailed, strPtr("recipient credit failed"))
		return nil, fmt.Errorf("failed to credit recipient: %w", err)
	}

	settledAt := time.Now().UTC()
	txn.Status = domain.TxnSettled
	txn.SettledAt = &settledAt
	s.txnRepo.UpdateStatus(ctx, txn.ID, domain.TxnSettled, nil)

	if err := s.eventPub.PublishTransactionSettled(ctx, txn); err != nil {
		log.Warn().Err(err).Msg("failed to publish transfer event")
	}

	return txn, nil
}

func (s *transactionService) Get(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return s.txnRepo.GetByID(ctx, id)
}

func (s *transactionService) ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Transaction, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.txnRepo.ListByWallet(ctx, wallet.ID, cursor, limit)
}

func strPtr(s string) *string {
	return &s
}
