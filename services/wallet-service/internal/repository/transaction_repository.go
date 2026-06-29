package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID, cursor time.Time, limit int) ([]*domain.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus, failureReason *string) error
}

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{pool: pool}
}

func (r *transactionRepository) Create(ctx context.Context, txn *domain.Transaction) error {
	query := `INSERT INTO transactions (id, idempotency_key, from_wallet_id, to_wallet_id, amount_cents, currency, type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query,
		txn.ID, txn.IdempotencyKey, txn.FromWalletID, txn.ToWalletID,
		txn.AmountCents, txn.Currency, txn.Type, txn.Status, txn.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	query := `SELECT id, idempotency_key, from_wallet_id, to_wallet_id, amount_cents, currency, type, status, failure_reason, created_at, settled_at
		FROM transactions WHERE id = $1`
	return r.scanTransaction(ctx, query, id)
}

func (r *transactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	query := `SELECT id, idempotency_key, from_wallet_id, to_wallet_id, amount_cents, currency, type, status, failure_reason, created_at, settled_at
		FROM transactions WHERE idempotency_key = $1`
	return r.scanTransaction(ctx, query, key)
}

func (r *transactionRepository) ListByWallet(ctx context.Context, walletID uuid.UUID, cursor time.Time, limit int) ([]*domain.Transaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, idempotency_key, from_wallet_id, to_wallet_id, amount_cents, currency, type, status, failure_reason, created_at, settled_at
		FROM transactions
		WHERE (from_wallet_id = $1 OR to_wallet_id = $1) AND created_at < $2
		ORDER BY created_at DESC LIMIT $3`, walletID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var txns []*domain.Transaction
	for rows.Next() {
		txn, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}
	if txns == nil {
		txns = []*domain.Transaction{}
	}
	return txns, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus, failureReason *string) error {
	query := `UPDATE transactions SET status = $2, failure_reason = $3, settled_at = CASE WHEN $2 = 'settled' THEN NOW() ELSE settled_at END
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status, failureReason)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *transactionRepository) scanTransaction(ctx context.Context, query string, args ...interface{}) (*domain.Transaction, error) {
	txn := &domain.Transaction{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&txn.ID, &txn.IdempotencyKey, &txn.FromWalletID, &txn.ToWalletID,
		&txn.AmountCents, &txn.Currency, &txn.Type, &txn.Status,
		&txn.FailureReason, &txn.CreatedAt, &txn.SettledAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan transaction: %w", err)
	}
	return txn, nil
}

func (r *transactionRepository) scanRow(rows pgx.Rows) (*domain.Transaction, error) {
	txn := &domain.Transaction{}
	err := rows.Scan(
		&txn.ID, &txn.IdempotencyKey, &txn.FromWalletID, &txn.ToWalletID,
		&txn.AmountCents, &txn.Currency, &txn.Type, &txn.Status,
		&txn.FailureReason, &txn.CreatedAt, &txn.SettledAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan transaction row: %w", err)
	}
	return txn, nil
}
