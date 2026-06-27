package domain

import (
	"time"
	"github.com/google/uuid"
)

type PortfolioItem struct {
	ID           uuid.UUID json:"id"
	CreatorID    uuid.UUID json:"creator_id"
	Title        string    json:"title"
	Description  string    json:"description,omitempty"
	MediaURL     string    json:"media_url"
	MediaType    string    json:"media_type"
	ThumbnailURL string    json:"thumbnail_url,omitempty"
	Featured     bool      json:"featured"
	SortOrder    int       json:"sort_order"
	CreatedAt    time.Time json:"created_at"
}

type CreatePortfolioItemRequest struct {
	Title        string json:"title" validate:"required,min=2,max=200"
	Description  string json:"description,omitempty"
	MediaURL     string json:"media_url" validate:"required,url"
	MediaType    string json:"media_type" validate:"required,oneof=video image audio"
	ThumbnailURL string json:"thumbnail_url,omitempty"
	Featured     bool   json:"featured"
	SortOrder    int    json:"sort_order"
}
