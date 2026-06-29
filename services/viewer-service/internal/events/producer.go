package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
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
	PublishWatchStarted(ctx context.Context, entry *domain.WatchHistory) error
	PublishWatchProgress(ctx context.Context, entry *domain.WatchHistory) error
	PublishWatchCompleted(ctx context.Context, entry *domain.WatchHistory) error
	PublishRatingSubmitted(ctx context.Context, rating *domain.Rating) error
	PublishReactionAdded(ctx context.Context, reaction *domain.Reaction) error
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
		source: "spark.viewer-service",
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishWatchStarted(ctx context.Context, entry *domain.WatchHistory) error {
	log.Debug().Str("content_id", entry.ContentID.String()).Msg("noop: watch started event")
	return nil
}

func (p *noopProducer) PublishWatchProgress(ctx context.Context, entry *domain.WatchHistory) error {
	log.Debug().Str("content_id", entry.ContentID.String()).Float64("progress", entry.Progress).Msg("noop: watch progress event")
	return nil
}

func (p *noopProducer) PublishWatchCompleted(ctx context.Context, entry *domain.WatchHistory) error {
	log.Debug().Str("content_id", entry.ContentID.String()).Msg("noop: watch completed event")
	return nil
}

func (p *noopProducer) PublishRatingSubmitted(ctx context.Context, rating *domain.Rating) error {
	log.Debug().Str("content_id", rating.ContentID.String()).Int("score", rating.Score).Msg("noop: rating submitted event")
	return nil
}

func (p *noopProducer) PublishReactionAdded(ctx context.Context, reaction *domain.Reaction) error {
	log.Debug().Str("content_id", reaction.ContentID.String()).Str("type", string(reaction.Type)).Msg("noop: reaction added event")
	return nil
}

func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishWatchStarted(ctx context.Context, entry *domain.WatchHistory) error {
	event := newCloudEvent("viewer.watch.started", map[string]interface{}{
		"viewer_id":  entry.ViewerID,
		"content_id": entry.ContentID,
		"type":      entry.ContentType,
		"timestamp": entry.WatchedAt,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishWatchProgress(ctx context.Context, entry *domain.WatchHistory) error {
	event := newCloudEvent("viewer.watch.progress", map[string]interface{}{
		"viewer_id":            entry.ViewerID,
		"content_id":           entry.ContentID,
		"content_type":         entry.ContentType,
		"progress":            entry.Progress,
		"watch_duration_seconds": entry.WatchDurationSeconds,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishWatchCompleted(ctx context.Context, entry *domain.WatchHistory) error {
	event := newCloudEvent("viewer.watch.completed", map[string]interface{}{
		"viewer_id":            entry.ViewerID,
		"content_id":           entry.ContentID,
		"content_type":         entry.ContentType,
		"watch_duration_seconds": entry.WatchDurationSeconds,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishRatingSubmitted(ctx context.Context, rating *domain.Rating) error {
	event := newCloudEvent("viewer.rating.submitted", map[string]interface{}{
		"viewer_id":  rating.ViewerID,
		"content_id": rating.ContentID,
		"score":     rating.Score,
	})
	return p.publish(ctx, event)
}

func (p *kafkaProducer) PublishReactionAdded(ctx context.Context, reaction *domain.Reaction) error {
	event := newCloudEvent("viewer.reaction.added", map[string]interface{}{
		"viewer_id":  reaction.ViewerID,
		"content_id": reaction.ContentID,
		"type":      reaction.Type,
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
		Source:          "spark.viewer-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
