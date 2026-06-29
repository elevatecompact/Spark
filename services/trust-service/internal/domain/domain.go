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

type TrustLevel string

const (
	TrustLow     TrustLevel = "low"
	TrustMedium  TrustLevel = "medium"
	TrustHigh    TrustLevel = "high"
	TrustVerified TrustLevel = "verified"
)

type ReputationScore struct {
	UserID                uuid.UUID   `json:"userId"`
	OverallScore          int         `json:"overallScore"`
	TrustLevel            TrustLevel  `json:"trustLevel"`
	PositiveSignalWeight  int         `json:"positiveSignalWeight"`
	NegativeSignalWeight  int         `json:"negativeSignalWeight"`
	ScoreDecayFactor      float64     `json:"scoreDecayFactor"`
	ModelVersion          string      `json:"modelVersion"`
	CalculatedAt          time.Time   `json:"calculatedAt"`
	NextRecalculationAt   time.Time   `json:"nextRecalculationAt"`
}

type SignalType string

const (
	SignalPositive SignalType = "positive"
	SignalNegative SignalType = "negative"
)

type SignalCategory string

const (
	SignalIdentity  SignalCategory = "identity"
	SignalPayment   SignalCategory = "payment"
	SignalContent   SignalCategory = "content"
	SignalCommunity SignalCategory = "community"
	SignalAccountAge SignalCategory = "account_age"
)

type TrustSignal struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"userId"`
	SignalType       SignalType     `json:"signalType"`
	Category         SignalCategory `json:"category"`
	Weight           int            `json:"weight"`
	Description      string         `json:"description"`
	SourceEntityType string         `json:"sourceEntityType"`
	SourceEntityID   string         `json:"sourceEntityId"`
	ExpiresAt        *time.Time     `json:"expiresAt,omitempty"`
	RecordedAt       time.Time      `json:"recordedAt"`
}

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type RecommendedAction string

const (
	ActionAllow    RecommendedAction = "allow"
	ActionChallenge RecommendedAction = "challenge"
	ActionBlock    RecommendedAction = "block"
	ActionReview   RecommendedAction = "review"
)

type RiskAssessment struct {
	ID                uuid.UUID          `json:"id"`
	UserID            uuid.UUID          `json:"userId"`
	ActionType        string             `json:"actionType"`
	Context           json.RawMessage    `json:"context"`
	RiskScore         float64            `json:"riskScore"`
	RiskLevel         RiskLevel          `json:"riskLevel"`
	TriggeredRules    []string           `json:"triggeredRules"`
	RecommendedAction RecommendedAction  `json:"recommendedAction"`
	AssessedAt        time.Time          `json:"assessedAt"`
}

type FraudCaseType string

const (
	FraudPaymentFraud        FraudCaseType = "payment_fraud"
	FraudAccountTakeover     FraudCaseType = "account_takeover"
	FraudCoordinatedBehavior FraudCaseType = "coordinated_behavior"
	FraudSpam                FraudCaseType = "spam"
)

type FraudCaseStatus string

const (
	FraudOpen         FraudCaseStatus = "open"
	FraudInvestigating FraudCaseStatus = "investigating"
	FraudConfirmed    FraudCaseStatus = "confirmed"
	FraudFalsePositive FraudCaseStatus = "false_positive"
	FraudResolved     FraudCaseStatus = "resolved"
)

type FraudCase struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"userId"`
	CaseType         FraudCaseType   `json:"caseType"`
	Status           FraudCaseStatus `json:"status"`
	Evidence         json.RawMessage `json:"evidence"`
	AutomatedDecision string         `json:"automatedDecision"`
	ReviewedBy       *uuid.UUID      `json:"reviewedBy,omitempty"`
	ResolvedAt       *time.Time      `json:"resolvedAt,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
}

type RiskRule struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Category       string          `json:"category"`
	Conditions     json.RawMessage `json:"conditions"`
	RiskScoreImpact float64        `json:"riskScoreImpact"`
	Action         string          `json:"action"`
	IsActive       bool            `json:"isActive"`
	Priority       int             `json:"priority"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type TrustDashboard struct {
	TotalUsers          int64   `json:"totalUsers"`
	AvgReputationScore  float64 `json:"avgReputationScore"`
	LowTrustUsers       int64   `json:"lowTrustUsers"`
	MediumTrustUsers    int64   `json:"mediumTrustUsers"`
	HighTrustUsers      int64   `json:"highTrustUsers"`
	VerifiedUsers       int64   `json:"verifiedUsers"`
	OpenFraudCases      int64   `json:"openFraudCases"`
	RiskAssessmentsToday int64  `json:"riskAssessmentsToday"`
}

type ScoreDistribution struct {
	Range0_200   int64 `json:"range0_200"`
	Range201_400 int64 `json:"range201_400"`
	Range401_600 int64 `json:"range401_600"`
	Range601_800 int64 `json:"range601_800"`
	Range801_1000 int64 `json:"range801_1000"`
}

type FraudDetectionResult struct {
	IsFraudulent bool              `json:"isFraudulent"`
	Score        float64           `json:"score"`
	Reasons      []string          `json:"reasons"`
	CaseID       *uuid.UUID        `json:"caseId,omitempty"`
}
