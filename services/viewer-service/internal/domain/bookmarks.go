package domain

import (
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	ID        uuid.UUID `json:"id"`
	ViewerID  uuid.UUID `json:"viewer_id"`
	ContentID uuid.UUID `json:"content_id"`
	Note      string    `json:"note,omitempty"`
	Folder    string    `json:"folder,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type WatchLaterItem struct {
	ID        uuid.UUID `json:"id"`
	ViewerID  uuid.UUID `json:"viewer_id"`
	ContentID uuid.UUID `json:"content_id"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}
