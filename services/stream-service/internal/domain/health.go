package domain

import (
	"time"

	"github.com/google/uuid"
)

type HealthReport struct {
	StreamID      uuid.UUID `json:"stream_id"`
	Bitrate       int       `json:"bitrate"`
	FPS           float64   `json:"fps"`
	PacketLoss    float64   `json:"packet_loss"`
	RoundTripTime float64   `json:"round_trip_time"`
	Jitter        float64   `json:"jitter"`
	Timestamp     time.Time `json:"timestamp"`
}

type StreamHealth struct {
	StreamID        uuid.UUID `json:"stream_id"`
	Bitrate         int       `json:"bitrate"`
	FPS             float64   `json:"fps"`
	PacketLoss      float64   `json:"packet_loss"`
	RoundTripTime   float64   `json:"round_trip_time"`
	Jitter          float64   `json:"jitter"`
	ConnectionScore int       `json:"connection_score"`
	Status          string    `json:"status"`
	LastUpdated     time.Time `json:"last_updated"`
}

type AnomalyReport struct {
	StreamID  uuid.UUID `json:"stream_id"`
	Anomalies []Anomaly `json:"anomalies"`
	Healthy   bool      `json:"healthy"`
}

type Anomaly struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Value       float64 `json:"value"`
	Threshold   float64 `json:"threshold"`
	Description string  `json:"description"`
}

type Recording struct {
	ID                 uuid.UUID  `json:"id"`
	StreamID           uuid.UUID  `json:"stream_id"`
	CreatorID          uuid.UUID  `json:"creator_id"`
	Title              string     `json:"title"`
	S3Key              string     `json:"s3_key"`
	Bucket             string     `json:"bucket"`
	Duration           int        `json:"duration"`
	FileSize           int64      `json:"file_size"`
	Width              int        `json:"width"`
	Height             int        `json:"height"`
	Codec              string     `json:"codec"`
	Status             string     `json:"status"`
	ProcessingProgress float64    `json:"processing_progress"`
	ThumbnailKey       string     `json:"thumbnail_key"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
