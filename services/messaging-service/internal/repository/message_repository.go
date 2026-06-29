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

type MessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	ListByConversation(ctx context.Context, convID uuid.UUID, cursor time.Time, limit int) ([]*domain.Message, error)
	UpdateContent(ctx context.Context, id uuid.UUID, content string) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	AddReaction(ctx context.Context, reaction *domain.Reaction) error
	RemoveReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error
	GetReactions(ctx context.Context, msgID uuid.UUID) ([]*domain.Reaction, error)
}

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(ctx context.Context, msg *domain.Message) error {
	query := `INSERT INTO messages (id, conversation_id, sender_id, content, content_type, reply_to, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.ContentType, msg.ReplyTo, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	query := `SELECT id, conversation_id, sender_id, content, content_type, reply_to, deleted_at, created_at
		FROM messages WHERE id = $1 AND deleted_at IS NULL`
	msg := &domain.Message{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.ContentType, &msg.ReplyTo, &msg.DeletedAt, &msg.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) ListByConversation(ctx context.Context, convID uuid.UUID, cursor time.Time, limit int) ([]*domain.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, conversation_id, sender_id, content, content_type, reply_to, deleted_at, created_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL AND created_at < $2
		ORDER BY created_at DESC LIMIT $3`, convID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	var msgs []*domain.Message
	for rows.Next() {
		msg := &domain.Message{}
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.ContentType, &msg.ReplyTo, &msg.DeletedAt, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		msgs = append(msgs, msg)
	}
	if msgs == nil {
		msgs = []*domain.Message{}
	}
	return msgs, nil
}

func (r *messageRepository) UpdateContent(ctx context.Context, id uuid.UUID, content string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE messages SET content = $2 WHERE id = $1 AND deleted_at IS NULL`, id, content)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *messageRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE messages SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`, id, now)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *messageRepository) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	query := `INSERT INTO message_reactions (message_id, user_id, emoji, created_at)
		VALUES ($1, $2, $3, $4) ON CONFLICT (message_id, user_id, emoji) DO NOTHING`
	tag, err := r.pool.Exec(ctx, query, reaction.MessageID, reaction.UserID, reaction.Emoji, reaction.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrDuplicateReaction
	}
	return nil
}

func (r *messageRepository) RemoveReaction(ctx context.Context, msgID, userID uuid.UUID, emoji string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM message_reactions WHERE message_id = $1 AND user_id = $2 AND emoji = $3`, msgID, userID, emoji)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *messageRepository) GetReactions(ctx context.Context, msgID uuid.UUID) ([]*domain.Reaction, error) {
	rows, err := r.pool.Query(ctx, `SELECT message_id, user_id, emoji, created_at FROM message_reactions WHERE message_id = $1 ORDER BY created_at`, msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}
	defer rows.Close()

	var reactions []*domain.Reaction
	for rows.Next() {
		rxn := &domain.Reaction{}
		if err := rows.Scan(&rxn.MessageID, &rxn.UserID, &rxn.Emoji, &rxn.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan reaction: %w", err)
		}
		reactions = append(reactions, rxn)
	}
	if reactions == nil {
		reactions = []*domain.Reaction{}
	}
	return reactions, nil
}
