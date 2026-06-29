package domain

import (
	"time"
	"github.com/google/uuid"
)

type CreatorStatus string

const (
	CreatorActive    CreatorStatus = "active"
	CreatorInactive  CreatorStatus = "inactive"
	CreatorSuspended CreatorStatus = "suspended"
)

type SocialLinks struct {
	Website   string `json:"website,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	TikTok    string `json:"tiktok,omitempty"`
	GitHub    string `json:"github,omitempty"`
	Discord   string `json:"discord,omitempty"`
}

type Creator struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	DisplayName     string       `json:"display_name"`
	Bio             string       `json:"bio"`
	AvatarURL       string       `json:"avatar_url,omitempty"`
	BannerURL       string       `json:"banner_url,omitempty"`
	Categories      []string     `json:"categories,omitempty"`
	Tags            []string     `json:"tags,omitempty"`
	Language        string       `json:"language"`
	Country         string       `json:"country"`
	Timezone        string       `json:"timezone,omitempty"`
	SocialLinks     SocialLinks  `json:"social_links,omitempty"`
	Verified        bool         `json:"verified"`
	VerifiedAt      *time.Time   `json:"verified_at,omitempty"`
	Status          CreatorStatus `json:"status"`
	FollowerCount   int          `json:"follower_count"`
	SubscriberCount int          `json:"subscriber_count"`
	TotalViews      int64        `json:"total_views"`
	TotalStreams    int          `json:"total_streams"`
	Level           int          `json:"level"`
	Rank            int          `json:"rank"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type CreateCreatorRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Bio         string `json:"bio" validate:"max=2000"`
	Language    string `json:"language" validate:"required,len=2"`
	Country     string `json:"country" validate:"required,len=2"`
	Timezone    string `json:"timezone"`
}

type UpdateCreatorRequest struct {
	DisplayName *string      `json:"display_name,omitempty"`
	Bio         *string      `json:"bio,omitempty"`
	AvatarURL   *string      `json:"avatar_url,omitempty"`
	BannerURL   *string      `json:"banner_url,omitempty"`
	Categories  *[]string    `json:"categories,omitempty"`
	Tags        *[]string    `json:"tags,omitempty"`
	Language    *string      `json:"language,omitempty"`
	Country     *string      `json:"country,omitempty"`
	Timezone    *string      `json:"timezone,omitempty"`
	SocialLinks *SocialLinks `json:"social_links,omitempty"`
}
