package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type EventStartedEvent struct {
	EventID    uuid.UUID `json:"eventId"`
	CreatorID  uuid.UUID `json:"creatorId"`
	StreamID   uuid.UUID `json:"streamId"`
	TicketCount int      `json:"ticketCount"`
	StartedAt  time.Time `json:"startedAt"`
}

type EventProducer interface {
	PublishEventCreated(ctx context.Context, eventID uuid.UUID) error
	PublishEventCancelled(ctx context.Context, eventID uuid.UUID) error
	PublishEventStarted(ctx context.Context, e *EventStartedEvent) error
	PublishTicketPurchased(ctx context.Context, ticketTierID, userID uuid.UUID) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "event-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishEventCreated(ctx context.Context, eventID uuid.UUID) error {
	log.Debug().Msg("noop: event.created")
	return nil
}
func (p *noopProducer) PublishEventCancelled(ctx context.Context, eventID uuid.UUID) error {
	log.Debug().Msg("noop: event.cancelled")
	return nil
}
func (p *noopProducer) PublishEventStarted(ctx context.Context, e *EventStartedEvent) error {
	log.Debug().Msg("noop: event.started")
	return nil
}
func (p *noopProducer) PublishTicketPurchased(ctx context.Context, ticketTierID, userID uuid.UUID) error {
	log.Debug().Msg("noop: event.ticket.purchased")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishEventCreated(ctx context.Context, eventID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishEventCancelled(ctx context.Context, eventID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishEventStarted(ctx context.Context, e *EventStartedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishTicketPurchased(ctx context.Context, ticketTierID, userID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
