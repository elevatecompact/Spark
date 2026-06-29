package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/events"
	"github.com/elevatecompact/spark/services/chat-service/internal/repository"
)

type ModerationService interface {
	MuteUser(ctx context.Context, roomID, userID uuid.UUID, duration time.Duration) error
	UnmuteUser(ctx context.Context, roomID, userID uuid.UUID) error
	BanUser(ctx context.Context, roomID, userID uuid.UUID, reason string, duration time.Duration) error
	UnbanUser(ctx context.Context, roomID, userID uuid.UUID) error
	SetSlowMode(ctx context.Context, roomID uuid.UUID, intervalSecs int) error
}

type moderationService struct {
	repo     repository.ModerationRepository
	eventPub events.EventProducer
}

func NewModerationService(repo repository.ModerationRepository, eventPub events.EventProducer) ModerationService {
	return &moderationService{repo: repo, eventPub: eventPub}
}

func (s *moderationService) MuteUser(ctx context.Context, roomID, userID uuid.UUID, duration time.Duration) error {
	if err := s.repo.MuteUser(ctx, roomID, userID, duration); err != nil {
		return err
	}
	return s.eventPub.PublishUserMuted(ctx, roomID, userID)
}

func (s *moderationService) UnmuteUser(ctx context.Context, roomID, userID uuid.UUID) error {
	return s.repo.UnmuteUser(ctx, roomID, userID)
}

func (s *moderationService) BanUser(ctx context.Context, roomID, userID uuid.UUID, reason string, duration time.Duration) error {
	if err := s.repo.BanUser(ctx, roomID, userID, reason, duration); err != nil {
		return err
	}
	return s.eventPub.PublishUserBanned(ctx, roomID, userID)
}

func (s *moderationService) UnbanUser(ctx context.Context, roomID, userID uuid.UUID) error {
	return s.repo.UnbanUser(ctx, roomID, userID)
}

func (s *moderationService) SetSlowMode(ctx context.Context, roomID uuid.UUID, intervalSecs int) error {
	return s.repo.SetSlowMode(ctx, roomID, intervalSecs)
}
