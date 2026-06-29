package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type ImpressionRecordedEvent struct {
	ImpressionID  uuid.UUID `json:"impressionId"`
	CampaignID    uuid.UUID `json:"campaignId"`
	AdUnitID      uuid.UUID `json:"adUnitId"`
	PlacementID   string    `json:"placementId"`
	UserID        *uuid.UUID `json:"userId"`
	CostCents     int64     `json:"costCents"`
	ServedAt      time.Time `json:"servedAt"`
}

type EventProducer interface {
	PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error
	PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error
	PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error
	PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error
	PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "advertising-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Msg("noop: advertising.campaign.created")
	return nil
}
func (p *noopProducer) PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Msg("noop: advertising.campaign.activated")
	return nil
}
func (p *noopProducer) PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Msg("noop: advertising.campaign.ended")
	return nil
}
func (p *noopProducer) PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error {
	log.Debug().Msg("noop: advertising.impression.recorded")
	return nil
}
func (p *noopProducer) PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error {
	log.Debug().Msg("noop: advertising.click.recorded")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error { return nil }
func (p *kafkaProducer) PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error { return nil }
func (p *kafkaProducer) PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error { return nil }
func (p *kafkaProducer) PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error { return nil }
func (p *kafkaProducer) PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error { return nil }
func (p *kafkaProducer) Close() error { return p.writer.Close() }
