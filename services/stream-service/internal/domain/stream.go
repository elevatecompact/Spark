package domain

import (
	"time"

	"github.com/google/uuid"
)

type StreamStatus string

const (
	StreamPending  StreamStatus = "pending"
	StreamStarting StreamStatus = "starting"
	StreamLive     StreamStatus = "live"
	StreamEnded    StreamStatus = "ended"
	StreamError    StreamStatus = "error"
	StreamArchived StreamStatus = "archived"
)

type Stream struct {
	ID                uuid.UUID    `json:"id"`
	CreatorID         uuid.UUID    `json:"creator_id"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	Category          string       `json:"category"`
	Tags              []string     `json:"tags"`
	ThumbnailURL      string       `json:"thumbnail_url"`
	StreamKey         string       `json:"stream_key"`
	RTMPEndpoint      string       `json:"rtmp_endpoint"`
	IngestProtocol    string       `json:"ingest_protocol"`
	Status            StreamStatus `json:"status"`
	StartedAt         *time.Time   `json:"started_at,omitempty"`
	EndedAt           *time.Time   `json:"ended_at,omitempty"`
	Duration          int          `json:"duration"`
	Width             int          `json:"width"`
	Height            int          `json:"height"`
	FrameRate         float64      `json:"frame_rate"`
	Bitrate           int          `json:"bitrate"`
	Codec             string       `json:"codec"`
	AvailableQualities []string    `json:"available_qualities"`
	ViewerCount       int          `json:"viewer_count"`
	PeakViewers       int          `json:"peak_viewers"`
	TotalViews        int64        `json:"total_views"`
	RecordEnabled     bool         `json:"record_enabled"`
	RecordingID       *uuid.UUID   `json:"recording_id,omitempty"`
	ChatEnabled       bool         `json:"chat_enabled"`
	AgeRestricted     bool         `json:"age_restricted"`
	DelaySeconds      int          `json:"delay_seconds"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type CreateStreamRequest struct {
	CreatorID     uuid.UUID `json:"creator_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	RecordEnabled bool      `json:"record_enabled"`
	ChatEnabled   bool      `json:"chat_enabled"`
	AgeRestricted bool      `json:"age_restricted"`
	DelaySeconds  int       `json:"delay_seconds"`
}

type UpdateStreamRequest struct {
	Title         *string   `json:"title,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Category      *string   `json:"category,omitempty"`
	Tags          *[]string `json:"tags,omitempty"`
	ThumbnailURL  *string   `json:"thumbnail_url,omitempty"`
	RecordEnabled *bool     `json:"record_enabled,omitempty"`
	ChatEnabled   *bool     `json:"chat_enabled,omitempty"`
	AgeRestricted *bool     `json:"age_restricted,omitempty"`
	DelaySeconds  *int      `json:"delay_seconds,omitempty"`
}

type StreamFilter struct {
	CreatorID *uuid.UUID
	Status    *StreamStatus
	Category  string
	Limit     int
	Offset    int
}

func (s *Stream) IsLive() bool {
	return s.Status == StreamLive
}

func (s *Stream) CanBeStarted() bool {
	return s.Status == StreamPending || s.Status == StreamStarting
}

func (s *Stream) GenerateRTMPEndpoint(rtmpDomain string) string {
	return "rtmp://" + rtmpDomain + "/live"
}
