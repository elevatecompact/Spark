package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
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
	PublishIntentCreated(ctx context.Context, intent *domain.PaymentIntent) error
	PublishIntentSucceeded(ctx context.Context, intent *domain.PaymentIntent) error
	PublishIntentFailed(ctx context.Context, intent *domain.PaymentIntent) error
	PublishRefundProcessed(ctx context.Context, intent *domain.PaymentIntent, refundID string) error
	PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "payment-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    100,
			Async:        false,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

type noopProducer struct{}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

func (p *noopProducer) PublishIntentCreated(ctx context.Context, intent *domain.PaymentIntent) error {
	log.Debug().Str("intent_id", intent.ID.String()).Msg("noop: payment.intent.created")
	return nil
}
func (p *noopProducer) PublishIntentSucceeded(ctx context.Context, intent *domain.PaymentIntent) error {
	log.Debug().Str("intent_id", intent.ID.String()).Msg("noop: payment.intent.succeeded")
	return nil
}
func (p *noopProducer) PublishIntentFailed(ctx context.Context, intent *domain.PaymentIntent) error {
	log.Debug().Str("intent_id", intent.ID.String()).Msg("noop: payment.intent.failed")
	return nil
}
func (p *noopProducer) PublishRefundProcessed(ctx context.Context, intent *domain.PaymentIntent, refundID string) error {
	log.Debug().Str("intent_id", intent.ID.String()).Msg("noop: payment.refund.processed")
	return nil
}
func (p *noopProducer) PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error {
	log.Debug().Str("payout_id", payout.ID.String()).Msg("noop: payment.payout.completed")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishIntentCreated(ctx context.Context, intent *domain.PaymentIntent) error {
	return p.publish(ctx, newEvent("payment.intent.created", intent))
}
func (p *kafkaProducer) PublishIntentSucceeded(ctx context.Context, intent *domain.PaymentIntent) error {
	return p.publish(ctx, newEvent("payment.intent.succeeded", intent))
}
func (p *kafkaProducer) PublishIntentFailed(ctx context.Context, intent *domain.PaymentIntent) error {
	return p.publish(ctx, newEvent("payment.intent.failed", intent))
}
func (p *kafkaProducer) PublishRefundProcessed(ctx context.Context, intent *domain.PaymentIntent, refundID string) error {
	return p.publish(ctx, newEvent("payment.refund.processed", map[string]interface{}{
		"intent_id": intent.ID,
		"refund_id": refundID,
	}))
}
func (p *kafkaProducer) PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error {
	return p.publish(ctx, newEvent("payment.payout.completed", payout))
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }

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
			{Key: "source", Value: []byte(event.Source)},
		},
	})
}

func newEvent(eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          "spark.payment-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
