package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
)

type PaymentIntentRepository interface {
	Create(ctx context.Context, intent *domain.PaymentIntent) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentIntent, error)
	GetByExternalID(ctx context.Context, externalID string) (*domain.PaymentIntent, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.PaymentIntent, error)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.PaymentIntent, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IntentStatus) error
	Update(ctx context.Context, intent *domain.PaymentIntent) error
}

type paymentIntentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentIntentRepository(pool *pgxpool.Pool) PaymentIntentRepository {
	return &paymentIntentRepository{pool: pool}
}

func (r *paymentIntentRepository) Create(ctx context.Context, intent *domain.PaymentIntent) error {
	query := `INSERT INTO payment_intents (id, user_id, external_id, processor, amount_cents, currency, status, idempotency_key, metadata, payment_method_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.pool.Exec(ctx, query, intent.ID, intent.UserID, intent.ExternalID, intent.Processor, intent.AmountCents, intent.Currency, intent.Status, intent.IdempotencyKey, intent.Metadata, intent.PaymentMethodID, intent.CreatedAt, intent.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment intent: %w", err)
	}
	return nil
}

func (r *paymentIntentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentIntent, error) {
	query := `SELECT id, user_id, external_id, processor, amount_cents, currency, status, idempotency_key, metadata, payment_method_id, created_at, updated_at
		FROM payment_intents WHERE id = $1`
	intent := &domain.PaymentIntent{}
	var metadata []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&intent.ID, &intent.UserID, &intent.ExternalID, &intent.Processor, &intent.AmountCents, &intent.Currency, &intent.Status, &intent.IdempotencyKey, &metadata, &intent.PaymentMethodID, &intent.CreatedAt, &intent.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrIntentNotFound
		}
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	intent.Metadata = json.RawMessage(metadata)
	return intent, nil
}

func (r *paymentIntentRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.PaymentIntent, error) {
	query := `SELECT id, user_id, external_id, processor, amount_cents, currency, status, idempotency_key, metadata, payment_method_id, created_at, updated_at
		FROM payment_intents WHERE external_id = $1`
	intent := &domain.PaymentIntent{}
	var metadata []byte
	err := r.pool.QueryRow(ctx, query, externalID).Scan(&intent.ID, &intent.UserID, &intent.ExternalID, &intent.Processor, &intent.AmountCents, &intent.Currency, &intent.Status, &intent.IdempotencyKey, &metadata, &intent.PaymentMethodID, &intent.CreatedAt, &intent.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrIntentNotFound
		}
		return nil, fmt.Errorf("failed to get intent by external id: %w", err)
	}
	intent.Metadata = json.RawMessage(metadata)
	return intent, nil
}

func (r *paymentIntentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.PaymentIntent, error) {
	query := `SELECT id, user_id, external_id, processor, amount_cents, currency, status, idempotency_key, metadata, payment_method_id, created_at, updated_at
		FROM payment_intents WHERE idempotency_key = $1`
	intent := &domain.PaymentIntent{}
	var metadata []byte
	err := r.pool.QueryRow(ctx, query, key).Scan(&intent.ID, &intent.UserID, &intent.ExternalID, &intent.Processor, &intent.AmountCents, &intent.Currency, &intent.Status, &intent.IdempotencyKey, &metadata, &intent.PaymentMethodID, &intent.CreatedAt, &intent.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get intent by idempotency key: %w", err)
	}
	intent.Metadata = json.RawMessage(metadata)
	return intent, nil
}

func (r *paymentIntentRepository) ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.PaymentIntent, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, external_id, processor, amount_cents, currency, status, idempotency_key, metadata, payment_method_id, created_at, updated_at
		FROM payment_intents WHERE user_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`, userID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list intents: %w", err)
	}
	defer rows.Close()

	var intents []*domain.PaymentIntent
	for rows.Next() {
		intent := &domain.PaymentIntent{}
		var metadata []byte
		if err := rows.Scan(&intent.ID, &intent.UserID, &intent.ExternalID, &intent.Processor, &intent.AmountCents, &intent.Currency, &intent.Status, &intent.IdempotencyKey, &metadata, &intent.PaymentMethodID, &intent.CreatedAt, &intent.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan intent: %w", err)
		}
		intent.Metadata = json.RawMessage(metadata)
		intents = append(intents, intent)
	}
	if intents == nil {
		intents = []*domain.PaymentIntent{}
	}
	return intents, nil
}

func (r *paymentIntentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IntentStatus) error {
	tag, err := r.pool.Exec(ctx, `UPDATE payment_intents SET status=$2, updated_at=NOW() WHERE id=$1`, id, status)
	if err != nil {
		return fmt.Errorf("failed to update intent status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrIntentNotFound
	}
	return nil
}

func (r *paymentIntentRepository) Update(ctx context.Context, intent *domain.PaymentIntent) error {
	query := `UPDATE payment_intents SET external_id=$2, status=$3, payment_method_id=$4, updated_at=NOW() WHERE id=$1`
	tag, err := r.pool.Exec(ctx, query, intent.ID, intent.ExternalID, intent.Status, intent.PaymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to update intent: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrIntentNotFound
	}
	return nil
}
