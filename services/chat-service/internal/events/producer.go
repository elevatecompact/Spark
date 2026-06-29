package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
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
	PublishRoomCreated(ctx context.Context, room *domain.ChatRoom) error
	PublishRoomClosed(ctx context.Context, room *domain.ChatRoom) error
	PublishMessageSent(ctx context.Context, msg *domain.ChatMessage) error
	PublishMessageDeleted(ctx context.Context, msg *domain.ChatMessage) error
	PublishUserMuted(ctx context.Context, roomID, userID uuid.UUID) error
	PublishUserBanned(ctx context.Context, roomID, userID uuid.UUID) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
	source string
}

func NewKafkaProducer(brokers []string, topic string) EventProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}
	return &kafkaProducer{
		writer: writer,
		topic:  topic,
		source: "spark.chat-service",
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishRoomCreated(ctx context.Context, room *domain.ChatRoom) error {
	log.Debug().Str("room_id", room.ID.String()).Msg("noop: room created")
	return nil
}
func (p *noopProducer) PublishRoomClosed(ctx context.Context, room *domain.ChatRoom) error {
	log.Debug().Str("room_id", room.ID.String()).Msg("noop: room closed")
	return nil
}
func (p *noopProducer) PublishMessageSent(ctx context.Context, msg *domain.ChatMessage) error {
	log.Debug().Str("msg_id", msg.ID.String()).Msg("noop: message sent")
	return nil
}
func (p *noopProducer) PublishMessageDeleted(ctx context.Context, msg *domain.ChatMessage) error {
	log.Debug().Str("msg_id", msg.ID.String()).Msg("noop: message deleted")
	return nil
}
func (p *noopProducer) PublishUserMuted(ctx context.Context, roomID, userID uuid.UUID) error {
	log.Debug().Str("user_id", userID.String()).Msg("noop: user muted")
	return nil
}
func (p *noopProducer) PublishUserBanned(ctx context.Context, roomID, userID uuid.UUID) error {
	log.Debug().Str("user_id", userID.String()).Msg("noop: user banned")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishRoomCreated(ctx context.Context, room *domain.ChatRoom) error {
	return p.publish(ctx, newCloudEvent("chat.room.created", room))
}
func (p *kafkaProducer) PublishRoomClosed(ctx context.Context, room *domain.ChatRoom) error {
	return p.publish(ctx, newCloudEvent("chat.room.closed", room))
}
func (p *kafkaProducer) PublishMessageSent(ctx context.Context, msg *domain.ChatMessage) error {
	return p.publish(ctx, newCloudEvent("chat.message.sent", msg))
}
func (p *kafkaProducer) PublishMessageDeleted(ctx context.Context, msg *domain.ChatMessage) error {
	return p.publish(ctx, newCloudEvent("chat.message.deleted", msg))
}
func (p *kafkaProducer) PublishUserMuted(ctx context.Context, roomID, userID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent("chat.user.muted", map[string]string{
		"room_id": roomID.String(),
		"user_id": userID.String(),
	}))
}
func (p *kafkaProducer) PublishUserBanned(ctx context.Context, roomID, userID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent("chat.user.banned", map[string]string{
		"room_id": roomID.String(),
		"user_id": userID.String(),
	}))
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}

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
	log.Debug().Str("event_id", event.ID).Str("type", event.Type).Msg("event published to kafka")
	return nil
}

func newCloudEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.chat-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
