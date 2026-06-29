package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventProducer interface {
	Publish(ctx context.Context, eventType string, data interface{})
	Close()
}

type Category struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Slug         string     `json:"slug"`
	Description  string     `json:"description"`
	ParentID     *uuid.UUID `json:"parentId,omitempty"`
	IconURL      string     `json:"iconUrl"`
	SortOrder    int        `json:"sortOrder"`
	IsActive     bool       `json:"isActive"`
	ContentCount int64      `json:"contentCount"`
}

type CollectionType string

const (
	CollectionEditorial CollectionType = "editorial"
	CollectionHoliday   CollectionType = "holiday"
	CollectionTheme     CollectionType = "theme"
)

type Collection struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Type          CollectionType `json:"type"`
	CoverImageURL string         `json:"coverImageUrl"`
	IsFeatured    bool           `json:"isFeatured"`
	StartAt       *time.Time     `json:"startAt,omitempty"`
	EndAt         *time.Time     `json:"endAt,omitempty"`
	CuratedBy     string         `json:"curatedBy"`
	Items         []CollectionItem `json:"items,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
}

type CollectionItem struct {
	CollectionID uuid.UUID `json:"collectionId"`
	ContentID    uuid.UUID `json:"contentId"`
	SortOrder    int       `json:"sortOrder"`
	AddedAt      time.Time `json:"addedAt"`
}

type PickType string

const (
	PickStaffPick  PickType = "staff_pick"
	PickSpotlight  PickType = "spotlight"
	PickHoliday    PickType = "holiday"
)

type EditorialPick struct {
	ContentID uuid.UUID `json:"contentId"`
	PickType  PickType  `json:"pickType"`
	Label     string    `json:"label"`
	Reason    string    `json:"reason"`
	PickedBy  string    `json:"pickedBy"`
	StartAt   time.Time `json:"startAt"`
	EndAt     time.Time `json:"endAt"`
	SortOrder int       `json:"sortOrder"`
}

type TrendingItem struct {
	ContentID uuid.UUID `json:"contentId"`
	Score     float64   `json:"score"`
	Rank      int       `json:"rank"`
}

type FeedType string

const (
	FeedHome     FeedType = "home"
	FeedTrending FeedType = "trending"
	FeedCategory FeedType = "category"
	FeedNew      FeedType = "new"
	FeedRelated  FeedType = "related"
)