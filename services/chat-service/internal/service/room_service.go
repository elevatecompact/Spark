package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
	"github.com/elevatecompact/spark/services/chat-service/internal/events"
	"github.com/elevatecompact/spark/services/chat-service/internal/repository"
)

type RoomService interface {
	Create(ctx context.Context, req domain.CreateRoomRequest) (*domain.ChatRoom, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error)
	Close(ctx context.Context, id uuid.UUID) error
}

type roomService struct {
	repo     repository.RoomRepository
	eventPub events.EventProducer
}

func NewRoomService(repo repository.RoomRepository, eventPub events.EventProducer) RoomService {
	return &roomService{repo: repo, eventPub: eventPub}
}

func (s *roomService) Create(ctx context.Context, req domain.CreateRoomRequest) (*domain.ChatRoom, error) {
	now := time.Now().UTC()
	room := &domain.ChatRoom{
		ID:              uuid.New(),
		Name:            req.Name,
		Type:            req.Type,
		OwnerID:         req.OwnerID,
		SlowModeSeconds: 0,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, room); err != nil {
		return nil, err
	}

	if err := s.eventPub.PublishRoomCreated(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) Get(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *roomService) Close(ctx context.Context, id uuid.UUID) error {
	room, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	room.IsActive = false
	room.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, room); err != nil {
		return err
	}

	return s.eventPub.PublishRoomClosed(ctx, room)
}
