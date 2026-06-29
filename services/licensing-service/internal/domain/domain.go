package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventProducer interface {
	Publish(ctx context.Context, eventType string, data interface{})
	Close()
}

type LicenseType string

const (
	LicenseExclusive    LicenseType = "exclusive"
	LicenseNonExclusive LicenseType = "non_exclusive"
	LicenseSync         LicenseType = "sync"
	LicensePerformance  LicenseType = "performance"
)

type LicenseScope string

const (
	LicenseScopePlatform LicenseScope = "platform"
	LicenseScopeTerritory LicenseScope = "territory"
	LicenseScopeGlobal    LicenseScope = "global"
)

type RateType string

const (
	RateTypeFlat          RateType = "flat"
	RateTypeFixedPerUse   RateType = "fixed_per_use"
	RateTypeRevenueShare  RateType = "revenue_share"
)

type LicenseStatus string

const (
	LicenseDraft      LicenseStatus = "draft"
	LicensePending    LicenseStatus = "pending"
	LicenseActive     LicenseStatus = "active"
	LicenseExpired    LicenseStatus = "expired"
	LicenseTerminated LicenseStatus = "terminated"
)

type License struct {
	ID                   uuid.UUID       `json:"id"`
	RightsHolderID       uuid.UUID       `json:"rightsHolderId"`
	LicenseeID           uuid.UUID       `json:"licenseeId"`
	ContentID            *uuid.UUID      `json:"contentId,omitempty"`
	Type                 LicenseType     `json:"type"`
	Scope                LicenseScope    `json:"scope"`
	Territory            []string        `json:"territory"`
	StartDate            time.Time       `json:"startDate"`
	EndDate              time.Time       `json:"endDate"`
	AutoRenew            bool            `json:"autoRenew"`
	RateType             RateType        `json:"rateType"`
	RateCents            int64           `json:"rateCents"`
	RevenueSharePercent  float64         `json:"revenueSharePercent"`
	MinGuaranteeCents    int64           `json:"minGuaranteeCents"`
	Status               LicenseStatus   `json:"status"`
	TermsURL             string          `json:"termsUrl"`
	CreatedAt            time.Time       `json:"createdAt"`
}

type ContentRight struct {
	ID              uuid.UUID       `json:"id"`
	ContentID       uuid.UUID       `json:"contentId"`
	RightsHolderID  uuid.UUID       `json:"rightsHolderId"`
	LicenseID       uuid.UUID       `json:"licenseId"`
	Restrictions    json.RawMessage `json:"restrictions"`
	RegisteredAt    time.Time       `json:"registeredAt"`
}

type UsageType string

const (
	UsageStream    UsageType = "stream"
	UsageDownload  UsageType = "download"
	UsagePerformance UsageType = "performance"
	UsageSync      UsageType = "sync"
)

type UsageLog struct {
	ID         uuid.UUID       `json:"id"`
	LicenseID  uuid.UUID       `json:"licenseId"`
	ContentID  uuid.UUID       `json:"contentId"`
	UsageType  UsageType       `json:"usageType"`
	Context    json.RawMessage `json:"context"`
	RecordedAt time.Time       `json:"recordedAt"`
}

type RoyaltyStatus string

const (
	RoyaltyPending  RoyaltyStatus = "pending"
	RoyaltyPaid     RoyaltyStatus = "paid"
	RoyaltyDisputed RoyaltyStatus = "disputed"
)

type RoyaltyStatement struct {
	ID              uuid.UUID     `json:"id"`
	LicenseID       uuid.UUID     `json:"licenseId"`
	RightsHolderID  uuid.UUID     `json:"rightsHolderId"`
	PeriodStart     time.Time     `json:"periodStart"`
	PeriodEnd       time.Time     `json:"periodEnd"`
	UsageCount      int64         `json:"usageCount"`
	RateApplied     int64         `json:"rateApplied"`
	TotalCents      int64         `json:"totalCents"`
	Status          RoyaltyStatus `json:"status"`
	PaidAt          *time.Time    `json:"paidAt,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
}

type RightsCheckResult struct {
	ContentID   uuid.UUID `json:"contentId"`
	Allowed     bool      `json:"allowed"`
	LicenseID   *uuid.UUID `json:"licenseId,omitempty"`
	Reason      string    `json:"reason"`
}

type ComplianceReport struct {
	TotalLicenses       int64 `json:"totalLicenses"`
	ActiveLicenses      int64 `json:"activeLicenses"`
	ExpiringLicenses    int64 `json:"expiringLicenses"`
	UsageRecords        int64 `json:"usageRecords"`
	PendingRoyalties    int64 `json:"pendingRoyalties"`
	PaidRoyalties       int64 `json:"paidRoyalties"`
	TotalRoyaltyCents   int64 `json:"totalRoyaltyCents"`
	FlagsRaised         int64 `json:"flagsRaised"`
}
