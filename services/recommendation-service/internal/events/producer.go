package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/recommendation-service/internal/domain"
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
	PublishFeedServed(ctx context.Context, feed *domain.Feed) error
	PublishFeedback(ctx context.Context, userID, contentID uuid.UUID, typ string) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "recommendation-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishFeedServed(ctx context.Context, feed *domain.Feed) error {
	log.Debug().Str("type", string(feed.Type)).Msg("noop: recommendation.feed.served")
	return nil
}
func (p *noopProducer) PublishFeedback(ctx context.Context, userID, contentID uuid.UUID, typ string) error {
	log.Debug().Str("type", typ).Msg("noop: recommendation.feedback")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishFeedServed(ctx context.Context, feed *domain.Feed) error {
	event := CloudEvent{
		ID:   uuid.New().String(),
		Source: "spark.recommendation-service",
		SpecVersion: "1.0",
		Type: "recommendation.feed.served",
		Time: time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            feed,
	}
	data, _ := json.Marshal(event)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(event.ID), Value: data})
}
func (p *kafkaProducer) PublishFeedback(ctx context.Context, userID, contentID uuid.UUID, typ string) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
