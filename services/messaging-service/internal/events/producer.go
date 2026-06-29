package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/messaging-service/internal/domain"
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
	PublishConversationCreated(ctx context.Context, conv *domain.Conversation) error
	PublishMessageSent(ctx context.Context, msg *domain.Message) error
	PublishMessageRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	source string
}

func NewKafkaProducer(brokers []string) EventProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        "messaging-events",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}
	return &kafkaProducer{
		writer: writer,
		source: "spark.messaging-service",
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishConversationCreated(ctx context.Context, conv *domain.Conversation) error {
	log.Debug().Str("conv_id", conv.ID.String()).Msg("noop: conversation created")
	return nil
}
func (p *noopProducer) PublishMessageSent(ctx context.Context, msg *domain.Message) error {
	log.Debug().Str("msg_id", msg.ID.String()).Msg("noop: message sent")
	return nil
}
func (p *noopProducer) PublishMessageRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error {
	log.Debug().Str("msg_id", msgID.String()).Msg("noop: message read")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishConversationCreated(ctx context.Context, conv *domain.Conversation) error {
	return p.publish(ctx, newCloudEvent("messaging.conversation.created", conv))
}
func (p *kafkaProducer) PublishMessageSent(ctx context.Context, msg *domain.Message) error {
	return p.publish(ctx, newCloudEvent("messaging.message.sent", msg))
}
func (p *kafkaProducer) PublishMessageRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent("messaging.message.read", map[string]string{
		"conversation_id": convID.String(),
		"user_id":         userID.String(),
		"message_id":      msgID.String(),
	}))
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }

func (p *kafkaProducer) publish(ctx context.Context, event CloudEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(event.Type)},
			{Key: "source", Value: []byte(event.Source)},
		},
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write kafka message: %w", err)
	}
	log.Debug().Str("event_id", event.ID).Str("type", event.Type).Msg("event published")
	return nil
}

func newCloudEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.messaging-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
