package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one broker is required")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{writer: w}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.Subject),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "event_id", Value: []byte(event.ID)},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (p *Producer) PublishBatch(ctx context.Context, topic string, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	messages := make([]kafka.Message, 0, len(events))
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}
		messages = append(messages, kafka.Message{
			Topic: topic,
			Key:   []byte(event.Subject),
			Value: data,
			Headers: []kafka.Header{
				{Key: "event_type", Value: []byte(event.Type)},
				{Key: "event_id", Value: []byte(event.ID)},
			},
		})
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("failed to publish events batch: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
