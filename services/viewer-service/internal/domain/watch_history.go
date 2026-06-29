package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentLive     ContentType = "live"
	ContentRecorded ContentType = "recorded"
	ContentClip     ContentType = "clip"
)

type WatchHistory struct {
	ID                  uuid.UUID   `json:"id"`
	ViewerID            uuid.UUID   `json:"viewer_id"`
	ContentID           uuid.UUID   `json:"content_id"`
	ContentType         ContentType `json:"content_type"`
	Progress            float64     `json:"progress"`
	WatchDurationSeconds int        `json:"watch_duration_seconds"`
	Completed           bool        `json:"completed"`
	WatchedAt           time.Time   `json:"watched_at"`
	CreatedAt           time.Time   `json:"created_at"`
}

type WatchProgressUpdate struct {
	ContentID           uuid.UUID   `json:"content_id"`
	ContentType         ContentType `json:"content_type"`
	Progress            float64     `json:"progress"`
	WatchDurationSeconds int        `json:"watch_duration_seconds"`
	Completed           bool        `json:"completed"`
}
