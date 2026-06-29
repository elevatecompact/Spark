package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/messaging-service/internal/domain"
	"github.com/elevatecompact/spark/services/messaging-service/internal/events"
	"github.com/elevatecompact/spark/services/messaging-service/internal/repository"
)

type MessageService interface {
	Send(ctx context.Context, convID, senderID uuid.UUID, req domain.SendMessageRequest) (*domain.Message, error)
	GetHistory(ctx context.Context, convID uuid.UUID, cursor time.Time, limit int) ([]*domain.Message, error)
	Edit(ctx context.Context, msgID uuid.UUID, content string) error
	Delete(ctx context.Context, msgID uuid.UUID) error
	AddReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error
}

type messageService struct {
	msgRepo  repository.MessageRepository
	convRepo repository.ConversationRepository
	eventPub events.EventProducer
	maxLen   int
}

func NewMessageService(
	msgRepo repository.MessageRepository,
	convRepo repository.ConversationRepository,
	eventPub events.EventProducer,
	maxLen int,
) MessageService {
	return &messageService{
		msgRepo:  msgRepo,
		convRepo: convRepo,
		eventPub: eventPub,
		maxLen:   maxLen,
	}
}

func (s *messageService) Send(ctx context.Context, convID, senderID uuid.UUID, req domain.SendMessageRequest) (*domain.Message, error) {
	member, err := s.convRepo.IsMember(ctx, convID, senderID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, domain.ErrNotMember
	}

	if len(req.Content) > s.maxLen {
		return nil, domain.ErrMsgTooLong
	}
	if req.ContentType == "" {
		req.ContentType = domain.MsgText
	}

	msg := &domain.Message{
		ID:             uuid.New(),
		ConversationID: convID,
		SenderID:       senderID,
		Content:        req.Content,
		ContentType:    req.ContentType,
		ReplyTo:        req.ReplyTo,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	if err := s.eventPub.PublishMessageSent(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to publish message event: %w", err)
	}

	return msg, nil
}

func (s *messageService) GetHistory(ctx context.Context, convID uuid.UUID, cursor time.Time, limit int) ([]*domain.Message, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.msgRepo.ListByConversation(ctx, convID, cursor, limit)
}

func (s *messageService) Edit(ctx context.Context, msgID uuid.UUID, content string) error {
	msg, err := s.msgRepo.GetByID(ctx, msgID)
	if err != nil {
		return err
	}
	if time.Since(msg.CreatedAt) > time.Hour {
		return domain.ErrEditWindowExpired
	}
	if len(content) > s.maxLen {
		return domain.ErrMsgTooLong
	}
	return s.msgRepo.UpdateContent(ctx, msgID, content)
}

func (s *messageService) Delete(ctx context.Context, msgID uuid.UUID) error {
	return s.msgRepo.SoftDelete(ctx, msgID)
}

func (s *messageService) AddReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error {
	rxn := &domain.Reaction{
		MessageID: msgID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now().UTC(),
	}
	return s.msgRepo.AddReaction(ctx, rxn)
}

func (s *messageService) RemoveReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error {
	return s.msgRepo.RemoveReaction(ctx, msgID, userID, emoji)
}
