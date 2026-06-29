package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
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
	PublishTransactionSettled(ctx context.Context, txn *domain.Transaction) error
	PublishTransactionFailed(ctx context.Context, txn *domain.Transaction) error
	PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error
	PublishBalanceLow(ctx context.Context, walletID uuid.UUID, balance int64) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	source string
}

func NewKafkaProducer(brokers []string) EventProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        "wallet-events",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}
	return &kafkaProducer{
		writer: writer,
		source: "spark.wallet-service",
	}
}

func NewNoopProducer() EventProducer {
	return &noopProducer{}
}

type noopProducer struct{}

func (p *noopProducer) PublishTransactionSettled(ctx context.Context, txn *domain.Transaction) error {
	log.Debug().Str("txn_id", txn.ID.String()).Msg("noop: transaction settled")
	return nil
}
func (p *noopProducer) PublishTransactionFailed(ctx context.Context, txn *domain.Transaction) error {
	log.Debug().Str("txn_id", txn.ID.String()).Msg("noop: transaction failed")
	return nil
}
func (p *noopProducer) PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error {
	log.Debug().Str("payout_id", payout.ID.String()).Msg("noop: payout completed")
	return nil
}
func (p *noopProducer) PublishBalanceLow(ctx context.Context, walletID uuid.UUID, balance int64) error {
	log.Debug().Str("wallet_id", walletID.String()).Int64("balance", balance).Msg("noop: balance low")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishTransactionSettled(ctx context.Context, txn *domain.Transaction) error {
	return p.publish(ctx, newCloudEvent("wallet.transaction.settled", txn))
}
func (p *kafkaProducer) PublishTransactionFailed(ctx context.Context, txn *domain.Transaction) error {
	return p.publish(ctx, newCloudEvent("wallet.transaction.failed", txn))
}
func (p *kafkaProducer) PublishPayoutCompleted(ctx context.Context, payout *domain.Payout) error {
	return p.publish(ctx, newCloudEvent("wallet.payout.completed", payout))
}
func (p *kafkaProducer) PublishBalanceLow(ctx context.Context, walletID uuid.UUID, balance int64) error {
	return p.publish(ctx, newCloudEvent("wallet.balance.low", map[string]interface{}{
		"wallet_id": walletID.String(),
		"balance":   balance,
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
		Source:          "spark.wallet-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
