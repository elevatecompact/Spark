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

type PayoutRepository interface {
	Create(ctx context.Context, payout *domain.Payout) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payout, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID, cursor time.Time, limit int) ([]*domain.Payout, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PayoutStatus, externalRef *string) error
}

type payoutRepository struct {
	pool *pgxpool.Pool
}

func NewPayoutRepository(pool *pgxpool.Pool) PayoutRepository {
	return &payoutRepository{pool: pool}
}

func (r *payoutRepository) Create(ctx context.Context, payout *domain.Payout) error {
	query := `INSERT INTO payouts (id, wallet_id, amount_cents, currency, method, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		payout.ID, payout.WalletID, payout.AmountCents,
		payout.Currency, payout.Method, payout.Status, payout.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payout: %w", err)
	}
	return nil
}

func (r *payoutRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payout, error) {
	query := `SELECT id, wallet_id, amount_cents, currency, method, status, external_ref, created_at, completed_at
		FROM payouts WHERE id = $1`
	payout := &domain.Payout{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&payout.ID, &payout.WalletID, &payout.AmountCents,
		&payout.Currency, &payout.Method, &payout.Status,
		&payout.ExternalRef, &payout.CreatedAt, &payout.CompletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan payout: %w", err)
	}
	return payout, nil
}

func (r *payoutRepository) ListByWallet(ctx context.Context, walletID uuid.UUID, cursor time.Time, limit int) ([]*domain.Payout, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, wallet_id, amount_cents, currency, method, status, external_ref, created_at, completed_at
		FROM payouts WHERE wallet_id = $1 AND created_at < $2
		ORDER BY created_at DESC LIMIT $3`, walletID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	defer rows.Close()

	var payouts []*domain.Payout
	for rows.Next() {
		p := &domain.Payout{}
		if err := rows.Scan(&p.ID, &p.WalletID, &p.AmountCents, &p.Currency, &p.Method, &p.Status, &p.ExternalRef, &p.CreatedAt, &p.CompletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan payout row: %w", err)
		}
		payouts = append(payouts, p)
	}
	if payouts == nil {
		payouts = []*domain.Payout{}
	}
	return payouts, nil
}

func (r *payoutRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PayoutStatus, externalRef *string) error {
	query := `UPDATE payouts SET status = $2, external_ref = $3, completed_at = CASE WHEN $2 = 'completed' THEN NOW() ELSE completed_at END
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status, externalRef)
	if err != nil {
		return fmt.Errorf("failed to update payout status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
