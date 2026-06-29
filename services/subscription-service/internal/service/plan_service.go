package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
	"github.com/elevatecompact/spark/services/subscription-service/internal/repository"
)

type PlanService interface {
	Create(ctx context.Context, req domain.CreatePlanRequest) (*domain.SubscriptionPlan, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error)
	List(ctx context.Context, creatorID *uuid.UUID, cursor time.Time, limit int) ([]*domain.SubscriptionPlan, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdatePlanRequest) (*domain.SubscriptionPlan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type planService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) PlanService {
	return &planService{repo: repo}
}

func (s *planService) Create(ctx context.Context, req domain.CreatePlanRequest) (*domain.SubscriptionPlan, error) {
	now := time.Now().UTC()
	plan := &domain.SubscriptionPlan{
		ID:            uuid.New(),
		CreatorID:     req.CreatorID,
		Name:          req.Name,
		PriceCents:    req.PriceCents,
		Currency:      req.Currency,
		BillingPeriod: req.BillingPeriod,
		Benefits:      req.Benefits,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.Create(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *planService) Get(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *planService) List(ctx context.Context, creatorID *uuid.UUID, cursor time.Time, limit int) ([]*domain.SubscriptionPlan, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.repo.List(ctx, creatorID, cursor, limit)
}

func (s *planService) Update(ctx context.Context, id uuid.UUID, req domain.UpdatePlanRequest) (*domain.SubscriptionPlan, error) {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.PriceCents != nil {
		plan.PriceCents = *req.PriceCents
	}
	if req.BillingPeriod != nil {
		plan.BillingPeriod = *req.BillingPeriod
	}
	if req.Benefits != nil {
		plan.Benefits = *req.Benefits
	}
	if req.IsActive != nil {
		plan.IsActive = *req.IsActive
	}
	if err := s.repo.Update(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *planService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}
