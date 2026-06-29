package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	UpdateBalance(ctx context.Context, wallet *domain.Wallet, delta int64) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.WalletStatus) error
}

type walletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) WalletRepository {
	return &walletRepository{pool: pool}
}

func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `INSERT INTO wallets (id, user_id, balance_cents, currency, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query,
		wallet.ID, wallet.UserID, wallet.BalanceCents,
		wallet.Currency, wallet.Status, wallet.Version,
		wallet.CreatedAt, wallet.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance_cents, currency, status, version, created_at, updated_at
		FROM wallets WHERE id = $1`
	return r.scanWallet(ctx, query, id)
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance_cents, currency, status, version, created_at, updated_at
		FROM wallets WHERE user_id = $1`
	return r.scanWallet(ctx, query, userID)
}

func (r *walletRepository) UpdateBalance(ctx context.Context, wallet *domain.Wallet, delta int64) error {
	query := `UPDATE wallets SET balance_cents = balance_cents + $2, version = version + 1, updated_at = NOW()
		WHERE id = $1 AND version = $3
		RETURNING balance_cents, version`
	err := r.pool.QueryRow(ctx, query, wallet.ID, delta, wallet.Version).Scan(&wallet.BalanceCents, &wallet.Version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrWalletNotFound
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (r *walletRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.WalletStatus) error {
	tag, err := r.pool.Exec(ctx, `UPDATE wallets SET status = $2, updated_at = NOW() WHERE id = $1`, id, status)
	if err != nil {
		return fmt.Errorf("failed to update wallet status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrWalletNotFound
	}
	return nil
}

func (r *walletRepository) scanWallet(ctx context.Context, query string, args ...interface{}) (*domain.Wallet, error) {
	wallet := &domain.Wallet{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&wallet.ID, &wallet.UserID, &wallet.BalanceCents,
		&wallet.Currency, &wallet.Status, &wallet.Version,
		&wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to scan wallet: %w", err)
	}
	return wallet, nil
}
