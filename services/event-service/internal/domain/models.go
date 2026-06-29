package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string
const (
	EventVirtual  EventType = "virtual"
	EventInPerson EventType = "inperson"
	EventHybrid   EventType = "hybrid"
)

type EventStatus string
const (
	StatusDraft     EventStatus = "draft"
	StatusPublished EventStatus = "published"
	StatusCancelled EventStatus = "cancelled"
	StatusCompleted EventStatus = "completed"
)

type AttendeeStatus string
const (
	AttendeeRegistered AttendeeStatus = "registered"
	AttendeeAttended   AttendeeStatus = "attended"
	AttendeeCancelled  AttendeeStatus = "cancelled"
	AttendeeNoShow     AttendeeStatus = "no_show"
)

type SeriesFrequency string
const (
	FreqDaily  SeriesFrequency = "daily"
	FreqWeekly SeriesFrequency = "weekly"
	FreqMonthly SeriesFrequency = "monthly"
)

type Event struct {
	ID             uuid.UUID    `json:"id"`
	CreatorID      uuid.UUID    `json:"creator_id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Category       string       `json:"category"`
	Type           EventType    `json:"type"`
	StartAt        time.Time    `json:"start_at"`
	EndAt          time.Time    `json:"end_at"`
	Timezone       string       `json:"timezone"`
	MaxAttendees   int          `json:"max_attendees"`
	StreamID       *uuid.UUID   `json:"stream_id,omitempty"`
	Status         EventStatus  `json:"status"`
	CoverImageURL  string       `json:"cover_image_url,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

type EventTicketTier struct {
	ID            uuid.UUID  `json:"id"`
	EventID       uuid.UUID  `json:"event_id"`
	Name          string     `json:"name"`
	PriceCents    int64      `json:"price_cents"`
	QuantityTotal int        `json:"quantity_total"`
	QuantitySold  int        `json:"quantity_sold"`
	Benefits      []string   `json:"benefits,omitempty"`
	SalesStartAt  *time.Time `json:"sales_start_at,omitempty"`
	SalesEndAt    *time.Time `json:"sales_end_at,omitempty"`
}

type EventAttendee struct {
	EventID      uuid.UUID       `json:"event_id"`
	TicketTierID *uuid.UUID      `json:"ticket_tier_id,omitempty"`
	UserID       uuid.UUID       `json:"user_id"`
	Status       AttendeeStatus  `json:"status"`
	RegisteredAt time.Time       `json:"registered_at"`
	AttendedAt   *time.Time      `json:"attended_at,omitempty"`
}

type EventSession struct {
	ID        uuid.UUID  `json:"id"`
	EventID   uuid.UUID  `json:"event_id"`
	Title     string     `json:"title"`
	Speaker   string     `json:"speaker"`
	StartAt   time.Time  `json:"start_at"`
	EndAt     time.Time  `json:"end_at"`
	StreamID  *uuid.UUID `json:"stream_id,omitempty"`
}

type EventSeries struct {
	ID           uuid.UUID       `json:"id"`
	CreatorID    uuid.UUID       `json:"creator_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Frequency    SeriesFrequency `json:"frequency"`
	DayOfWeek    *int            `json:"day_of_week,omitempty"`
	StartTime    string          `json:"start_time"`
	Timezone     string          `json:"timezone"`
	NextEventAt  *time.Time      `json:"next_event_at,omitempty"`
	IsActive     bool            `json:"is_active"`
}

type EventAdminStats struct {
	TotalEvents      int     `json:"total_events"`
	PublishedEvents  int     `json:"published_events"`
	TotalAttendees   int     `json:"total_attendees"`
	RevenueCents     int64   `json:"revenue_cents"`
}
