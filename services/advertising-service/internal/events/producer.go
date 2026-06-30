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

type ImpressionRecordedEvent struct {
	ImpressionID uuid.UUID `json:"impressionId"`
	CampaignID   uuid.UUID `json:"campaignId"`
	AdUnitID     uuid.UUID `json:"adUnitID"`
	PlacementID  string    `json:"placementId"`
	UserID       *uuid.UUID `json:"userId"`
	CostCents    int64     `json:"costCents"`
	ServedAt     time.Time `json:"servedAt"`
}

type EventProducer interface {
	PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error
	PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error
	PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error
	PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error
	PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	source string
}

type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	if len(brokers) == 0 {
		return NewNoopProducer()
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        "advertising-events",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}
	return &kafkaProducer{
		writer: writer,
		source: "spark.advertising-service",
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Str("campaign_id", campaignID.String()).Msg("noop: advertising.campaign.created")
	return nil
}
func (p *noopProducer) PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Str("campaign_id", campaignID.String()).Msg("noop: advertising.campaign.activated")
	return nil
}
func (p *noopProducer) PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error {
	log.Debug().Str("campaign_id", campaignID.String()).Msg("noop: advertising.campaign.ended")
	return nil
}
func (p *noopProducer) PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error {
	log.Debug().Str("impression_id", e.ImpressionID.String()).Msg("noop: advertising.impression.recorded")
	return nil
}
func (p *noopProducer) PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error {
	log.Debug().Str("impression_id", impressionID.String()).Msg("noop: advertising.click.recorded")
	return nil
}
func (p *noopProducer) Close() error { return nil }

type CloudEvent struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	SpecVersion     string      `json:"specversion"`
	Type            string      `json:"type"`
	Time            string      `json:"time"`
	DataContentType string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
}

func (p *kafkaProducer) publish(ctx context.Context, event CloudEvent) error {
	if p == nil || p.writer == nil {
		return fmt.Errorf("kafka producer not initialised")
	}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal cloud event: %w", err)
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

func newCloudEvent(source, eventType string, data interface{}) CloudEvent {
	return CloudEvent{
		ID:              uuid.New().String(),
		Source:          source,
		SpecVersion:     "1.0",
		Type:            eventType,
		Time:            time.Now().UTC().Format(time.RFC3339),
		DataContentType: "application/json",
		Data:            data,
	}
}

func (p *kafkaProducer) PublishCampaignCreated(ctx context.Context, campaignID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent(p.source, "advertising.campaign.created", map[string]interface{}{
		"campaignId": campaignID,
	}))
}

func (p *kafkaProducer) PublishCampaignActivated(ctx context.Context, campaignID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent(p.source, "advertising.campaign.activated", map[string]interface{}{
		"campaignId": campaignID,
	}))
}

func (p *kafkaProducer) PublishCampaignEnded(ctx context.Context, campaignID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent(p.source, "advertising.campaign.ended", map[string]interface{}{
		"campaignId": campaignID,
	}))
}

func (p *kafkaProducer) PublishImpressionRecorded(ctx context.Context, e *ImpressionRecordedEvent) error {
	return p.publish(ctx, newCloudEvent(p.source, "advertising.impression.recorded", e))
}

func (p *kafkaProducer) PublishClickRecorded(ctx context.Context, impressionID uuid.UUID) error {
	return p.publish(ctx, newCloudEvent(p.source, "advertising.click.recorded", map[string]interface{}{
		"impressionId": impressionID,
	}))
}

func (p *kafkaProducer) Close() error {
	if p == nil || p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
