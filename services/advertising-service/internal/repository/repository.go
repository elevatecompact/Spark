package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/advertising-service/internal/domain"
)

type AdvertisingRepository interface {
	CreateCampaign(ctx context.Context, c *domain.Campaign) error
	GetCampaign(ctx context.Context, id uuid.UUID) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, c *domain.Campaign) error
	ListCampaigns(ctx context.Context, advertiserID uuid.UUID) ([]domain.Campaign, error)

	CreateAdUnit(ctx context.Context, u *domain.AdUnit) error
	GetAdUnit(ctx context.Context, id uuid.UUID) (*domain.AdUnit, error)
	UpdateAdUnit(ctx context.Context, u *domain.AdUnit) error
	DeleteAdUnit(ctx context.Context, id uuid.UUID) error
	ListAdUnits(ctx context.Context, campaignID uuid.UUID) ([]domain.AdUnit, error)
	ApproveAdUnit(ctx context.Context, id uuid.UUID, approved bool, note string) error

	RecordImpression(ctx context.Context, i *domain.Impression) error
	RecordClick(ctx context.Context, c *domain.Click) error

	GetActiveAds(ctx context.Context, placementID string, limit int) ([]domain.AdUnit, error)
	GetCampaignPerformance(ctx context.Context, id uuid.UUID) (*domain.CampaignPerformance, error)
	GetRevenueStats(ctx context.Context) (*domain.RevenueStats, error)
}

type adRepo struct{ pool *pgxpool.Pool }

func NewAdvertisingRepository(pool *pgxpool.Pool) AdvertisingRepository { return &adRepo{pool} }

func (r *adRepo) CreateCampaign(ctx context.Context, c *domain.Campaign) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO campaigns (id, advertiser_id, name, budget_cents, spent_cents, daily_budget_cents, status, start_at, end_at, targeting, bid_strategy, created_at) VALUES ($1,$2,$3,$4,0,$5,$6,$7,$8,$9,$10,NOW())`,
		c.ID, c.AdvertiserID, c.Name, c.BudgetCents, c.DailyBudgetCents, c.Status, c.StartAt, c.EndAt, c.Targeting, c.BidStrategy)
	return err
}

func (r *adRepo) GetCampaign(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	c := &domain.Campaign{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, advertiser_id, name, budget_cents, spent_cents, daily_budget_cents, status, start_at, end_at, targeting, bid_strategy, created_at FROM campaigns WHERE id=$1`, id).
		Scan(&c.ID, &c.AdvertiserID, &c.Name, &c.BudgetCents, &c.SpentCents, &c.DailyBudgetCents, &c.Status, &c.StartAt, &c.EndAt, &c.Targeting, &c.BidStrategy, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return c, err
}

func (r *adRepo) UpdateCampaign(ctx context.Context, c *domain.Campaign) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE campaigns SET name=$2, budget_cents=$3, spent_cents=$4, daily_budget_cents=$5, status=$6, start_at=$7, end_at=$8, targeting=$9, bid_strategy=$10 WHERE id=$1`,
		c.ID, c.Name, c.BudgetCents, c.SpentCents, c.DailyBudgetCents, c.Status, c.StartAt, c.EndAt, c.Targeting, c.BidStrategy)
	return err
}

func (r *adRepo) ListCampaigns(ctx context.Context, advertiserID uuid.UUID) ([]domain.Campaign, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, advertiser_id, name, budget_cents, spent_cents, daily_budget_cents, status, start_at, end_at, targeting, bid_strategy, created_at FROM campaigns WHERE advertiser_id=$1 ORDER BY created_at DESC`, advertiserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var camps []domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		if err := rows.Scan(&c.ID, &c.AdvertiserID, &c.Name, &c.BudgetCents, &c.SpentCents, &c.DailyBudgetCents, &c.Status, &c.StartAt, &c.EndAt, &c.Targeting, &c.BidStrategy, &c.CreatedAt); err != nil {
			return nil, err
		}
		camps = append(camps, c)
	}
	if camps == nil {
		camps = []domain.Campaign{}
	}
	return camps, nil
}

func (r *adRepo) CreateAdUnit(ctx context.Context, u *domain.AdUnit) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO ad_units (id, campaign_id, type, format, creative_url, destination_url, width, height, duration_seconds, status, approval_note) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		u.ID, u.CampaignID, u.Type, u.Format, u.CreativeURL, u.DestinationURL, u.Width, u.Height, u.DurationSec, u.Status, u.ApprovalNote)
	return err
}

func (r *adRepo) GetAdUnit(ctx context.Context, id uuid.UUID) (*domain.AdUnit, error) {
	u := &domain.AdUnit{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, campaign_id, type, format, creative_url, destination_url, width, height, duration_seconds, status, approval_note FROM ad_units WHERE id=$1`, id).
		Scan(&u.ID, &u.CampaignID, &u.Type, &u.Format, &u.CreativeURL, &u.DestinationURL, &u.Width, &u.Height, &u.DurationSec, &u.Status, &u.ApprovalNote)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *adRepo) UpdateAdUnit(ctx context.Context, u *domain.AdUnit) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE ad_units SET type=$2, format=$3, creative_url=$4, destination_url=$5, width=$6, height=$7, duration_seconds=$8 WHERE id=$1`,
		u.ID, u.Type, u.Format, u.CreativeURL, u.DestinationURL, u.Width, u.Height, u.DurationSec)
	return err
}

func (r *adRepo) DeleteAdUnit(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM ad_units WHERE id=$1`, id)
	return err
}

func (r *adRepo) ListAdUnits(ctx context.Context, campaignID uuid.UUID) ([]domain.AdUnit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, campaign_id, type, format, creative_url, destination_url, width, height, duration_seconds, status, approval_note FROM ad_units WHERE campaign_id=$1`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []domain.AdUnit
	for rows.Next() {
		var u domain.AdUnit
		if err := rows.Scan(&u.ID, &u.CampaignID, &u.Type, &u.Format, &u.CreativeURL, &u.DestinationURL, &u.Width, &u.Height, &u.DurationSec, &u.Status, &u.ApprovalNote); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	if units == nil {
		units = []domain.AdUnit{}
	}
	return units, nil
}

func (r *adRepo) ApproveAdUnit(ctx context.Context, id uuid.UUID, approved bool, note string) error {
	status := domain.AdRejected
	if approved {
		status = domain.AdApproved
	}
	_, err := r.pool.Exec(ctx, `UPDATE ad_units SET status=$2, approval_note=$3 WHERE id=$1`, id, status, note)
	return err
}

func (r *adRepo) RecordImpression(ctx context.Context, i *domain.Impression) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO impressions (id, campaign_id, ad_unit_id, placement_id, user_id, cost_micro_cents, device_type, geo, served_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())`,
		i.ID, i.CampaignID, i.AdUnitID, i.PlacementID, i.UserID, i.CostMicroCents, i.DeviceType, i.Geo)
	if err == nil {
		r.pool.Exec(ctx, `UPDATE campaigns SET spent_cents = spent_cents + $2 WHERE id=$1`, i.CampaignID, i.CostMicroCents/100000)
	}
	return err
}

func (r *adRepo) RecordClick(ctx context.Context, c *domain.Click) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO clicks (id, impression_id, clicked_at) VALUES ($1,$2,NOW())`, c.ID, c.ImpressionID)
	return err
}

func (r *adRepo) GetActiveAds(ctx context.Context, placementID string, limit int) ([]domain.AdUnit, error) {
	if limit <= 0 || limit > 10 {
		limit = 5
	}
	rows, err := r.pool.Query(ctx,
		`SELECT au.id, au.campaign_id, au.type, au.format, au.creative_url, au.destination_url, au.width, au.height, au.duration_seconds, au.status, au.approval_note
		 FROM ad_units au JOIN campaigns c ON au.campaign_id = c.id
		 WHERE au.status='approved' AND c.status='active' AND c.spent_cents < c.budget_cents
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var units []domain.AdUnit
	for rows.Next() {
		var u domain.AdUnit
		if err := rows.Scan(&u.ID, &u.CampaignID, &u.Type, &u.Format, &u.CreativeURL, &u.DestinationURL, &u.Width, &u.Height, &u.DurationSec, &u.Status, &u.ApprovalNote); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	if units == nil {
		units = []domain.AdUnit{}
	}
	return units, nil
}

func (r *adRepo) GetCampaignPerformance(ctx context.Context, id uuid.UUID) (*domain.CampaignPerformance, error) {
	p := &domain.CampaignPerformance{}
	err := r.pool.QueryRow(ctx,
		`SELECT c.id,
		        (SELECT COUNT(*) FROM impressions WHERE campaign_id=c.id) as imps,
		        (SELECT COUNT(*) FROM impressions i JOIN clicks cl ON i.id=cl.impression_id WHERE i.campaign_id=c.id) as clks,
		        c.spent_cents, c.budget_cents
		 FROM campaigns c WHERE c.id=$1`, id).
		Scan(&p.CampaignID, &p.Impressions, &p.Clicks, &p.SpentCents, &p.BudgetCents)
	if err != nil {
		return nil, err
	}
	if p.Impressions > 0 {
		p.CTR = float64(p.Clicks) / float64(p.Impressions)
	}
	return p, nil
}

func (r *adRepo) GetRevenueStats(ctx context.Context) (*domain.RevenueStats, error) {
	s := &domain.RevenueStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(spent_cents),0) FROM campaigns`).Scan(&s.TotalRevenueCents)
	if err != nil {
		return nil, err
	}
	s.PlatformShareCents = int64(float64(s.TotalRevenueCents) * 0.3)
	s.CreatorShareCents = int64(float64(s.TotalRevenueCents) * 0.7)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM campaigns WHERE status='active'`).Scan(&s.ActiveCampaigns)
	return s, nil
}
