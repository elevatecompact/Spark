package domain

import (
	"time"

	"github.com/google/uuid"
)

type Emote struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	ImageURL  string    `json:"image_url"`
	IsGlobal  bool      `json:"is_global"`
	RoomID    *uuid.UUID `json:"room_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
