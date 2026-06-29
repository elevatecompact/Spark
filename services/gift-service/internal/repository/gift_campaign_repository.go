package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
)

type GiftCampaignRepository interface {
	Create(ctx context.Context, campaign *domain.GiftCampaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftCampaign, error)
	ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*domain.GiftCampaign, error)
	ListActive(ctx context.Context) ([]*domain.GiftCampaign, error)
	AddMatchAmount(ctx context.Context, id uuid.UUID, amountCents int64) error
	CountByCreatorSince(ctx context.Context, creatorID uuid.UUID, since time.Time) (int, error)
}

type giftCampaignRepository struct {
	pool *pgxpool.Pool
}

func NewGiftCampaignRepository(pool *pgxpool.Pool) GiftCampaignRepository {
	return &giftCampaignRepository{pool: pool}
}

func (r *giftCampaignRepository) Create(ctx context.Context, campaign *domain.GiftCampaign) error {
	query := `INSERT INTO gift_campaigns (id, creator_id, match_ratio, max_match_cents, total_matched, start_at, end_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, campaign.ID, campaign.CreatorID, campaign.MatchRatio, campaign.MaxMatchCents, campaign.TotalMatched, campaign.StartAt, campaign.EndAt, campaign.CreatedAt, campaign.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	return nil
}

func (r *giftCampaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftCampaign, error) {
	query := `SELECT id, creator_id, match_ratio, max_match_cents, total_matched, start_at, end_at, created_at, updated_at
		FROM gift_campaigns WHERE id = $1`
	c := &domain.GiftCampaign{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.CreatorID, &c.MatchRatio, &c.MaxMatchCents, &c.TotalMatched, &c.StartAt, &c.EndAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}
	return c, nil
}

func (r *giftCampaignRepository) ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*domain.GiftCampaign, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, creator_id, match_ratio, max_match_cents, total_matched, start_at, end_at, created_at, updated_at
		FROM gift_campaigns WHERE creator_id = $1 ORDER BY created_at DESC`, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}
	defer rows.Close()
	return scanCampaigns(rows)
}

func (r *giftCampaignRepository) ListActive(ctx context.Context) ([]*domain.GiftCampaign, error) {
	now := time.Now().UTC()
	rows, err := r.pool.Query(ctx, `SELECT id, creator_id, match_ratio, max_match_cents, total_matched, start_at, end_at, created_at, updated_at
		FROM gift_campaigns WHERE start_at <= $1 AND end_at >= $1 AND total_matched < max_match_cents ORDER BY created_at DESC`, now)
	if err != nil {
		return nil, fmt.Errorf("failed to list active campaigns: %w", err)
	}
	defer rows.Close()
	return scanCampaigns(rows)
}

func (r *giftCampaignRepository) AddMatchAmount(ctx context.Context, id uuid.UUID, amountCents int64) error {
	tag, err := r.pool.Exec(ctx, `UPDATE gift_campaigns SET total_matched = total_matched + $2, updated_at = NOW() WHERE id = $1 AND total_matched + $2 <= max_match_cents`, id, amountCents)
	if err != nil {
		return fmt.Errorf("failed to add match amount: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCampaignBudgetExhausted
	}
	return nil
}

func (r *giftCampaignRepository) CountByCreatorSince(ctx context.Context, creatorID uuid.UUID, since time.Time) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM gift_campaigns WHERE creator_id = $1 AND created_at >= $2`, creatorID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count campaigns: %w", err)
	}
	return count, nil
}

func scanCampaigns(rows pgx.Rows) ([]*domain.GiftCampaign, error) {
	var campaigns []*domain.GiftCampaign
	for rows.Next() {
		c := &domain.GiftCampaign{}
		if err := rows.Scan(&c.ID, &c.CreatorID, &c.MatchRatio, &c.MaxMatchCents, &c.TotalMatched, &c.StartAt, &c.EndAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, c)
	}
	if campaigns == nil {
		campaigns = []*domain.GiftCampaign{}
	}
	return campaigns, nil
}
