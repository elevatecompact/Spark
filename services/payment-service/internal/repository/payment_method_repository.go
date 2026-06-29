package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
)

type PaymentMethodRepository interface {
	Create(ctx context.Context, method *domain.PaymentMethod) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentMethod, error)
	GetByFingerprint(ctx context.Context, fingerprint string) (*domain.PaymentMethod, error)
	SetDefault(ctx context.Context, id, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type paymentMethodRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentMethodRepository(pool *pgxpool.Pool) PaymentMethodRepository {
	return &paymentMethodRepository{pool: pool}
}

func (r *paymentMethodRepository) Create(ctx context.Context, method *domain.PaymentMethod) error {
	query := `INSERT INTO payment_methods (id, user_id, external_id, processor, type, fingerprint, last4, exp_month, exp_year, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, method.ID, method.UserID, method.ExternalID, method.Processor, method.Type, method.Fingerprint, method.Last4, method.ExpMonth, method.ExpYear, method.IsDefault, method.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment method: %w", err)
	}
	return nil
}

func (r *paymentMethodRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error) {
	query := `SELECT id, user_id, external_id, processor, type, fingerprint, last4, exp_month, exp_year, is_default, created_at
		FROM payment_methods WHERE id = $1`
	m := &domain.PaymentMethod{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.UserID, &m.ExternalID, &m.Processor, &m.Type, &m.Fingerprint, &m.Last4, &m.ExpMonth, &m.ExpYear, &m.IsDefault, &m.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrMethodNotFound
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	return m, nil
}

func (r *paymentMethodRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PaymentMethod, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, external_id, processor, type, fingerprint, last4, exp_month, exp_year, is_default, created_at
		FROM payment_methods WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}
	defer rows.Close()

	var methods []*domain.PaymentMethod
	for rows.Next() {
		m := &domain.PaymentMethod{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.ExternalID, &m.Processor, &m.Type, &m.Fingerprint, &m.Last4, &m.ExpMonth, &m.ExpYear, &m.IsDefault, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}
		methods = append(methods, m)
	}
	if methods == nil {
		methods = []*domain.PaymentMethod{}
	}
	return methods, nil
}

func (r *paymentMethodRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*domain.PaymentMethod, error) {
	query := `SELECT id, user_id, external_id, processor, type, fingerprint, last4, exp_month, exp_year, is_default, created_at
		FROM payment_methods WHERE fingerprint = $1`
	m := &domain.PaymentMethod{}
	err := r.pool.QueryRow(ctx, query, fingerprint).Scan(&m.ID, &m.UserID, &m.ExternalID, &m.Processor, &m.Type, &m.Fingerprint, &m.Last4, &m.ExpMonth, &m.ExpYear, &m.IsDefault, &m.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get method by fingerprint: %w", err)
	}
	return m, nil
}

func (r *paymentMethodRepository) SetDefault(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE payment_methods SET is_default=false WHERE user_id=$1`, userID)
	if err != nil {
		return fmt.Errorf("failed to clear default methods: %w", err)
	}
	tag, err := r.pool.Exec(ctx, `UPDATE payment_methods SET is_default=true WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to set default method: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMethodNotFound
	}
	return nil
}

func (r *paymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM payment_methods WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMethodNotFound
	}
	return nil
}
