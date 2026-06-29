package domain

import (
	"time"
	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description,omitempty"`
	IconURL     string     `json:"icon_url,omitempty"`
	Color       string     `json:"color,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   int        `json:"sort_order"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	Slug        string     `json:"slug" validate:"required,min=2,max=100"`
	Description string     `json:"description,omitempty"`
	IconURL     string     `json:"icon_url,omitempty"`
	Color       string     `json:"color,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   int        `json:"sort_order"`
}
