package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
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
	PublishSent(ctx context.Context, gift *domain.Gift) error
	PublishReceived(ctx context.Context, gift *domain.Gift) error
	PublishSubscriptionGifted(ctx context.Context, senderID, recipientID, planID uuid.UUID) error
	PublishCampaignMatch(ctx context.Context, gift *domain.Gift, campaign *domain.GiftCampaign, matchAmount int64) error
	PublishCardPurchased(ctx context.Context, card *domain.GiftCard) error
	PublishCardRedeemed(ctx context.Context, card *domain.GiftCard) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "gift-events",
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

func (p *noopProducer) PublishSent(ctx context.Context, gift *domain.Gift) error {
	log.Debug().Str("gift_id", gift.ID.String()).Msg("noop: gift.sent")
	return nil
}
func (p *noopProducer) PublishReceived(ctx context.Context, gift *domain.Gift) error {
	log.Debug().Str("gift_id", gift.ID.String()).Msg("noop: gift.received")
	return nil
}
func (p *noopProducer) PublishSubscriptionGifted(ctx context.Context, senderID, recipientID, planID uuid.UUID) error {
	log.Debug().Str("sender", senderID.String()).Msg("noop: gift.subscription.gifted")
	return nil
}
func (p *noopProducer) PublishCampaignMatch(ctx context.Context, gift *domain.Gift, campaign *domain.GiftCampaign, matchAmount int64) error {
	log.Debug().Str("gift_id", gift.ID.String()).Msg("noop: gift.campaign.match")
	return nil
}
func (p *noopProducer) PublishCardPurchased(ctx context.Context, card *domain.GiftCard) error {
	log.Debug().Str("card_id", card.ID.String()).Msg("noop: gift.card.purchased")
	return nil
}
func (p *noopProducer) PublishCardRedeemed(ctx context.Context, card *domain.GiftCard) error {
	log.Debug().Str("card_id", card.ID.String()).Msg("noop: gift.card.redeemed")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishSent(ctx context.Context, gift *domain.Gift) error {
	return p.publish(ctx, newEvent("gift.sent", gift))
}
func (p *kafkaProducer) PublishReceived(ctx context.Context, gift *domain.Gift) error {
	return p.publish(ctx, newEvent("gift.received", gift))
}
func (p *kafkaProducer) PublishSubscriptionGifted(ctx context.Context, senderID, recipientID, planID uuid.UUID) error {
	return p.publish(ctx, newEvent("gift.subscription.gifted", map[string]interface{}{
		"sender_id":    senderID,
		"recipient_id": recipientID,
		"plan_id":      planID,
	}))
}
func (p *kafkaProducer) PublishCampaignMatch(ctx context.Context, gift *domain.Gift, campaign *domain.GiftCampaign, matchAmount int64) error {
	return p.publish(ctx, newEvent("gift.campaign.match", map[string]interface{}{
		"gift_id":      gift.ID,
		"campaign_id":  campaign.ID,
		"match_amount": matchAmount,
	}))
}
func (p *kafkaProducer) PublishCardPurchased(ctx context.Context, card *domain.GiftCard) error {
	return p.publish(ctx, newEvent("gift.card.purchased", card))
}
func (p *kafkaProducer) PublishCardRedeemed(ctx context.Context, card *domain.GiftCard) error {
	return p.publish(ctx, newEvent("gift.card.redeemed", card))
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
		Source:          "spark.gift-service",
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}
