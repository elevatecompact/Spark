package domain

import (
	"time"

	"github.com/google/uuid"
)

type RoomType string

const (
	RoomTypeStream   RoomType = "stream"
	RoomTypeChannel  RoomType = "channel"
	RoomTypeDM       RoomType = "dm"
	RoomTypeSystem   RoomType = "system"
)

type ChatRoom struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Type            RoomType  `json:"type"`
	OwnerID         uuid.UUID `json:"owner_id"`
	SlowModeSeconds int       `json:"slow_mode_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateRoomRequest struct {
	Name    string   `json:"name"`
	Type    RoomType `json:"type"`
	OwnerID uuid.UUID `json:"owner_id"`
}
