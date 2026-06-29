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
	ID          string      `json:"id"`
	Source      string      `json:"source"`
	SpecVersion string      `json:"specversion"`
	Type        string      `json:"type"`
	Time        string      `json:"time"`
	Data        interface{} `json:"data"`
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "discovery-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, eventType string, data interface{}) {
	event := CloudEvent{
		ID:          uuid.New().String(),
		Source:      "discovery-service",
		SpecVersion: "1.0",
		Type:        eventType,
		Time:        time.Now().UTC().Format(time.RFC3339),
		Data:        data,
	}
	body, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Str("type", eventType).Msg("failed to marshal event")
		return
	}
	if err := p.writer.WriteMessages(ctx, kafka.Message{Value: body}); err != nil {
		log.Error().Err(err).Str("type", eventType).Msg("failed to publish event")
		return
	}
	log.Info().Str("type", eventType).Msg("event published")
}

func (p *KafkaProducer) Close() {
	if p.writer != nil {
		p.writer.Close()
	}
}

type NoopProducer struct{}

func NewNoopProducer() *NoopProducer {
	return &NoopProducer{}
}

func (p *NoopProducer) Publish(ctx context.Context, eventType string, data interface{}) {
	log.Info().Str("type", eventType).Msg("noop event")
}

func (p *NoopProducer) Close() {}
