package domain

import (
	"time"

	"github.com/google/uuid"
)

type CampaignStatus string
const (
	CampDraft   CampaignStatus = "draft"
	CampActive  CampaignStatus = "active"
	CampPaused  CampaignStatus = "paused"
	CampEnded   CampaignStatus = "ended"
)

type BidStrategy string
const (
	BidCPM BidStrategy = "cpm"
	BidCPC BidStrategy = "cpc"
)

type AdType string
const (
	AdPreRoll    AdType = "preroll"
	AdMidRoll    AdType = "midroll"
	AdDisplay    AdType = "display"
	AdSponsored AdType = "sponsored"
)

type AdFormat string
const (
	FmtVideo AdFormat = "video"
	FmtImage AdFormat = "image"
	FmtText  AdFormat = "text"
)

type AdStatus string
const (
	AdPending  AdStatus = "pending"
	AdApproved AdStatus = "approved"
	AdRejected AdStatus = "rejected"
)

type Campaign struct {
	ID              uuid.UUID      `json:"id"`
	AdvertiserID    uuid.UUID      `json:"advertiser_id"`
	Name            string         `json:"name"`
	BudgetCents     int64          `json:"budget_cents"`
	SpentCents      int64          `json:"spent_cents"`
	DailyBudgetCents int64         `json:"daily_budget_cents"`
	Status          CampaignStatus `json:"status"`
	StartAt         *time.Time     `json:"start_at,omitempty"`
	EndAt           *time.Time     `json:"end_at,omitempty"`
	Targeting       map[string]interface{} `json:"targeting,omitempty"`
	BidStrategy     BidStrategy    `json:"bid_strategy"`
	CreatedAt       time.Time      `json:"created_at"`
}

type AdUnit struct {
	ID             uuid.UUID `json:"id"`
	CampaignID     uuid.UUID `json:"campaign_id"`
	Type           AdType    `json:"type"`
	Format         AdFormat  `json:"format"`
	CreativeURL    string    `json:"creative_url"`
	DestinationURL string    `json:"destination_url"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	DurationSec    int       `json:"duration_seconds,omitempty"`
	Status         AdStatus  `json:"status"`
	ApprovalNote   string    `json:"approval_note,omitempty"`
}

type Impression struct {
	ID            uuid.UUID `json:"id"`
	CampaignID    uuid.UUID `json:"campaign_id"`
	AdUnitID      uuid.UUID `json:"ad_unit_id"`
	PlacementID   string    `json:"placement_id"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	CostMicroCents int64    `json:"cost_micro_cents"`
	DeviceType    string    `json:"device_type"`
	Geo           string    `json:"geo"`
	ServedAt      time.Time `json:"served_at"`
}

type Click struct {
	ID           uuid.UUID `json:"id"`
	ImpressionID uuid.UUID `json:"impression_id"`
	ClickedAt    time.Time `json:"clicked_at"`
}

type AdInventory struct {
	PlacementID        string    `json:"placement_id"`
	ContentType        string    `json:"content_type"`
	AvailableFrom      time.Time `json:"available_from"`
	AvailableTo        time.Time `json:"available_to"`
	FloorPriceMicroCents int64   `json:"floor_price_micro_cents"`
	IsActive           bool      `json:"is_active"`
}

type CampaignPerformance struct {
	CampaignID      uuid.UUID `json:"campaign_id"`
	Impressions     int64     `json:"impressions"`
	Clicks          int64     `json:"clicks"`
	CTR             float64   `json:"ctr"`
	SpentCents      int64     `json:"spent_cents"`
	BudgetCents     int64     `json:"budget_cents"`
}

type CreatorRevenue struct {
	CreatorID  uuid.UUID `json:"creator_id"`
	RevenueCents int64   `json:"revenue_cents"`
	Period     string    `json:"period"`
}

type RevenueStats struct {
	TotalRevenueCents int64 `json:"total_revenue_cents"`
	PlatformShareCents int64 `json:"platform_share_cents"`
	CreatorShareCents  int64 `json:"creator_share_cents"`
	ActiveCampaigns    int   `json:"active_campaigns"`
}
