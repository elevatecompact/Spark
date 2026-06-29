package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.SubscriptionPlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error)
	List(ctx context.Context, creatorID *uuid.UUID, cursor time.Time, limit int) ([]*domain.SubscriptionPlan, error)
	Update(ctx context.Context, plan *domain.SubscriptionPlan) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type planRepository struct {
	pool *pgxpool.Pool
}

func NewPlanRepository(pool *pgxpool.Pool) PlanRepository {
	return &planRepository{pool: pool}
}

func (r *planRepository) Create(ctx context.Context, plan *domain.SubscriptionPlan) error {
	query := `INSERT INTO subscription_plans (id, creator_id, name, price_cents, currency, billing_period, benefits, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, query, plan.ID, plan.CreatorID, plan.Name, plan.PriceCents, plan.Currency, plan.BillingPeriod, plan.Benefits, plan.IsActive, plan.CreatedAt, plan.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}
	return nil
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	query := `SELECT id, creator_id, name, price_cents, currency, billing_period, benefits, is_active, created_at, updated_at
		FROM subscription_plans WHERE id = $1`
	plan := &domain.SubscriptionPlan{}
	var benefits []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&plan.ID, &plan.CreatorID, &plan.Name, &plan.PriceCents, &plan.Currency, &plan.BillingPeriod, &benefits, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPlanNotFound
		}
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}
	plan.Benefits = json.RawMessage(benefits)
	return plan, nil
}

func (r *planRepository) List(ctx context.Context, creatorID *uuid.UUID, cursor time.Time, limit int) ([]*domain.SubscriptionPlan, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var rows pgx.Rows
	var err error
	if creatorID != nil {
		rows, err = r.pool.Query(ctx, `
			SELECT id, creator_id, name, price_cents, currency, billing_period, benefits, is_active, created_at, updated_at
			FROM subscription_plans WHERE creator_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`, *creatorID, cursor, limit)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, creator_id, name, price_cents, currency, billing_period, benefits, is_active, created_at, updated_at
			FROM subscription_plans WHERE created_at < $1 ORDER BY created_at DESC LIMIT $2`, cursor, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	defer rows.Close()

	var plans []*domain.SubscriptionPlan
	for rows.Next() {
		p := &domain.SubscriptionPlan{}
		var benefits []byte
		if err := rows.Scan(&p.ID, &p.CreatorID, &p.Name, &p.PriceCents, &p.Currency, &p.BillingPeriod, &benefits, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		p.Benefits = json.RawMessage(benefits)
		plans = append(plans, p)
	}
	if plans == nil {
		plans = []*domain.SubscriptionPlan{}
	}
	return plans, nil
}

func (r *planRepository) Update(ctx context.Context, plan *domain.SubscriptionPlan) error {
	query := `UPDATE subscription_plans SET name=$2, price_cents=$3, currency=$4, billing_period=$5, benefits=$6, is_active=$7, updated_at=NOW() WHERE id=$1`
	tag, err := r.pool.Exec(ctx, query, plan.ID, plan.Name, plan.PriceCents, plan.Currency, plan.BillingPeriod, plan.Benefits, plan.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlanNotFound
	}
	return nil
}

func (r *planRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `UPDATE subscription_plans SET is_active=false, updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPlanNotFound
	}
	return nil
}
