package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.ChatRoom) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error)
	Update(ctx context.Context, room *domain.ChatRoom) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type roomRepository struct {
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) RoomRepository {
	return &roomRepository{pool: pool}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.ChatRoom) error {
	query := `INSERT INTO chat_rooms (id, name, type, owner_id, slow_mode_seconds, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		room.ID, room.Name, room.Type, room.OwnerID,
		room.SlowModeSeconds, room.IsActive, room.CreatedAt, room.UpdatedAt,
	).Scan(&room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error) {
	query := `SELECT id, name, type, owner_id, slow_mode_seconds, is_active, created_at, updated_at
		FROM chat_rooms WHERE id = $1`

	room := &domain.ChatRoom{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&room.ID, &room.Name, &room.Type, &room.OwnerID,
		&room.SlowModeSeconds, &room.IsActive, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRoomNotFound
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}
	return room, nil
}

func (r *roomRepository) Update(ctx context.Context, room *domain.ChatRoom) error {
	query := `UPDATE chat_rooms SET name = $2, slow_mode_seconds = $3, is_active = $4, updated_at = NOW() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, room.ID, room.Name, room.SlowModeSeconds, room.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRoomNotFound
	}
	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM chat_rooms WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRoomNotFound
	}
	return nil
}
