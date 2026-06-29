package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/elevatecompact/spark/services/search-service/internal/domain"
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

type SearchQueryExecutedEvent struct {
	QueryID     string            `json:"queryId"`
	Query       string            `json:"query"`
	Filters     map[string]string `json:"filters"`
	ResultCount int               `json:"resultCount"`
	TopResultIDs []string         `json:"topResultIds"`
	ClickedResultID *string       `json:"clickedResultId"`
	LatencyMs   int64             `json:"latencyMs"`
	UserID      uuid.UUID         `json:"userId"`
	Timestamp   time.Time         `json:"timestamp"`
}

type EventProducer interface {
	PublishQueryExecuted(ctx context.Context, e *SearchQueryExecutedEvent) error
	PublishIndexUpdated(ctx context.Context, contentType domain.ContentType, docID uuid.UUID) error
	PublishSuggestionClicked(ctx context.Context, userID uuid.UUID, suggestion string) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "search-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishQueryExecuted(ctx context.Context, e *SearchQueryExecutedEvent) error {
	log.Debug().Str("query", e.Query).Msg("noop: search.query.executed")
	return nil
}
func (p *noopProducer) PublishIndexUpdated(ctx context.Context, ct domain.ContentType, docID uuid.UUID) error {
	log.Debug().Str("type", string(ct)).Msg("noop: search.index.updated")
	return nil
}
func (p *noopProducer) PublishSuggestionClicked(ctx context.Context, userID uuid.UUID, suggestion string) error {
	log.Debug().Str("suggestion", suggestion).Msg("noop: search.suggestions.clicked")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishQueryExecuted(ctx context.Context, e *SearchQueryExecutedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishIndexUpdated(ctx context.Context, ct domain.ContentType, docID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishSuggestionClicked(ctx context.Context, userID uuid.UUID, suggestion string) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
