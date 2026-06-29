package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentCreators  ContentType = "creators"
	ContentStreams   ContentType = "streams"
	ContentRecordings ContentType = "recordings"
	ContentClips     ContentType = "clips"
)

type SortBy string

const (
	SortRelevance   SortBy = "relevance"
	SortDate        SortBy = "date"
	SortPopularity  SortBy = "popularity"
)

type SearchDocument struct {
	ID              uuid.UUID   `json:"id"`
	ContentType     ContentType `json:"content_type"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	CreatorName     string      `json:"creator_name,omitempty"`
	Category        string      `json:"category,omitempty"`
	Tags            []string    `json:"tags,omitempty"`
	ViewCount       int64       `json:"view_count"`
	FollowerCount   int64       `json:"follower_count,omitempty"`
	Duration        int64       `json:"duration,omitempty"`
	Status          string      `json:"status,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
}

type SearchQuery struct {
	Query    string            `json:"q"`
	Type     ContentType       `json:"type"`
	Filters  map[string]string `json:"filters,omitempty"`
	Sort     SortBy            `json:"sort"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
	UserID   uuid.UUID         `json:"user_id,omitempty"`
}

type SearchResult struct {
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	Size    int               `json:"size"`
	Results []SearchDocument  `json:"results"`
}

type AutocompleteSuggestion struct {
	Text    string      `json:"text"`
	Type    ContentType `json:"type"`
	Score   float64     `json:"score"`
}

type SynonymSet struct {
	ID      string   `json:"id"`
	Synonyms []string `json:"synonyms"`
}

type IndexStats struct {
	ContentType   ContentType `json:"content_type"`
	DocCount      int64       `json:"doc_count"`
	IndexSize     string      `json:"index_size"`
}

type SearchAnalytics struct {
	Query     string    `json:"query"`
	ResultIDs []string  `json:"result_ids"`
	LatencyMs int64     `json:"latency_ms"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type ESHealth struct {
	Status      string `json:"status"`
	NodeCount   int    `json:"node_count"`
	ActiveShards int   `json:"active_shards"`
}
