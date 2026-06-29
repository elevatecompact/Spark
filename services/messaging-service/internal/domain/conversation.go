package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConversationType string

const (
	ConvDirect ConversationType = "direct"
	ConvGroup  ConversationType = "group"
)

type MemberRole string

const (
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

type Conversation struct {
	ID        uuid.UUID        `json:"id"`
	Type      ConversationType `json:"type"`
	Name      *string          `json:"name,omitempty"`
	IconURL   *string          `json:"icon_url,omitempty"`
	CreatedBy uuid.UUID        `json:"created_by"`
	IsActive  bool             `json:"is_active"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type ConversationMember struct {
	ConversationID  uuid.UUID  `json:"conversation_id"`
	UserID          uuid.UUID  `json:"user_id"`
	Role            MemberRole `json:"role"`
	LastReadMsgID   *uuid.UUID `json:"last_read_message_id,omitempty"`
	JoinedAt        time.Time  `json:"joined_at"`
}

type CreateConversationRequest struct {
	Type       ConversationType `json:"type"`
	Name       *string          `json:"name,omitempty"`
	MemberIDs  []uuid.UUID      `json:"member_ids"`
}
