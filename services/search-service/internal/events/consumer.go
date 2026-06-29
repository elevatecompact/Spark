package events

import (
	"context"

	"github.com/rs/zerolog/log"
)

type EventConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type noopConsumer struct{}

func NewNoopConsumer() EventConsumer { return &noopConsumer{} }

func (c *noopConsumer) Start(ctx context.Context) error {
	log.Info().Msg("noop consumer started (placeholder for creator.channel.created, viewer.rating.submitted, media.content.uploaded, moderation.content.flagged, identity.user.deleted)")
	return nil
}
func (c *noopConsumer) Stop() error {
	log.Info().Msg("noop consumer stopped")
	return nil
}
