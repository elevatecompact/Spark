package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentText  ContentType = "text"
	ContentEmote ContentType = "emote"
	ContentMedia ContentType = "media"
	ContentSystem ContentType = "system"
)

type ModerationStatus string

const (
	ModPending  ModerationStatus = "pending"
	ModApproved ModerationStatus = "approved"
	ModRejected ModerationStatus = "rejected"
)

type ChatMessage struct {
	ID               uuid.UUID        `json:"id"`
	RoomID           uuid.UUID        `json:"room_id"`
	UserID           uuid.UUID        `json:"user_id"`
	Username         string           `json:"username"`
	Content          string           `json:"content"`
	ContentType      ContentType      `json:"content_type"`
	ModerationStatus ModerationStatus `json:"moderation_status"`
	Emotes           []string         `json:"emotes,omitempty"`
	EditedAt         *time.Time       `json:"edited_at,omitempty"`
	DeletedAt        *time.Time       `json:"deleted_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}

type SendMessageRequest struct {
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
}

type MessageCursor struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit"`
}
