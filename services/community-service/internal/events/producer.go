package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type MemberJoinedEvent struct {
	CommunityID uuid.UUID `json:"communityId"`
	UserID      uuid.UUID `json:"userId"`
	JoinedAt    time.Time `json:"joinedAt"`
	MemberCount int       `json:"memberCount"`
	Role        string    `json:"role"`
}

type EventProducer interface {
	PublishCommunityCreated(ctx context.Context, communityID uuid.UUID) error
	PublishMemberJoined(ctx context.Context, e *MemberJoinedEvent) error
	PublishMemberLeft(ctx context.Context, communityID, userID uuid.UUID) error
	PublishPostCreated(ctx context.Context, postID uuid.UUID) error
	PublishRoleChanged(ctx context.Context, communityID, userID uuid.UUID, role string) error
	Close() error
}

type kafkaProducer struct{ writer *kafka.Writer }
type noopProducer struct{}

func NewKafkaProducer(brokers []string) EventProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "community-events",
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func NewNoopProducer() EventProducer { return &noopProducer{} }

func (p *noopProducer) PublishCommunityCreated(ctx context.Context, communityID uuid.UUID) error {
	log.Debug().Str("id", communityID.String()).Msg("noop: community.created")
	return nil
}
func (p *noopProducer) PublishMemberJoined(ctx context.Context, e *MemberJoinedEvent) error {
	log.Debug().Str("role", e.Role).Msg("noop: community.member.joined")
	return nil
}
func (p *noopProducer) PublishMemberLeft(ctx context.Context, communityID, userID uuid.UUID) error {
	log.Debug().Msg("noop: community.member.left")
	return nil
}
func (p *noopProducer) PublishPostCreated(ctx context.Context, postID uuid.UUID) error {
	log.Debug().Msg("noop: community.post.created")
	return nil
}
func (p *noopProducer) PublishRoleChanged(ctx context.Context, communityID, userID uuid.UUID, role string) error {
	log.Debug().Str("role", role).Msg("noop: community.role.changed")
	return nil
}
func (p *noopProducer) Close() error { return nil }

func (p *kafkaProducer) PublishCommunityCreated(ctx context.Context, communityID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishMemberJoined(ctx context.Context, e *MemberJoinedEvent) error {
	return nil
}
func (p *kafkaProducer) PublishMemberLeft(ctx context.Context, communityID, userID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishPostCreated(ctx context.Context, postID uuid.UUID) error {
	return nil
}
func (p *kafkaProducer) PublishRoleChanged(ctx context.Context, communityID, userID uuid.UUID, role string) error {
	return nil
}
func (p *kafkaProducer) Close() error { return p.writer.Close() }
