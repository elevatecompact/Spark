package types

import "time"

type MediaStatus string

const (
	MediaPending    MediaStatus = "pending"
	MediaProcessing MediaStatus = "processing"
	MediaReady      MediaStatus = "ready"
	MediaFailed     MediaStatus = "failed"
)

type MediaType string

const (
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeFile     MediaType = "file"
)

type Media struct {
	ID           ID            `json:"id"`
	UserID       ID            `json:"user_id"`
	StreamID     *ID           `json:"stream_id,omitempty"`
	Type         MediaType     `json:"type"`
	Status       MediaStatus   `json:"status"`
	URL          string        `json:"url"`
	OriginalName string        `json:"original_name"`
	MimeType     string        `json:"mime_type"`
	Size         int64         `json:"size"`
	Width        int           `json:"width,omitempty"`
	Height       int           `json:"height,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	ThumbnailURL string        `json:"thumbnail_url,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type UploadPolicy struct {
	MaxFileSize    int64    `json:"max_file_size"`
	AllowedTypes   []string `json:"allowed_types"`
	MaxDuration    int      `json:"max_duration,omitempty"`
	RequireAuth    bool     `json:"require_auth"`
	RequireApproval bool    `json:"require_approval"`
}

type MediaUploadRequest struct {
	Filename     string            `json:"filename"`
	ContentType  string            `json:"content_type"`
	Size         int64             `json:"size"`
	Type         MediaType         `json:"type"`
	StreamID     *ID               `json:"stream_id,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type MediaUploadResponse struct {
	ID       ID     `json:"id"`
	UploadURL string `json:"upload_url"`
	Policy   UploadPolicy `json:"policy"`
}
