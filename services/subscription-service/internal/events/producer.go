package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
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
	PublishActivated(ctx context.Context, sub *domain.Subscription) error
	PublishCancelled(ctx context.Context, sub *domain.Subscription) error
	PublishUpgraded(ctx context.Context, sub *domain.Subscription) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "subscription-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    100,
			Async:        false,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishActivated(ctx context.Context, sub *domain.Subscription) error {
	log.Debug().Str("sub_id", sub.ID.String()).Msg("noop: subscription.activated")
	return nil
}
func (p *noopProducer) PublishCancelled(ctx context.Context, sub *domain.Subscription) error {
	log.Debug().Str("sub_id", sub.ID.String()).Msg("noop: subscription.cancelled")
	return nil
}
func (p *noopProducer) PublishUpgraded(ctx context.Context, sub *domain.Subscription) error {
	log.Debug().Str("sub_id", sub.ID.String()).Msg("noop: subscription.upgraded")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishActivated(ctx context.Context, sub *domain.Subscription) error {
	return p.publish(ctx, newCloudEvent("subscription.activated", sub))
}
func (p *kafkaProducer) PublishCancelled(ctx context.Context, sub *domain.Subscription) error {
	return p.publish(ctx, newCloudEvent("subscription.cancelled", sub))
}
func (p *kafkaProducer) PublishUpgraded(ctx context.Context, sub *domain.Subscription) error {
	return p.publish(ctx, newCloudEvent("subscription.upgraded", sub))
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }

func (p *kafkaProducer) publish(ctx context.Context, event CloudEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(event.Type)},
			{Key: "source", Value: []byte(event.Source)},
		},
	})
}

func newCloudEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.subscription-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
