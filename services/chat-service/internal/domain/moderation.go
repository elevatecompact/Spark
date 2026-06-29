package domain

import (
	"time"

	"github.com/google/uuid"
)

type MutedUser struct {
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type BannedUser struct {
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SlowModeConfig struct {
	RoomID       uuid.UUID `json:"room_id"`
	IntervalSecs int       `json:"interval_seconds"`
	Enabled      bool      `json:"enabled"`
}
