package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type EventType string

const (
	EventStreamCreated       EventType = "stream.created"
	EventStreamStarted       EventType = "stream.started"
	EventStreamEnded         EventType = "stream.ended"
	EventStreamUpdated       EventType = "stream.updated"
	EventStreamError         EventType = "stream.error"
	EventViewerJoined        EventType = "stream.viewer.joined"
	EventViewerLeft          EventType = "stream.viewer.left"
	EventStreamHealthChanged EventType = "stream.health.changed"
	EventRecordingStarted    EventType = "stream.recording.started"
	EventRecordingCompleted  EventType = "stream.recording.completed"
	EventRecordingFailed     EventType = "stream.recording.failed"
)

type StreamEvent struct {
	ID        uuid.UUID       `json:"id"`
	Type      EventType       `json:"type"`
	StreamID  uuid.UUID       `json:"stream_id"`
	CreatorID uuid.UUID       `json:"creator_id"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data,omitempty"`
}

type EventProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewEventProducer(brokers []string, topic string) *EventProducer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           10 * time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}

	return &EventProducer{writer: w, topic: topic}
}

func (p *EventProducer) Close() error {
	return p.writer.Close()
}

func (p *EventProducer) publish(ctx context.Context, eventType EventType, streamID, creatorID uuid.UUID, data interface{}) error {
	var rawData json.RawMessage
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshal event data: %w", err)
		}
		rawData = b
	}

	event := StreamEvent{
		ID:        uuid.New(),
		Type:      eventType,
		StreamID:  streamID,
		CreatorID: creatorID,
		Timestamp: time.Now().UTC(),
		Data:      rawData,
	}

	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(streamID.String()),
		Value: b,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(eventType)},
			{Key: "stream_id", Value: []byte(streamID.String())},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Error().Err(err).Str("event_type", string(eventType)).Str("stream_id", streamID.String()).Msg("Failed to publish event")
		return fmt.Errorf("write kafka message: %w", err)
	}

	log.Debug().Str("event_type", string(eventType)).Str("stream_id", streamID.String()).Msg("Event published")
	return nil
}

func (p *EventProducer) StreamCreated(ctx context.Context, stream *domain.Stream) error {
	return p.publish(ctx, EventStreamCreated, stream.ID, stream.CreatorID, map[string]interface{}{
		"title":    stream.Title,
		"category": stream.Category,
		"status":   stream.Status,
	})
}

func (p *EventProducer) StreamStarted(ctx context.Context, stream *domain.Stream) error {
	return p.publish(ctx, EventStreamStarted, stream.ID, stream.CreatorID, map[string]interface{}{
		"title":      stream.Title,
		"started_at": stream.StartedAt,
	})
}

func (p *EventProducer) StreamEnded(ctx context.Context, stream *domain.Stream) error {
	return p.publish(ctx, EventStreamEnded, stream.ID, stream.CreatorID, map[string]interface{}{
		"title":    stream.Title,
		"duration": stream.Duration,
		"ended_at": stream.EndedAt,
	})
}

func (p *EventProducer) StreamUpdated(ctx context.Context, stream *domain.Stream, changes []string) error {
	return p.publish(ctx, EventStreamUpdated, stream.ID, stream.CreatorID, map[string]interface{}{
		"changes": changes,
	})
}

func (p *EventProducer) ViewerJoined(ctx context.Context, streamID, viewerID uuid.UUID) error {
	return p.publish(ctx, EventViewerJoined, streamID, viewerID, map[string]interface{}{
		"viewer_id": viewerID,
	})
}

func (p *EventProducer) ViewerLeft(ctx context.Context, streamID, viewerID uuid.UUID) error {
	return p.publish(ctx, EventViewerLeft, streamID, viewerID, map[string]interface{}{
		"viewer_id": viewerID,
	})
}

func (p *EventProducer) StreamHealthChanged(ctx context.Context, streamID uuid.UUID, health *domain.StreamHealth) error {
	return p.publish(ctx, EventStreamHealthChanged, streamID, uuid.Nil, health)
}

func (p *EventProducer) RecordingStarted(ctx context.Context, streamID, creatorID uuid.UUID, recordingID uuid.UUID) error {
	return p.publish(ctx, EventRecordingStarted, streamID, creatorID, map[string]interface{}{
		"recording_id": recordingID,
	})
}

func (p *EventProducer) RecordingCompleted(ctx context.Context, streamID, creatorID uuid.UUID, recordingID uuid.UUID, s3Key string) error {
	return p.publish(ctx, EventRecordingCompleted, streamID, creatorID, map[string]interface{}{
		"recording_id": recordingID,
		"s3_key":       s3Key,
	})
}

func (p *EventProducer) RecordingFailed(ctx context.Context, streamID, creatorID uuid.UUID, recordingID uuid.UUID, reason string) error {
	return p.publish(ctx, EventRecordingFailed, streamID, creatorID, map[string]interface{}{
		"recording_id": recordingID,
		"reason":       reason,
	})
}
