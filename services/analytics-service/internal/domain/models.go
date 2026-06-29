package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DashboardType string

const (
	DashCreator DashboardType = "creator"
	DashViewer  DashboardType = "viewer"
	DashAdmin   DashboardType = "admin"
)

type ReportStatus string

const (
	ReportPending    ReportStatus = "pending"
	ReportGenerating ReportStatus = "generating"
	ReportReady      ReportStatus = "ready"
	ReportFailed     ReportStatus = "failed"
)

type TrackedEvent struct {
	ID          uuid.UUID       `json:"id"`
	EventName   string          `json:"event_name"`
	UserID      uuid.UUID       `json:"user_id"`
	SessionID   string          `json:"session_id"`
	Properties  json.RawMessage `json:"properties"`
	Context     json.RawMessage `json:"context"`
	EventTime   time.Time       `json:"event_time"`
	CreatedAt   time.Time       `json:"created_at"`
}

type MetricAggregate struct {
	MetricName  string          `json:"metric_name"`
	TimeBucket  time.Time       `json:"time_bucket"`
	Dimensions  json.RawMessage `json:"dimensions"`
	Count       int64           `json:"count"`
	Sum         float64         `json:"sum"`
	Avg         float64         `json:"avg"`
	P50         float64         `json:"p50"`
	P95         float64         `json:"p95"`
	P99         float64         `json:"p99"`
}

type Dashboard struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	DashType   DashboardType   `json:"dash_type"`
	Config     json.RawMessage `json:"config"`
	Data       json.RawMessage `json:"data,omitempty"`
	CacheUntil *time.Time      `json:"cache_until,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type Report struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Config       json.RawMessage `json:"config"`
	Status       ReportStatus    `json:"status"`
	DownloadURL  string          `json:"download_url,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
}

type ReportTemplate struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
}

type Funnel struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Name      string          `json:"name"`
	Steps     json.RawMessage `json:"steps"`
	Results   json.RawMessage `json:"results,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type TrackEventRequest struct {
	EventName  string          `json:"event_name"`
	SessionID  string          `json:"session_id"`
	Properties json.RawMessage `json:"properties,omitempty"`
	Context    json.RawMessage `json:"context,omitempty"`
	EventTime  *time.Time      `json:"event_time,omitempty"`
}

type MetricQuery struct {
	MetricName string            `json:"metric_name"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Granularity string           `json:"granularity"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
	Aggregates []string          `json:"aggregates"`
}

type FunnelStep struct {
	Name  string `json:"name"`
	Event string `json:"event"`
	Order int    `json:"order"`
}

type FunnelDefinition struct {
	Name  string       `json:"name"`
	Steps []FunnelStep `json:"steps"`
}
