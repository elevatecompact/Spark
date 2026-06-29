package domain

import (
	"time"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelPush Channel = "push"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelInApp Channel = "inapp"
)

type NotificationType string

const (
	NotifNewSubscriber  NotificationType = "new_subscriber"
	NotifGiftReceived   NotificationType = "gift_received"
	NotifStreamLive     NotificationType = "stream_live"
	NotifPayoutComplete NotificationType = "payout_complete"
	NotifMention        NotificationType = "mention"
	NotifDirectMessage  NotificationType = "direct_message"
	NotifWelcome        NotificationType = "welcome"
)

type Platform string

const (
	PlatformIOS  Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWeb    Platform = "web"
)

type Notification struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
	Data      string          `json:"data"`
	Channel   Channel         `json:"channel"`
	ReadAt    *time.Time      `json:"read_at,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type NotificationPreference struct {
	UserID      uuid.UUID `json:"user_id"`
	Preferences string    `json:"preferences"`
}

type PushDevice struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Platform  Platform  `json:"platform"`
	Token     string    `json:"token"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Template struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type"`
	SubjectTemplate string    `json:"subject_template"`
	BodyTemplate    string    `json:"body_template"`
	Channels        []string  `json:"channels"`
	CreatedAt       time.Time `json:"created_at"`
}

type SendNotificationRequest struct {
	UserID uuid.UUID       `json:"user_id"`
	Type   NotificationType `json:"type"`
	Title  string          `json:"title"`
	Body   string          `json:"body"`
	Data   string          `json:"data"`
}

type RegisterDeviceRequest struct {
	Platform Platform `json:"platform"`
	Token    string   `json:"token"`
}

type UpdatePreferencesRequest struct {
	Preferences string `json:"preferences"`
}
