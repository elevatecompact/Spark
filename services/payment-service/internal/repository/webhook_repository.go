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

type WebhookRepository interface {
	Create(ctx context.Context, event *domain.WebhookEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WebhookEvent, error)
	GetByExternalEventID(ctx context.Context, processor domain.PaymentProcessor, externalID string) (*domain.WebhookEvent, error)
	ListByStatus(ctx context.Context, status domain.WebhookStatus, limit int) ([]*domain.WebhookEvent, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.WebhookStatus) error
}

type webhookRepository struct {
	pool *pgxpool.Pool
}

func NewWebhookRepository(pool *pgxpool.Pool) WebhookRepository {
	return &webhookRepository{pool: pool}
}

func (r *webhookRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := `INSERT INTO webhook_events (id, processor, external_event_id, type, body, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, event.ID, event.Processor, event.ExternalEventID, event.Type, event.Body, event.Status, event.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}
	return nil
}

func (r *webhookRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WebhookEvent, error) {
	query := `SELECT id, processor, external_event_id, type, body, status, created_at, processed_at
		FROM webhook_events WHERE id = $1`
	e := &domain.WebhookEvent{}
	var body []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&e.ID, &e.Processor, &e.ExternalEventID, &e.Type, &body, &e.Status, &e.CreatedAt, &e.ProcessedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrWebhookNotFound
		}
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}
	e.Body = json.RawMessage(body)
	return e, nil
}

func (r *webhookRepository) GetByExternalEventID(ctx context.Context, processor domain.PaymentProcessor, externalID string) (*domain.WebhookEvent, error) {
	query := `SELECT id, processor, external_event_id, type, body, status, created_at, processed_at
		FROM webhook_events WHERE processor = $1 AND external_event_id = $2`
	e := &domain.WebhookEvent{}
	var body []byte
	err := r.pool.QueryRow(ctx, query, processor, externalID).Scan(&e.ID, &e.Processor, &e.ExternalEventID, &e.Type, &body, &e.Status, &e.CreatedAt, &e.ProcessedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get webhook by external id: %w", err)
	}
	e.Body = json.RawMessage(body)
	return e, nil
}

func (r *webhookRepository) ListByStatus(ctx context.Context, status domain.WebhookStatus, limit int) ([]*domain.WebhookEvent, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id, processor, external_event_id, type, body, status, created_at, processed_at
		FROM webhook_events WHERE status = $1 ORDER BY created_at DESC LIMIT $2`, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook events: %w", err)
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		e := &domain.WebhookEvent{}
		var body []byte
		if err := rows.Scan(&e.ID, &e.Processor, &e.ExternalEventID, &e.Type, &body, &e.Status, &e.CreatedAt, &e.ProcessedAt); err != nil {
			return nil, fmt.Errorf("failed to scan webhook event: %w", err)
		}
		e.Body = json.RawMessage(body)
		events = append(events, e)
	}
	if events == nil {
		events = []*domain.WebhookEvent{}
	}
	return events, nil
}

func (r *webhookRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.WebhookStatus) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE webhook_events SET status=$2, processed_at=CASE WHEN $2='processed' THEN $3 ELSE processed_at END WHERE id=$1`, id, status, now)
	if err != nil {
		return fmt.Errorf("failed to update webhook status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrWebhookNotFound
	}
	return nil
}
