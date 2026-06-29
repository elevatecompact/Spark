package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/moderation-service/internal/domain"
)

type ContentFlaggedEvent struct {
	ContentID    uuid.UUID                   `json:"contentId"`
	ContentType  string                      `json:"contentType"`
	ScanResults  []domain.ScanViolation      `json:"scanResults"`
	AutoAction   *domain.ModerationAction    `json:"autoAction"`
	Timestamp    time.Time                   `json:"timestamp"`
}

type ActionTakenEvent struct {
	Action     domain.ModerationAction `json:"action"`
	Timestamp  time.Time               `json:"timestamp"`
}

type EventProducer interface {
	PublishContentFlagged(ctx context.Context, e *ContentFlaggedEvent) error
	PublishActionTaken(ctx context.Context, e *ActionTakenEvent) error
	PublishReviewCompleted(ctx context.Context, itemID uuid.UUID, resolution string) error
	PublishRuleUpdated(ctx context.Context, ruleID uuid.UUID) error
	PublishReportSubmitted(ctx context.Context, reportID uuid.UUID) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "moderation-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishContentFlagged(ctx context.Context, e *ContentFlaggedEvent) error {
	log.Debug().Str("type", e.ContentType).Msg("noop: moderation.content.flagged")
	return nil
}
func (p *noopProducer) PublishActionTaken(ctx context.Context, e *ActionTakenEvent) error {
	log.Debug().Str("action", string(e.Action.ActionType)).Msg("noop: moderation.action.taken")
	return nil
}
func (p *noopProducer) PublishReviewCompleted(ctx context.Context, itemID uuid.UUID, resolution string) error {
	log.Debug().Msg("noop: moderation.review.completed")
	return nil
}
func (p *noopProducer) PublishRuleUpdated(ctx context.Context, ruleID uuid.UUID) error {
	log.Debug().Msg("noop: moderation.rule.updated")
	return nil
}
func (p *noopProducer) PublishReportSubmitted(ctx context.Context, reportID uuid.UUID) error {
	log.Debug().Msg("noop: moderation.report.submitted")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishContentFlagged(ctx context.Context, e *ContentFlaggedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishActionTaken(ctx context.Context, e *ActionTakenEvent) error {
	return nil
}
func (p *kafkaProducer) PublishReviewCompleted(ctx context.Context, itemID uuid.UUID, resolution string) error {
	return nil
}
func (p *kafkaProducer) PublishRuleUpdated(ctx context.Context, ruleID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishReportSubmitted(ctx context.Context, reportID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
