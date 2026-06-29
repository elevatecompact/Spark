package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type EventType string

const (
	CreatorCreated  EventType = "creator.created"
	CreatorUpdated  EventType = "creator.updated"
	CreatorVerified EventType = "creator.verified"
	CreatorFollowed EventType = "creator.followed"
	CreatorDeleted  EventType = "creator.deleted"
)

type CloudEvent struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	SpecVersion     string      `json:"specversion"`
	Type            string      `json:"type"`
	Subject         string      `json:"subject"`
	Time            string      `json:"time"`
	DataContentType string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	return &Producer{writer: w}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) emit(ctx context.Context, eventType EventType, subject string, data interface{}) error {
	event := CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.creator-service",
		SpecVersion:     "1.0",
		Type:            string(eventType),
		Subject:         subject,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(subject),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(eventType)},
			{Key: "source", Value: []byte("spark.creator-service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	log.Info().Str("type", string(eventType)).Str("subject", subject).Msg("Event emitted")
	return nil
}

func (p *Producer) CreatorCreated(ctx context.Context, creatorID, userID, displayName string) error {
	return p.emit(ctx, CreatorCreated, creatorID, map[string]interface{}{
		"creator_id":   creatorID,
		"user_id":      userID,
		"display_name": displayName,
	})
}

func (p *Producer) CreatorUpdated(ctx context.Context, creatorID string, changes map[string]interface{}) error {
	return p.emit(ctx, CreatorUpdated, creatorID, map[string]interface{}{
		"creator_id": creatorID,
		"changes":    changes,
	})
}

func (p *Producer) CreatorVerified(ctx context.Context, creatorID, adminID string) error {
	return p.emit(ctx, CreatorVerified, creatorID, map[string]interface{}{
		"creator_id": creatorID,
		"admin_id":   adminID,
		"verified":   true,
	})
}

func (p *Producer) CreatorFollowed(ctx context.Context, followerID, creatorID string) error {
	return p.emit(ctx, CreatorFollowed, creatorID, map[string]interface{}{
		"follower_id": followerID,
		"creator_id":  creatorID,
	})
}

func (p *Producer) CreatorDeleted(ctx context.Context, creatorID string) error {
	return p.emit(ctx, CreatorDeleted, creatorID, map[string]interface{}{
		"creator_id": creatorID,
	})
}
