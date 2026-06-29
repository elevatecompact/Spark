package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type TranslationCompletedEvent struct {
	TranslationID  uuid.UUID `json:"translationId"`
	SourceText     string    `json:"sourceText"`
	TranslatedText string    `json:"translatedText"`
	SourceLang     string    `json:"sourceLang"`
	TargetLang     string    `json:"targetLang"`
	Provider       string    `json:"provider"`
	LatencyMs      int64     `json:"latencyMs"`
	CharCount      int       `json:"charCount"`
	Timestamp      time.Time `json:"timestamp"`
}

type EventProducer interface {
	PublishTranslationCompleted(ctx context.Context, e *TranslationCompletedEvent) error
	PublishBatchCompleted(ctx context.Context, jobID uuid.UUID, count int) error
	PublishReviewReady(ctx context.Context, entryID uuid.UUID) error
	PublishProviderSwitched(ctx context.Context, provider string) error
	PublishMemoryUpdated(ctx context.Context, entryID uuid.UUID) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "translation-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishTranslationCompleted(ctx context.Context, e *TranslationCompletedEvent) error {
	log.Debug().Str("lang", e.SourceLang+"->"+e.TargetLang).Msg("noop: translation.completed")
	return nil
}
func (p *noopProducer) PublishBatchCompleted(ctx context.Context, jobID uuid.UUID, count int) error {
	log.Debug().Int("count", count).Msg("noop: translation.batch.completed")
	return nil
}
func (p *noopProducer) PublishReviewReady(ctx context.Context, entryID uuid.UUID) error {
	log.Debug().Msg("noop: translation.review.ready")
	return nil
}
func (p *noopProducer) PublishProviderSwitched(ctx context.Context, provider string) error {
	log.Debug().Str("provider", provider).Msg("noop: translation.provider.switched")
	return nil
}
func (p *noopProducer) PublishMemoryUpdated(ctx context.Context, entryID uuid.UUID) error {
	log.Debug().Msg("noop: translation.memory.updated")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishTranslationCompleted(ctx context.Context, e *TranslationCompletedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishBatchCompleted(ctx context.Context, jobID uuid.UUID, count int) error {
	return nil
}
func (p *kafkaProducer) PublishReviewReady(ctx context.Context, entryID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishProviderSwitched(ctx context.Context, provider string) error {
	return nil
}
func (p *kafkaProducer) PublishMemoryUpdated(ctx context.Context, entryID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
