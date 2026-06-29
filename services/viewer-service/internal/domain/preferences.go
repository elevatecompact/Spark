package domain

import (
	"time"

	"github.com/google/uuid"
)

type ViewerPreferences struct {
	ViewerID            uuid.UUID              `json:"viewer_id"`
	PreferredCategories []uuid.UUID            `json:"preferred_categories"`
	ContentLanguage     string                 `json:"content_language"`
	Autoplay            bool                   `json:"autoplay"`
	MatureContentAllowed bool                  `json:"mature_content_allowed"`
	NotificationPrefs   map[string]interface{} `json:"notification_prefs"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

type UpdatePreferences struct {
	PreferredCategories  *[]uuid.UUID          `json:"preferred_categories,omitempty"`
	ContentLanguage      *string               `json:"content_language,omitempty"`
	Autoplay             *bool                 `json:"autoplay,omitempty"`
	MatureContentAllowed *bool                 `json:"mature_content_allowed,omitempty"`
	NotificationPrefs    *map[string]interface{} `json:"notification_prefs,omitempty"`
}
