package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *domain.ChatMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error)
	ListByRoom(ctx context.Context, roomID uuid.UUID, cursor time.Time, limit int) ([]*domain.ChatMessage, error)
	Update(ctx context.Context, msg *domain.ChatMessage) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	DeleteOlderThan(ctx context.Context, days int) error
}

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(ctx context.Context, msg *domain.ChatMessage) error {
	query := `INSERT INTO chat_messages (id, room_id, user_id, username, content, content_type, moderation_status, emote_codes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		msg.ID, msg.RoomID, msg.UserID, msg.Username, msg.Content,
		msg.ContentType, msg.ModerationStatus, msg.Emotes, msg.CreatedAt,
	).Scan(&msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	query := `SELECT id, room_id, user_id, username, content, content_type, moderation_status, emote_codes, edited_at, deleted_at, created_at
		FROM chat_messages WHERE id = $1 AND deleted_at IS NULL`

	msg := &domain.ChatMessage{}
	var editedAt, deletedAt *time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username, &msg.Content,
		&msg.ContentType, &msg.ModerationStatus, &msg.Emotes,
		&editedAt, &deletedAt, &msg.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	msg.EditedAt = editedAt
	msg.DeletedAt = deletedAt
	return msg, nil
}

func (r *messageRepository) ListByRoom(ctx context.Context, roomID uuid.UUID, cursor time.Time, limit int) ([]*domain.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var rows pgx.Rows
	var err error

	if cursor.IsZero() {
		rows, err = r.pool.Query(ctx, `
			SELECT id, room_id, user_id, username, content, content_type, moderation_status, emote_codes, edited_at, deleted_at, created_at
			FROM chat_messages
			WHERE room_id = $1 AND deleted_at IS NULL
			ORDER BY created_at DESC LIMIT $2`, roomID, limit)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, room_id, user_id, username, content, content_type, moderation_status, emote_codes, edited_at, deleted_at, created_at
			FROM chat_messages
			WHERE room_id = $1 AND deleted_at IS NULL AND created_at < $2
			ORDER BY created_at DESC LIMIT $3`, roomID, cursor, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.ChatMessage
	for rows.Next() {
		msg := &domain.ChatMessage{}
		var editedAt, deletedAt *time.Time
		err := rows.Scan(
			&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username, &msg.Content,
			&msg.ContentType, &msg.ModerationStatus, &msg.Emotes,
			&editedAt, &deletedAt, &msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		msg.EditedAt = editedAt
		msg.DeletedAt = deletedAt
		messages = append(messages, msg)
	}
	if messages == nil {
		messages = []*domain.ChatMessage{}
	}
	return messages, nil
}

func (r *messageRepository) Update(ctx context.Context, msg *domain.ChatMessage) error {
	now := time.Now().UTC()
	query := `UPDATE chat_messages SET content = $2, content_type = $3, moderation_status = $4, emotes = $5, edited_at = $6 WHERE id = $1 AND deleted_at IS NULL`
	tag, err := r.pool.Exec(ctx, query, msg.ID, msg.Content, msg.ContentType, msg.ModerationStatus, msg.Emotes, now)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	msg.EditedAt = &now
	return nil
}

func (r *messageRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE chat_messages SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`, id, now)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *messageRepository) DeleteOlderThan(ctx context.Context, days int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM chat_messages WHERE created_at < NOW() - INTERVAL '1 day' * $1`, days)
	if err != nil {
		return fmt.Errorf("failed to delete old messages: %w", err)
	}
	return nil
}
