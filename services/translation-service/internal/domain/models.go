package domain

import (
	"time"

	"github.com/google/uuid"
)

type TranslationJobStatus string

const (
	JobPending    TranslationJobStatus = "pending"
	JobProcessing TranslationJobStatus = "processing"
	JobCompleted  TranslationJobStatus = "completed"
	JobFailed     TranslationJobStatus = "failed"
)

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending"
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "rejected"
)

type TranslationRequest struct {
	Text         string `json:"text"`
	SourceLang   string `json:"sourceLang,omitempty"`
	TargetLang   string `json:"targetLang"`
}

type TranslationResult struct {
	TranslatedText string  `json:"translatedText"`
	SourceLang     string  `json:"sourceLang"`
	TargetLang     string  `json:"targetLang"`
	Provider       string  `json:"provider"`
	Confidence     float64 `json:"confidence"`
}

type DetectionResult struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
}

type TranslationMemoryEntry struct {
	ID             uuid.UUID `json:"id"`
	SourceHash     string    `json:"sourceHash"`
	SourceText     string    `json:"sourceText"`
	TranslatedText string    `json:"translatedText"`
	SourceLang     string    `json:"sourceLang"`
	TargetLang     string    `json:"targetLang"`
	Provider       string    `json:"provider"`
	QualityScore   float64   `json:"qualityScore"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type TranslationJob struct {
	ID          uuid.UUID             `json:"id"`
	ContentType string                `json:"contentType"`
	ContentID   uuid.UUID             `json:"contentId"`
	Status      TranslationJobStatus  `json:"status"`
	Languages   []string              `json:"languages"`
	CreatedAt   time.Time             `json:"createdAt"`
}

type ReviewEntry struct {
	ID                uuid.UUID    `json:"id"`
	TranslationID     uuid.UUID    `json:"translationId"`
	OriginalText      string       `json:"originalText"`
	TranslatedText    string       `json:"translatedText"`
	SourceLang        string       `json:"sourceLang"`
	TargetLang        string       `json:"targetLang"`
	ReviewerID        *uuid.UUID   `json:"reviewerId,omitempty"`
	Status            ReviewStatus `json:"status"`
	CorrectedText     string       `json:"correctedText,omitempty"`
	ReviewedAt        *time.Time   `json:"reviewedAt,omitempty"`
}

type SupportedLanguage struct {
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Coverage float64 `json:"coverage"`
}

type ProviderUsage struct {
	Provider      string `json:"provider"`
	RequestCount  int64  `json:"requestCount"`
	CharCount     int64  `json:"charCount"`
}
