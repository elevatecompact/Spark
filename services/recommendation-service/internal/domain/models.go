package domain

import (
	"time"

	"github.com/google/uuid"
)

type FeedType string

const (
	FeedHome     FeedType = "home"
	FeedTrending FeedType = "trending"
	FeedUpNext   FeedType = "up_next"
	FeedSimilar  FeedType = "similar"
	FeedCreator  FeedType = "creator"
)

type InteractionType string

const (
	InteractionClick      InteractionType = "click"
	InteractionWatch      InteractionType = "watch"
	InteractionRate       InteractionType = "rate"
	InteractionSubscribe  InteractionType = "subscribe"
	InteractionDismiss    InteractionType = "dismiss"
)

type UserEmbedding struct {
	UserID       uuid.UUID `json:"user_id"`
	Embedding    []float64 `json:"embedding"`
	ModelVersion string    `json:"model_version"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ContentEmbedding struct {
	ContentID    uuid.UUID `json:"content_id"`
	Embedding    []float64 `json:"embedding"`
	ModelVersion string    `json:"model_version"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserContentInteraction struct {
	UserID          uuid.UUID       `json:"user_id"`
	ContentID       uuid.UUID       `json:"content_id"`
	InteractionType InteractionType `json:"interaction_type"`
	Weight          float64         `json:"weight"`
	Timestamp       time.Time       `json:"timestamp"`
}

type Recommendation struct {
	ContentID uuid.UUID `json:"content_id"`
	Score     float64   `json:"score"`
	Reason    string    `json:"reason"`
	ReasonSet []string  `json:"reason_set,omitempty"`
}

type Feed struct {
	Type     FeedType         `json:"type"`
	UserID   uuid.UUID        `json:"user_id"`
	Items    []Recommendation `json:"items"`
	ServedAt time.Time        `json:"served_at"`
}

type ModelInfo struct {
	Version    string    `json:"version"`
	DeployedAt time.Time `json:"deployed_at"`
	Metrics    string    `json:"metrics"`
	IsActive   bool      `json:"is_active"`
}
