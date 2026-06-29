package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
)

type InvoiceRepository interface {
	Create(ctx context.Context, inv *domain.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)
	ListBySubscription(ctx context.Context, subID uuid.UUID) ([]*domain.Invoice, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.InvoiceStatus) error
}

type invoiceRepository struct {
	pool *pgxpool.Pool
}

func NewInvoiceRepository(pool *pgxpool.Pool) InvoiceRepository {
	return &invoiceRepository{pool: pool}
}

func (r *invoiceRepository) Create(ctx context.Context, inv *domain.Invoice) error {
	query := `INSERT INTO invoices (id, subscription_id, amount_cents, currency, status, period_start, period_end, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, inv.ID, inv.SubscriptionID, inv.AmountCents, inv.Currency, inv.Status, inv.PeriodStart, inv.PeriodEnd, inv.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	query := `SELECT id, subscription_id, amount_cents, currency, status, paid_at, period_start, period_end, created_at
		FROM invoices WHERE id = $1`
	inv := &domain.Invoice{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&inv.ID, &inv.SubscriptionID, &inv.AmountCents, &inv.Currency, &inv.Status, &inv.PaidAt, &inv.PeriodStart, &inv.PeriodEnd, &inv.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}
	return inv, nil
}

func (r *invoiceRepository) ListBySubscription(ctx context.Context, subID uuid.UUID) ([]*domain.Invoice, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, subscription_id, amount_cents, currency, status, paid_at, period_start, period_end, created_at
		FROM invoices WHERE subscription_id = $1 ORDER BY created_at DESC`, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		inv := &domain.Invoice{}
		if err := rows.Scan(&inv.ID, &inv.SubscriptionID, &inv.AmountCents, &inv.Currency, &inv.Status, &inv.PaidAt, &inv.PeriodStart, &inv.PeriodEnd, &inv.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		invoices = append(invoices, inv)
	}
	if invoices == nil {
		invoices = []*domain.Invoice{}
	}
	return invoices, nil
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.InvoiceStatus) error {
	query := `UPDATE invoices SET status = $2, paid_at = CASE WHEN $2 = 'paid' THEN NOW() ELSE paid_at END WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrInvoiceNotFound
	}
	return nil
}
