package processor

import (
	"context"

	"github.com/rs/zerolog/log"
)

type PushProcessor interface {
	Send(ctx context.Context, deviceToken string, title, body string, data map[string]string) error
}

type EmailProcessor interface {
	Send(ctx context.Context, to, subject, body string) error
}

type SMSProcessor interface {
	Send(ctx context.Context, to, body string) error
}

type noopPush struct{}

func NewNoopPush() PushProcessor { return &noopPush{} }

func (p *noopPush) Send(ctx context.Context, deviceToken string, title, body string, data map[string]string) error {
	log.Debug().Str("token", deviceToken[:min(8, len(deviceToken))]+"...").Str("title", title).Msg("noop push sent")
	return nil
}

type noopEmail struct{}

func NewNoopEmail() EmailProcessor { return &noopEmail{} }

func (e *noopEmail) Send(ctx context.Context, to, subject, body string) error {
	log.Debug().Str("to", to).Str("subject", subject).Msg("noop email sent")
	return nil
}

type noopSMS struct{}

func NewNoopSMS() SMSProcessor { return &noopSMS{} }

func (s *noopSMS) Send(ctx context.Context, to, body string) error {
	log.Debug().Str("to", to).Msg("noop sms sent")
	return nil
}
