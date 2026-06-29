package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/messaging-service/internal/domain"
)

type ConversationRepository interface {
	Create(ctx context.Context, conv *domain.Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Conversation, error)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Conversation, error)
	Update(ctx context.Context, conv *domain.Conversation) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	AddMember(ctx context.Context, member *domain.ConversationMember) error
	RemoveMember(ctx context.Context, convID, userID uuid.UUID) error
	GetMembers(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error)
	IsMember(ctx context.Context, convID, userID uuid.UUID) (bool, error)
	UpdateLastRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error
}

type conversationRepository struct {
	pool *pgxpool.Pool
}

func NewConversationRepository(pool *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{pool: pool}
}

func (r *conversationRepository) Create(ctx context.Context, conv *domain.Conversation) error {
	query := `INSERT INTO conversations (id, type, name, icon_url, created_by, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, conv.ID, conv.Type, conv.Name, conv.IconURL, conv.CreatedBy, conv.IsActive, conv.CreatedAt, conv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	return nil
}

func (r *conversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	query := `SELECT id, type, name, icon_url, created_by, is_active, created_at, updated_at
		FROM conversations WHERE id = $1`
	conv := &domain.Conversation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.IconURL, &conv.CreatedBy, &conv.IsActive, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrConvNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return conv, nil
}

func (r *conversationRepository) ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.type, c.name, c.icon_url, c.created_by, c.is_active, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_members cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1 AND c.is_active = true AND c.created_at < $2
		ORDER BY c.created_at DESC LIMIT $3`, userID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var convs []*domain.Conversation
	for rows.Next() {
		conv := &domain.Conversation{}
		if err := rows.Scan(&conv.ID, &conv.Type, &conv.Name, &conv.IconURL, &conv.CreatedBy, &conv.IsActive, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		convs = append(convs, conv)
	}
	if convs == nil {
		convs = []*domain.Conversation{}
	}
	return convs, nil
}

func (r *conversationRepository) Update(ctx context.Context, conv *domain.Conversation) error {
	query := `UPDATE conversations SET name = $2, icon_url = $3, updated_at = NOW() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, conv.ID, conv.Name, conv.IconURL)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrConvNotFound
	}
	return nil
}

func (r *conversationRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `UPDATE conversations SET is_active = false, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrConvNotFound
	}
	return nil
}

func (r *conversationRepository) AddMember(ctx context.Context, member *domain.ConversationMember) error {
	query := `INSERT INTO conversation_members (conversation_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, member.ConversationID, member.UserID, member.Role, member.JoinedAt)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}
	return nil
}

func (r *conversationRepository) RemoveMember(ctx context.Context, convID, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM conversation_members WHERE conversation_id = $1 AND user_id = $2`, convID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotMember
	}
	return nil
}

func (r *conversationRepository) GetMembers(ctx context.Context, convID uuid.UUID) ([]*domain.ConversationMember, error) {
	rows, err := r.pool.Query(ctx, `SELECT conversation_id, user_id, role, last_read_message_id, joined_at FROM conversation_members WHERE conversation_id = $1`, convID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	var members []*domain.ConversationMember
	for rows.Next() {
		m := &domain.ConversationMember{}
		if err := rows.Scan(&m.ConversationID, &m.UserID, &m.Role, &m.LastReadMsgID, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, m)
	}
	if members == nil {
		members = []*domain.ConversationMember{}
	}
	return members, nil
}

func (r *conversationRepository) UpdateLastRead(ctx context.Context, convID, userID uuid.UUID, msgID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE conversation_members SET last_read_message_id = $3 WHERE conversation_id = $1 AND user_id = $2`, convID, userID, msgID)
	if err != nil {
		return fmt.Errorf("failed to update last read: %w", err)
	}
	return nil
}

func (r *conversationRepository) IsMember(ctx context.Context, convID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id = $1 AND user_id = $2)`, convID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	return exists, nil
}
