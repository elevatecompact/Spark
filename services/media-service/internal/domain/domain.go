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

type ContentType string

const (
	ContentTypeVideo ContentType = "video"
	ContentTypeImage ContentType = "image"
	ContentTypeAudio ContentType = "audio"
)

type MediaStatus string

const (
	MediaStatusUploading  MediaStatus = "uploading"
	MediaStatusProcessing MediaStatus = "processing"
	MediaStatusReady      MediaStatus = "ready"
	MediaStatusFailed     MediaStatus = "failed"
	MediaStatusDeleted    MediaStatus = "deleted"
)

type MediaAsset struct {
	ID             uuid.UUID   `json:"id"`
	UploaderID     uuid.UUID   `json:"uploaderId"`
	ContentType    ContentType `json:"contentType"`
	SourceFilename string      `json:"sourceFilename"`
	FileSizeBytes  int64       `json:"fileSizeBytes"`
	MimeType       string      `json:"mimeType"`
	Status         MediaStatus `json:"status"`
	StoragePath    string      `json:"storagePath"`
	CDNURL         string      `json:"cdnUrl"`
	DurationSecs   float64     `json:"durationSeconds"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Checksum       string      `json:"checksum"`
	CreatedAt      time.Time   `json:"createdAt"`
}

type RenditionProfile string

const (
	RenditionThumbnail RenditionProfile = "thumbnail"
	Rendition720p      RenditionProfile = "720p"
	Rendition1080p     RenditionProfile = "1080p"
	RenditionSource    RenditionProfile = "source"
)

type RenditionFormat string

const (
	RenditionHLS   RenditionFormat = "hls"
	RenditionDASH  RenditionFormat = "dash"
	RenditionMP4   RenditionFormat = "mp4"
	RenditionWebP  RenditionFormat = "webp"
	RenditionJPG   RenditionFormat = "jpg"
)

type MediaRendition struct {
	ID           uuid.UUID        `json:"id"`
	MediaID      uuid.UUID        `json:"mediaId"`
	Profile      RenditionProfile `json:"profile"`
	Format       RenditionFormat  `json:"format"`
	FileSizeBytes int64           `json:"fileSizeBytes"`
	StoragePath  string           `json:"storagePath"`
	CDNURL       string           `json:"cdnUrl"`
	Status       MediaStatus      `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
}

type TranscodingJobStatus string

const (
	TranscodingPending    TranscodingJobStatus = "pending"
	TranscodingProcessing TranscodingJobStatus = "processing"
	TranscodingCompleted  TranscodingJobStatus = "completed"
	TranscodingFailed     TranscodingJobStatus = "failed"
)

type TranscodingJob struct {
	ID             uuid.UUID            `json:"id"`
	MediaID        uuid.UUID            `json:"mediaId"`
	Profiles       []RenditionProfile   `json:"profiles"`
	Status         TranscodingJobStatus `json:"status"`
	WorkerID       string               `json:"workerId"`
	StartedAt      *time.Time           `json:"startedAt,omitempty"`
	CompletedAt    *time.Time           `json:"completedAt,omitempty"`
	ErrorMessage   string               `json:"errorMessage,omitempty"`
	CreatedAt      time.Time            `json:"createdAt"`
}

type KeySystem string

const (
	KeySystemWidevine KeySystem = "widevine"
	KeySystemFairPlay KeySystem = "fairplay"
)

type DRMPolicy struct {
	ID                   uuid.UUID  `json:"id"`
	Name                 string     `json:"name"`
	ContentID            *uuid.UUID `json:"contentId,omitempty"`
	KeySystem            KeySystem  `json:"keySystem"`
	LicenseDurationSecs  int64      `json:"licenseDurationSeconds"`
	SecurityLevel        string     `json:"securityLevel"`
	IsActive             bool       `json:"isActive"`
	CreatedAt            time.Time  `json:"createdAt"`
}

type UploadSession struct {
	ID            uuid.UUID `json:"id"`
	UploaderID    uuid.UUID `json:"uploaderId"`
	Filename      string    `json:"filename"`
	FileSizeBytes int64     `json:"fileSizeBytes"`
	ContentType   string    `json:"contentType"`
	ChunksTotal   int       `json:"chunksTotal"`
	ChunksDone    int       `json:"chunksDone"`
	Checksum      string    `json:"checksum"`
	Status        string    `json:"status"`
	StoragePath   string    `json:"storagePath"`
	CreatedAt     time.Time `json:"createdAt"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

type StorageUsage struct {
	TotalBytes   int64 `json:"totalBytes"`
	UploadBytes  int64 `json:"uploadBytes"`
	RenditionBytes int64 `json:"renditionBytes"`
	ThumbnailBytes int64 `json:"thumbnailBytes"`
	AssetCount   int64 `json:"assetCount"`
}
