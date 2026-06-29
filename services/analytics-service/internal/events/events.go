package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/analytics-service/internal/domain"
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
	PublishReportReady(ctx context.Context, report *domain.Report) error
	Close() error
}

type EventConsumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, event CloudEvent) error) error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

type kafkaConsumer struct {
	reader *kafka.Reader
}

type noopProducer struct{}
type noopConsumer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "analytics-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    100,
			Async:        false,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewKafkaConsumer(brokers []string, groupID string) EventConsumer {
	return &kafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    "platform-events",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func NewNoopProducer() EventProducer  { return &noopProducer{} }
func NewNoopConsumer() EventConsumer { return &noopConsumer{} }

func (p *noopProducer) PublishReportReady(ctx context.Context, report *domain.Report) error {
	log.Debug().Str("report_id", report.ID.String()).Msg("noop: analytics.report.ready")
	return nil
}
func (p *noopProducer) Close() error { return nil }
func (c *noopConsumer) Consume(ctx context.Context, handler func(ctx context.Context, event CloudEvent) error) error {
	log.Debug().Msg("noop consumer: no events to consume")
	return nil
}

func (p *kafkaProducer) PublishReportReady(ctx context.Context, report *domain.Report) error {
	return p.publish(ctx, newEvent("analytics.report.ready", report))
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }

func (c *kafkaConsumer) Consume(ctx context.Context, handler func(ctx context.Context, event CloudEvent) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var event CloudEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Warn().Err(err).Msg("failed to unmarshal event")
			continue
		}
		if err := handler(ctx, event); err != nil {
			log.Warn().Err(err).Str("type", event.Type).Msg("failed to handle event")
		}
	}
}

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
		},
	})
}

func newEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.analytics-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
