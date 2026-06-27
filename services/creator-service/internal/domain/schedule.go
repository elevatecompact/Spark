package domain

import (
	"time"
	"github.com/google/uuid"
)

type ScheduleSlot struct {
	ID        uuid.UUID json:"id"
	CreatorID uuid.UUID json:"creator_id"
	DayOfWeek int       json:"day_of_week"
	StartTime string    json:"start_time"
	EndTime   string    json:"end_time"
	Title     string    json:"title,omitempty"
	Recurring bool      json:"recurring"
	Active    bool      json:"active"
	CreatedAt time.Time json:"created_at"
}

type CreateScheduleSlotRequest struct {
	DayOfWeek int    json:"day_of_week" validate:"required,min=0,max=6"
	StartTime string json:"start_time" validate:"required,pattern=^\\d{2}:\\d{2}$"
	EndTime   string json:"end_time" validate:"required,pattern=^\\d{2}:\\d{2}$"
	Title     string json:"title,omitempty"
	Recurring bool   json:"recurring"
}
