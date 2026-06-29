package domain

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID        uuid.UUID `json:"id"`
	ViewerID  uuid.UUID `json:"viewer_id"`
	ContentID uuid.UUID `json:"content_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
)

type Reaction struct {
	ID        uuid.UUID    `json:"id"`
	ViewerID  uuid.UUID    `json:"viewer_id"`
	ContentID uuid.UUID    `json:"content_id"`
	Type      ReactionType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type ReportType string

const (
	ReportSpam       ReportType = "spam"
	ReportHarassment ReportType = "harassment"
	ReportCopyright  ReportType = "copyright"
	ReportOther      ReportType = "other"
)

type Report struct {
	ID          uuid.UUID  `json:"id"`
	ViewerID    uuid.UUID  `json:"viewer_id"`
	ContentID   uuid.UUID  `json:"content_id"`
	Type        ReportType `json:"type"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
