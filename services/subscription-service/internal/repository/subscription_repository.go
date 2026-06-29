package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	GetByUserAndPlan(ctx context.Context, userID, planID uuid.UUID) (*domain.Subscription, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error)
	ListActive(ctx context.Context) ([]*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	CountActiveByUser(ctx context.Context, userID uuid.UUID) (int, error)
}

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `INSERT INTO subscriptions (id, user_id, plan_id, status, current_period_start, current_period_end, grace_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, sub.ID, sub.UserID, sub.PlanID, sub.Status, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.GracePeriodEnd, sub.CreatedAt, sub.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, user_id, plan_id, status, current_period_start, current_period_end, cancelled_at, grace_period_end, created_at, updated_at
		FROM subscriptions WHERE id = $1`
	sub := &domain.Subscription{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelledAt, &sub.GracePeriodEnd, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepository) GetByUserAndPlan(ctx context.Context, userID, planID uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, user_id, plan_id, status, current_period_start, current_period_end, cancelled_at, grace_period_end, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 AND plan_id = $2`
	sub := &domain.Subscription{}
	err := r.pool.QueryRow(ctx, query, userID, planID).Scan(&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelledAt, &sub.GracePeriodEnd, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription by user and plan: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, plan_id, status, current_period_start, current_period_end, cancelled_at, grace_period_end, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		sub := &domain.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelledAt, &sub.GracePeriodEnd, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	return subs, nil
}

func (r *subscriptionRepository) ListActive(ctx context.Context) ([]*domain.Subscription, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, plan_id, status, current_period_start, current_period_end, cancelled_at, grace_period_end, created_at, updated_at
		FROM subscriptions WHERE status = 'active' OR status = 'grace_period'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		sub := &domain.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelledAt, &sub.GracePeriodEnd, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	return subs, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	query := `UPDATE subscriptions SET plan_id=$2, status=$3, current_period_start=$4, current_period_end=$5, cancelled_at=$6, grace_period_end=$7, updated_at=NOW() WHERE id=$1`
	tag, err := r.pool.Exec(ctx, query, sub.ID, sub.PlanID, sub.Status, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CancelledAt, sub.GracePeriodEnd)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}
	return nil
}

func (r *subscriptionRepository) CountActiveByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscriptions WHERE user_id = $1 AND status IN ('active', 'grace_period')`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active subscriptions: %w", err)
	}
	return count, nil
}
