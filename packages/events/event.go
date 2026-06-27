package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	Subject         string      `json:"subject"`
	Data            interface{} `json:"data"`
	Time            time.Time   `json:"time"`
	SpecVersion     string      `json:"specversion"`
	DataContentType string      `json:"datacontenttype"`
}

func New(source, eventType, subject string, data interface{}) *Event {
	return &Event{
		ID:              uuid.New().String(),
		Source:          source,
		Type:            eventType,
		Subject:         subject,
		Data:            data,
		Time:            time.Now().UTC(),
		SpecVersion:     "1.0",
		DataContentType: "application/json",
	}
}
