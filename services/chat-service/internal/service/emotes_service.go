package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
	"github.com/elevatecompact/spark/services/chat-service/internal/repository"
)

type EmoteService interface {
	GetGlobal(ctx context.Context) ([]*domain.Emote, error)
	GetByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.Emote, error)
}

type emoteService struct {
	repo repository.EmoteRepository
}

func NewEmoteService(repo repository.EmoteRepository) EmoteService {
	return &emoteService{repo: repo}
}

func (s *emoteService) GetGlobal(ctx context.Context) ([]*domain.Emote, error) {
	return s.repo.GetGlobal(ctx)
}

func (s *emoteService) GetByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.Emote, error) {
	return s.repo.GetByRoom(ctx, roomID)
}
