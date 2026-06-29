package domain

import (
	"time"

	"github.com/google/uuid"
)

type GiftCategory string

const (
	GiftCategoryEmote  GiftCategory = "emote"
	GiftCategoryBadge  GiftCategory = "badge"
	GiftCategoryEffect GiftCategory = "effect"
	GiftCategorySub    GiftCategory = "sub"
)

type GiftStatus string

const (
	GiftPending   GiftStatus = "pending"
	GiftCompleted GiftStatus = "completed"
	GiftRefunded  GiftStatus = "refunded"
)

type GiftItem struct {
	ID         uuid.UUID    `json:"id"`
	Name       string       `json:"name"`
	PriceCents int64        `json:"price_cents"`
	ImageURL   string       `json:"image_url"`
	Category   GiftCategory `json:"category"`
	IsActive   bool         `json:"is_active"`
	SortOrder  int          `json:"sort_order"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type Gift struct {
	ID          uuid.UUID  `json:"id"`
	SenderID    uuid.UUID  `json:"sender_id"`
	RecipientID uuid.UUID  `json:"recipient_id"`
	GiftItemID  *uuid.UUID `json:"gift_item_id,omitempty"`
	AmountCents int64      `json:"amount_cents"`
	Message     string     `json:"message"`
	CampaignID  *uuid.UUID `json:"campaign_id,omitempty"`
	IsAnonymous bool       `json:"is_anonymous"`
	Status      GiftStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type GiftCard struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"`
	PurchaserID  uuid.UUID  `json:"purchaser_id"`
	BalanceCents int64      `json:"balance_cents"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RedeemedAt   *time.Time `json:"redeemed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type GiftCampaign struct {
	ID             uuid.UUID `json:"id"`
	CreatorID      uuid.UUID `json:"creator_id"`
	MatchRatio     float64   `json:"match_ratio"`
	MaxMatchCents  int64     `json:"max_match_cents"`
	TotalMatched   int64     `json:"total_matched"`
	StartAt        time.Time `json:"start_at"`
	EndAt          time.Time `json:"end_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	UserID    uuid.UUID `json:"user_id"`
	GiftCount int       `json:"gift_count"`
	TotalCents int64    `json:"total_cents"`
	Rank      int       `json:"rank"`
}

type SendGiftRequest struct {
	RecipientID uuid.UUID `json:"recipient_id"`
	GiftItemID  *uuid.UUID `json:"gift_item_id,omitempty"`
	AmountCents int64     `json:"amount_cents"`
	Message     string    `json:"message"`
	CampaignID  *uuid.UUID `json:"campaign_id,omitempty"`
	IsAnonymous bool      `json:"is_anonymous"`
}

type SendBatchGiftRequest struct {
	Gifts []SendGiftRequest `json:"gifts"`
}

type SendSubscriptionGiftRequest struct {
	RecipientID uuid.UUID `json:"recipient_id"`
	PlanID      uuid.UUID `json:"plan_id"`
	Message     string    `json:"message"`
	IsAnonymous bool      `json:"is_anonymous"`
}

type PurchaseGiftCardRequest struct {
	AmountCents int64 `json:"amount_cents"`
}

type RedeemGiftCardRequest struct {
	Code string `json:"code"`
}

type CreateCampaignRequest struct {
	MatchRatio    float64   `json:"match_ratio"`
	MaxMatchCents int64     `json:"max_match_cents"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
}
