package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
)

type CloudEvent struct {
	ID          string      `json:"id"`
	Source      string      `json:"source"`
	SpecVersion string      `json:"specversion"`
	Type        string      `json:"type"`
	Time        string      `json:"time"`
	DataContentType string  `json:"datacontenttype"`
	Data        interface{} `json:"data"`
}

type EventProducer interface {
	PublishUserCreated(ctx context.Context, user *domain.User) error
	PublishUserLoggedIn(ctx context.Context, session *domain.Session) error
	PublishUserLoggedOut(ctx context.Context, sessionID string) error
	PublishUserUpdated(ctx context.Context, user *domain.User) error
	PublishUserDeleted(ctx context.Context, userID string) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
	source string
}

func NewKafkaProducer(brokers []string, topic string) EventProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}

	return &kafkaProducer{
		writer: writer,
		topic:  topic,
		source: "spark.identity-service",
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishUserCreated(ctx context.Context, user *domain.User) error {
	log.Debug().Str("user_id", user.ID.String()).Msg("noop: user created event")
	return nil
}

func (p *noopProducer) PublishUserLoggedIn(ctx context.Context, session *domain.Session) error {
	log.Debug().Str("user_id", session.UserID.String()).Msg("noop: user logged in event")
	return nil
}

func (p *noopProducer) PublishUserLoggedOut(ctx context.Context, sessionID string) error {
	log.Debug().Str("session_id", sessionID).Msg("noop: user logged out event")
	return nil
}

func (p *noopProducer) PublishUserUpdated(ctx context.Context, user *domain.User) error {
	log.Debug().Str("user_id", user.ID.String()).Msg("noop: user updated event")
	return nil
}

func (p *noopProducer) PublishUserDeleted(ctx context.Context, userID string) error {
	log.Debug().Str("user_id", userID).Msg("noop: user deleted event")
	return nil
}

func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishUserCreated(ctx context.Context, user *domain.User) error {
	event := newCloudEvent("user.created", map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"role":         user.Role,
		"status":       user.Status,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishUserLoggedIn(ctx context.Context, session *domain.Session) error {
	event := newCloudEvent("user.logged_in", map[string]interface{}{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"ip_address": session.IPAddress,
		"user_agent": session.UserAgent,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishUserLoggedOut(ctx context.Context, sessionID string) error {
	event := newCloudEvent("user.logged_out", map[string]interface{}{
		"session_id": sessionID,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishUserUpdated(ctx context.Context, user *domain.User) error {
	event := newCloudEvent("user.updated", map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"role":         user.Role,
		"status":       user.Status,
		"verified":     user.Verified,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishUserDeleted(ctx context.Context, userID string) error {
	event := newCloudEvent("user.deleted", map[string]interface{}{
		"user_id": userID,
	})
	return p.publish(ctx, event)
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

	log.Debug().
		Str("event_id", event.ID).
		Str("type", event.Type).
		Msg("event published to kafka")

	return nil
}

func newCloudEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.identity-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
