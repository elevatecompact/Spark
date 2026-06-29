package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
)

type EmoteRepository interface {
	Create(ctx context.Context, emote *domain.Emote) error
	GetGlobal(ctx context.Context) ([]*domain.Emote, error)
	GetByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.Emote, error)
}

type emoteRepository struct {
	pool *pgxpool.Pool
}

func NewEmoteRepository(pool *pgxpool.Pool) EmoteRepository {
	return &emoteRepository{pool: pool}
}

func (r *emoteRepository) Create(ctx context.Context, emote *domain.Emote) error {
	query := `INSERT INTO chat_emotes (id, code, image_url, is_global, room_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`
	err := r.pool.QueryRow(ctx, query,
		emote.ID, emote.Code, emote.ImageURL, emote.IsGlobal, emote.RoomID, emote.CreatedAt,
	).Scan(&emote.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create emote: %w", err)
	}
	return nil
}

func (r *emoteRepository) GetGlobal(ctx context.Context) ([]*domain.Emote, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, code, image_url, is_global, room_id, created_at FROM chat_emotes WHERE is_global = true ORDER BY code ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list global emotes: %w", err)
	}
	defer rows.Close()

	var emotes []*domain.Emote
	for rows.Next() {
		e := &domain.Emote{}
		if err := rows.Scan(&e.ID, &e.Code, &e.ImageURL, &e.IsGlobal, &e.RoomID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan emote: %w", err)
		}
		emotes = append(emotes, e)
	}
	if emotes == nil {
		emotes = []*domain.Emote{}
	}
	return emotes, nil
}

func (r *emoteRepository) GetByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.Emote, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, code, image_url, is_global, room_id, created_at FROM chat_emotes WHERE room_id = $1 OR is_global = true ORDER BY is_global DESC, code ASC`, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list room emotes: %w", err)
	}
	defer rows.Close()

	var emotes []*domain.Emote
	for rows.Next() {
		e := &domain.Emote{}
		if err := rows.Scan(&e.ID, &e.Code, &e.ImageURL, &e.IsGlobal, &e.RoomID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan emote: %w", err)
		}
		emotes = append(emotes, e)
	}
	if emotes == nil {
		emotes = []*domain.Emote{}
	}
	return emotes, nil
}
