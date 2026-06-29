package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
	"github.com/elevatecompact/spark/services/chat-service/internal/events"
	"github.com/elevatecompact/spark/services/chat-service/internal/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, roomID, userID uuid.UUID, username string, req domain.SendMessageRequest) (*domain.ChatMessage, error)
	GetHistory(ctx context.Context, roomID uuid.UUID, cursor time.Time, limit int) ([]*domain.ChatMessage, error)
	EditMessage(ctx context.Context, msgID uuid.UUID, content string) (*domain.ChatMessage, error)
	DeleteMessage(ctx context.Context, msgID uuid.UUID) error
}

type messageService struct {
	msgRepo      repository.MessageRepository
	roomRepo     repository.RoomRepository
	modRepo      repository.ModerationRepository
	eventPub     events.EventProducer
	maxMsgLen    int
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	roomRepo repository.RoomRepository,
	modRepo repository.ModerationRepository,
	eventPub events.EventProducer,
	maxMsgLen int,
) MessageService {
	return &messageService{
		msgRepo:   msgRepo,
		roomRepo:  roomRepo,
		modRepo:   modRepo,
		eventPub:  eventPub,
		maxMsgLen: maxMsgLen,
	}
}

func (s *messageService) SendMessage(ctx context.Context, roomID, userID uuid.UUID, username string, req domain.SendMessageRequest) (*domain.ChatMessage, error) {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !room.IsActive {
		return nil, domain.NewDomainErrorMsg(domain.ErrForbidden, "room is closed", 403)
	}

	banned, err := s.modRepo.IsUserBanned(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ban status: %w", err)
	}
	if banned {
		return nil, domain.ErrUserBanned
	}

	muted, err := s.modRepo.IsUserMuted(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check mute status: %w", err)
	}
	if muted {
		return nil, domain.ErrUserMuted
	}

	if len(req.Content) > s.maxMsgLen {
		return nil, domain.ErrMessageTooLong
	}

	if req.ContentType == "" {
		req.ContentType = domain.ContentText
	}

	now := time.Now().UTC()
	msg := &domain.ChatMessage{
		ID:               uuid.New(),
		RoomID:           roomID,
		UserID:           userID,
		Username:         username,
		Content:          req.Content,
		ContentType:      req.ContentType,
		ModerationStatus: domain.ModApproved,
		Emotes:           []string{},
		CreatedAt:        now,
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	if err := s.eventPub.PublishMessageSent(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to publish message event: %w", err)
	}

	return msg, nil
}

func (s *messageService) GetHistory(ctx context.Context, roomID uuid.UUID, cursor time.Time, limit int) ([]*domain.ChatMessage, error) {
	return s.msgRepo.ListByRoom(ctx, roomID, cursor, limit)
}

func (s *messageService) EditMessage(ctx context.Context, msgID uuid.UUID, content string) (*domain.ChatMessage, error) {
	msg, err := s.msgRepo.GetByID(ctx, msgID)
	if err != nil {
		return nil, err
	}

	if len(content) > s.maxMsgLen {
		return nil, domain.ErrMessageTooLong
	}

	msg.Content = content
	if err := s.msgRepo.Update(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	return msg, nil
}

func (s *messageService) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	msg, err := s.msgRepo.GetByID(ctx, msgID)
	if err != nil {
		return err
	}

	if err := s.msgRepo.SoftDelete(ctx, msgID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return s.eventPub.PublishMessageDeleted(ctx, msg)
}
