package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenHash    string    `json:"-"`
	RefreshHash  string    `json:"-"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
