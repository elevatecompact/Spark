package domain

import (
	"time"

	"github.com/google/uuid"
)

type RuleCategory string

const (
	RuleHarassment RuleCategory = "harassment"
	RuleSpam       RuleCategory = "spam"
	RuleNSFW       RuleCategory = "nsfw"
	RuleViolence   RuleCategory = "violence"
	RuleHateSpeech RuleCategory = "hate_speech"
)

type Severity string

const (
	SevWarn    Severity = "warn"
	SevRestrict Severity = "restrict"
	SevRemove  Severity = "remove"
	SevSuspend Severity = "suspend"
)

type ActionStatus string

const (
	ActionPending  ActionStatus = "pending"
	ActionApplied  ActionStatus = "applied"
	ActionAppealed ActionStatus = "appealed"
	ActionReversed ActionStatus = "reversed"
)

type ReviewStatus string

const (
	ReviewPending   ReviewStatus = "pending"
	ReviewReviewing ReviewStatus = "reviewing"
	ReviewResolved  ReviewStatus = "resolved"
)

type ReportStatus string

const (
	ReportOpen         ReportStatus = "open"
	ReportInvestigating ReportStatus = "investigating"
	ReportResolved     ReportStatus = "resolved"
)

type ModerationRule struct {
	ID         uuid.UUID    `json:"id"`
	Name       string       `json:"name"`
	Category   RuleCategory `json:"category"`
	Severity   Severity     `json:"severity"`
	Conditions map[string]interface{} `json:"conditions"`
	IsActive   bool         `json:"is_active"`
	Priority   int          `json:"priority"`
	CreatedAt  time.Time    `json:"created_at"`
}

type ModerationAction struct {
	ID         uuid.UUID    `json:"id"`
	UserID     uuid.UUID    `json:"user_id"`
	ContentID  *uuid.UUID   `json:"content_id,omitempty"`
	RuleID     uuid.UUID    `json:"rule_id"`
	ActionType Severity     `json:"action_type"`
	Status     ActionStatus `json:"status"`
	AppliedBy  string       `json:"applied_by"`
	Reason     string       `json:"reason"`
	Duration   *int         `json:"duration,omitempty"`
	AppliedAt  time.Time    `json:"applied_at"`
}

type ReviewItem struct {
	ID                uuid.UUID    `json:"id"`
	ContentType       string       `json:"content_type"`
	ContentID         uuid.UUID    `json:"content_id"`
	FlaggedBy         string       `json:"flagged_by"`
	Reasons           []string     `json:"reasons"`
	AssignedModerator *uuid.UUID   `json:"assigned_moderator,omitempty"`
	Status            ReviewStatus `json:"status"`
	Resolution        *string      `json:"resolution,omitempty"`
	ResolvedAt        *time.Time   `json:"resolved_at,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
}

type ContentReport struct {
	ID            uuid.UUID    `json:"id"`
	ReporterID    uuid.UUID    `json:"reporter_id"`
	ContentType   string       `json:"content_type"`
	ContentID     uuid.UUID    `json:"content_id"`
	Reason        string       `json:"reason"`
	Description   string       `json:"description"`
	Status        ReportStatus `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
}

type ScanResult struct {
	ContentID    uuid.UUID         `json:"content_id"`
	ContentType  string            `json:"content_type"`
	Violations   []ScanViolation   `json:"violations"`
	AutoAction   *ModerationAction `json:"auto_action,omitempty"`
	NeedsReview  bool              `json:"needs_review"`
}

type ScanViolation struct {
	RuleID     uuid.UUID    `json:"rule_id"`
	Severity   Severity     `json:"severity"`
	Category   RuleCategory `json:"category"`
	Confidence float64      `json:"confidence"`
	Matched    string       `json:"matched"`
}

type QueueStats struct {
	ContentType string `json:"content_type"`
	Pending     int    `json:"pending"`
	Reviewing   int    `json:"reviewing"`
	Resolved    int    `json:"resolved"`
}

type AdminStats struct {
	TotalScans     int64 `json:"total_scans"`
	TotalActions   int64 `json:"total_actions"`
	QueueDepth     int   `json:"queue_depth"`
	AutoActionRate float64 `json:"auto_action_rate"`
}
