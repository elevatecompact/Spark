package domain

import (
	"time"

	"github.com/google/uuid"
)

type CompetitionFormat string
const (
	FormatElimination CompetitionFormat = "elimination"
	FormatScore       CompetitionFormat = "score"
	FormatLeaderboard CompetitionFormat = "leaderboard"
)

type BracketType string
const (
	BracketSingle BracketType = "single"
	BracketDouble BracketType = "double"
)

type CompetitionStatus string
const (
	CompDraft     CompetitionStatus = "draft"
	CompOpen      CompetitionStatus = "open"
	CompStarted   CompetitionStatus = "started"
	CompCompleted CompetitionStatus = "completed"
	CompCancelled CompetitionStatus = "cancelled"
)

type ParticipantStatus string
const (
	PartRegistered ParticipantStatus = "registered"
	PartActive     ParticipantStatus = "active"
	PartEliminated ParticipantStatus = "eliminated"
	PartWithdrawn  ParticipantStatus = "withdrawn"
	PartWinner     ParticipantStatus = "winner"
)

type MatchStatus string
const (
	MatchScheduled  MatchStatus = "scheduled"
	MatchCompleted  MatchStatus = "completed"
	MatchDisputed   MatchStatus = "disputed"
)

type SubmissionStatus string
const (
	SubPending  SubmissionStatus = "pending"
	SubScored   SubmissionStatus = "scored"
	SubConfirmed SubmissionStatus = "confirmed"
)

type Competition struct {
	ID                 uuid.UUID         `json:"id"`
	CreatorID          uuid.UUID         `json:"creator_id"`
	Title              string            `json:"title"`
	Description        string            `json:"description"`
	Format             CompetitionFormat `json:"format"`
	BracketType        BracketType       `json:"bracket_type,omitempty"`
	Status             CompetitionStatus `json:"status"`
	MaxParticipants    int               `json:"max_participants"`
	RegistrationStart  *time.Time        `json:"registration_start,omitempty"`
	RegistrationEnd    *time.Time        `json:"registration_end,omitempty"`
	CompetitionStart   *time.Time        `json:"competition_start,omitempty"`
	CompetitionEnd     *time.Time        `json:"competition_end,omitempty"`
	PrizePoolCents     int64             `json:"prize_pool_cents"`
	Rules              []string          `json:"rules,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

type CompetitionParticipant struct {
	ID              uuid.UUID         `json:"id"`
	CompetitionID   uuid.UUID         `json:"competition_id"`
	UserID          uuid.UUID         `json:"user_id"`
	SeedNumber      int               `json:"seed_number"`
	Status          ParticipantStatus `json:"status"`
	BracketPosition *int              `json:"bracket_position,omitempty"`
	RegisteredAt    time.Time         `json:"registered_at"`
}

type Match struct {
	ID             uuid.UUID   `json:"id"`
	CompetitionID  uuid.UUID   `json:"competition_id"`
	Round          int         `json:"round"`
	BracketPosition int        `json:"bracket_position"`
	Participant1ID *uuid.UUID  `json:"participant1_id,omitempty"`
	Participant2ID *uuid.UUID  `json:"participant2_id,omitempty"`
	WinnerID       *uuid.UUID  `json:"winner_id,omitempty"`
	Scores         map[string]interface{} `json:"scores,omitempty"`
	Status         MatchStatus `json:"status"`
	ScheduledAt    *time.Time  `json:"scheduled_at,omitempty"`
	CompletedAt    *time.Time  `json:"completed_at,omitempty"`
}

type JudgingSubmission struct {
	ID              uuid.UUID       `json:"id"`
	CompetitionID   uuid.UUID       `json:"competition_id"`
	ParticipantID   uuid.UUID       `json:"participant_id"`
	Scores          map[string]float64 `json:"scores"`
	AverageScore    float64         `json:"average_score"`
	Status          SubmissionStatus `json:"status"`
}

type LeaderboardEntry struct {
	ParticipantID uuid.UUID `json:"participant_id"`
	UserID        uuid.UUID `json:"user_id"`
	Score         float64   `json:"score"`
	Rank          int       `json:"rank"`
}

type CompetitionAdminStats struct {
	TotalCompetitions int `json:"total_competitions"`
	TotalParticipants int `json:"total_participants"`
	ActiveCompetitions int `json:"active_competitions"`
	PrizePoolTotal   int64 `json:"prize_pool_total"`
}
