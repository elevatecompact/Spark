package domain

import (
	"time"

	"github.com/google/uuid"
)

type Passkey struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	CredentialID    string     `json:"credential_id"`
	PublicKey       []byte     `json:"-"`
	AttestationType string     `json:"attestation_type"`
	Transports      []string   `json:"transports"`
	AAGUID          *uuid.UUID `json:"aaguid,omitempty"`
	SignCount       int64      `json:"sign_count"`
	Name            string     `json:"name"`
	DeviceType      string     `json:"device_type"`
	BackedUp        bool       `json:"backed_up"`
	CreatedAt       time.Time  `json:"created_at"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
}
