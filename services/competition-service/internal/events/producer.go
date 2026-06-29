package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type MatchCompletedEvent struct {
	CompetitionID  uuid.UUID              `json:"competitionId"`
	MatchID        uuid.UUID              `json:"matchId"`
	WinnerID       uuid.UUID              `json:"winnerId"`
	LoserID        uuid.UUID              `json:"loserId"`
	Scores         map[string]interface{} `json:"scores"`
	Round          int                    `json:"round"`
	BracketPosition int                   `json:"bracketPosition"`
	CompletedAt    time.Time              `json:"completedAt"`
}

type EventProducer interface {
	PublishCompetitionCreated(ctx context.Context, compID uuid.UUID) error
	PublishCompetitionStarted(ctx context.Context, compID uuid.UUID) error
	PublishCompetitionEnded(ctx context.Context, compID uuid.UUID) error
	PublishParticipantRegistered(ctx context.Context, compID, userID uuid.UUID) error
	PublishMatchCompleted(ctx context.Context, e *MatchCompletedEvent) error
	PublishPrizeDistributed(ctx context.Context, compID uuid.UUID) error
	PublishLeaderboardUpdated(ctx context.Context, compID uuid.UUID) error
	PublishMatchDisputed(ctx context.Context, matchID uuid.UUID) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "competition-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishCompetitionCreated(ctx context.Context, compID uuid.UUID) error {
	log.Debug().Msg("noop: competition.created")
	return nil
}
func (p *noopProducer) PublishCompetitionStarted(ctx context.Context, compID uuid.UUID) error {
	log.Debug().Msg("noop: competition.started")
	return nil
}
func (p *noopProducer) PublishCompetitionEnded(ctx context.Context, compID uuid.UUID) error {
	log.Debug().Msg("noop: competition.ended")
	return nil
}
func (p *noopProducer) PublishParticipantRegistered(ctx context.Context, compID, userID uuid.UUID) error {
	log.Debug().Msg("noop: competition.participant.registered")
	return nil
}
func (p *noopProducer) PublishMatchCompleted(ctx context.Context, e *MatchCompletedEvent) error {
	log.Debug().Msg("noop: competition.match.completed")
	return nil
}
func (p *noopProducer) PublishPrizeDistributed(ctx context.Context, compID uuid.UUID) error {
	log.Debug().Msg("noop: competition.prize.distributed")
	return nil
}
func (p *noopProducer) PublishLeaderboardUpdated(ctx context.Context, compID uuid.UUID) error {
	log.Debug().Msg("noop: competition.leaderboard.updated")
	return nil
}
func (p *noopProducer) PublishMatchDisputed(ctx context.Context, matchID uuid.UUID) error {
	log.Debug().Msg("noop: competition.match.disputed")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishCompetitionCreated(ctx context.Context, compID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishCompetitionStarted(ctx context.Context, compID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishCompetitionEnded(ctx context.Context, compID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishParticipantRegistered(ctx context.Context, compID, userID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishMatchCompleted(ctx context.Context, e *MatchCompletedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishPrizeDistributed(ctx context.Context, compID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishLeaderboardUpdated(ctx context.Context, compID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishMatchDisputed(ctx context.Context, matchID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
