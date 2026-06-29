package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type CloudEvent struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	SpecVersion     string      `json:"specversion"`
	Type            string      `json:"type"`
	Time            string      `json:"time"`
	DataContentType string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
}

type EventProducer interface {
	PublishDelivered(ctx context.Context, notifID, userID uuid.UUID, channel, status string) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "notification-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishDelivered(ctx context.Context, notifID, userID uuid.UUID, channel, status string) error {
	log.Debug().Str("notif_id", notifID.String()).Msg("noop: notification.delivered")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishDelivered(ctx context.Context, notifID, userID uuid.UUID, channel, status string) error {
	return p.publish(ctx, newEvent("notification.delivered", map[string]interface{}{
		"notification_id": notifID,
		"user_id":         userID,
		"channel":         channel,
		"status":          status,
	}))
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }

func (p *kafkaProducer) publish(ctx context.Context, event CloudEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(event.Type)},
		},
	})
}

func newEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.notification-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
