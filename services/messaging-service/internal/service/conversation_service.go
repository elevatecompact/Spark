package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/messaging-service/internal/domain"
	"github.com/elevatecompact/spark/services/messaging-service/internal/events"
	"github.com/elevatecompact/spark/services/messaging-service/internal/repository"
)

type ConversationService interface {
	Create(ctx context.Context, userID uuid.UUID, req domain.CreateConversationRequest) (*domain.Conversation, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Conversation, error)
	List(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Conversation, error)
	Update(ctx context.Context, id uuid.UUID, name *string, iconURL *string) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddMember(ctx context.Context, convID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, convID, userID uuid.UUID) error
	GetMembers(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error)
	MarkRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error
	GetReadStatus(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error)
}

type conversationService struct {
	convRepo repository.ConversationRepository
	eventPub events.EventProducer
}

func NewConversationService(convRepo repository.ConversationRepository, eventPub events.EventProducer) ConversationService {
	return &conversationService{convRepo: convRepo, eventPub: eventPub}
}

func (s *conversationService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateConversationRequest) (*domain.Conversation, error) {
	now := time.Now().UTC()
	conv := &domain.Conversation{
		ID:        uuid.New(),
		Type:      req.Type,
		Name:      req.Name,
		CreatedBy: userID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.convRepo.Create(ctx, conv); err != nil {
		return nil, err
	}

	member := &domain.ConversationMember{
		ConversationID: conv.ID,
		UserID:         userID,
		Role:           domain.RoleAdmin,
		JoinedAt:       now,
	}
	s.convRepo.AddMember(ctx, member)

	for _, mid := range req.MemberIDs {
		m := &domain.ConversationMember{
			ConversationID: conv.ID,
			UserID:         mid,
			Role:           domain.RoleMember,
			JoinedAt:       now,
		}
		s.convRepo.AddMember(ctx, m)
	}

	if err := s.eventPub.PublishConversationCreated(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *conversationService) Get(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	return s.convRepo.GetByID(ctx, id)
}

func (s *conversationService) List(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Conversation, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.convRepo.ListByUser(ctx, userID, cursor, limit)
}

func (s *conversationService) Update(ctx context.Context, id uuid.UUID, name *string, iconURL *string) error {
	conv, err := s.convRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	conv.Name = name
	conv.IconURL = iconURL
	return s.convRepo.Update(ctx, conv)
}

func (s *conversationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.convRepo.SoftDelete(ctx, id)
}

func (s *conversationService) AddMember(ctx context.Context, convID, userID uuid.UUID) error {
	member := &domain.ConversationMember{
		ConversationID: convID,
		UserID:         userID,
		Role:           domain.RoleMember,
		JoinedAt:       time.Now().UTC(),
	}
	return s.convRepo.AddMember(ctx, member)
}

func (s *conversationService) RemoveMember(ctx context.Context, convID, userID uuid.UUID) error {
	return s.convRepo.RemoveMember(ctx, convID, userID)
}

func (s *conversationService) GetMembers(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error) {
	return s.convRepo.GetMembers(ctx, convID)
}

func (s *conversationService) MarkRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error {
	return s.convRepo.UpdateLastRead(ctx, convID, userID, msgID)
}

func (s *conversationService) GetReadStatus(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error) {
	return s.convRepo.GetMembers(ctx, convID)
}
