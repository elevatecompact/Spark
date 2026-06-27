package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type EventHandler func(ctx context.Context, event *Event) error

type Consumer struct {
	brokers []string
	reader  *kafka.Reader
	groupID string
	topics  []string
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one broker is required")
	}
	if groupID == "" {
		return nil, fmt.Errorf("consumer group ID is required")
	}

	return &Consumer{
		brokers: brokers,
		groupID: groupID,
	}, nil
}

func (c *Consumer) Subscribe(ctx context.Context, topics []string, handler EventHandler) error {
	if len(topics) == 0 {
		return fmt.Errorf("at least one topic is required")
	}
	if handler == nil {
		return fmt.Errorf("event handler cannot be nil")
	}

	c.topics = topics

	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:                c.readerBrokers(),
		GroupID:                c.groupID,
		GroupTopics:            topics,
		MinBytes:               10e3,
		MaxBytes:               10e6,
		MaxWait:                1 * time.Second,
		CommitInterval:         time.Second,
		PartitionWatchInterval: 5 * time.Second,
		HeartbeatInterval:      3 * time.Second,
		SessionTimeout:         30 * time.Second,
		RebalanceTimeout:       30 * time.Second,
		StartOffset:            kafka.LastOffset,
		RetentionTime:          24 * time.Hour,
	})

	c.wg.Add(1)
	go c.consumeLoop(ctx, handler)

	return nil
}

func (c *Consumer) readerBrokers() []string {
	return c.brokers
}

func (c *Consumer) consumeLoop(ctx context.Context, handler EventHandler) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
			log.Printf("kafka read error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			continue
		}

		if err := handler(ctx, &event); err != nil {
			log.Printf("event handler error: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()

	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
