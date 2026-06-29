package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ModerationRepository interface {
	MuteUser(ctx context.Context, roomID, userID uuid.UUID, duration time.Duration) error
	UnmuteUser(ctx context.Context, roomID, userID uuid.UUID) error
	IsUserMuted(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
	BanUser(ctx context.Context, roomID, userID uuid.UUID, reason string, duration time.Duration) error
	UnbanUser(ctx context.Context, roomID, userID uuid.UUID) error
	IsUserBanned(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
	SetSlowMode(ctx context.Context, roomID uuid.UUID, intervalSecs int) error
	GetSlowMode(ctx context.Context, roomID uuid.UUID) (int, error)
}

type moderationRepository struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func NewModerationRepository(pool *pgxpool.Pool, rdb *redis.Client) ModerationRepository {
	return &moderationRepository{pool: pool, rdb: rdb}
}

func (r *moderationRepository) MuteUser(ctx context.Context, roomID, userID uuid.UUID, duration time.Duration) error {
	key := fmt.Sprintf("mute:%s:%s", roomID.String(), userID.String())
	return r.rdb.Set(ctx, key, "1", duration).Err()
}

func (r *moderationRepository) UnmuteUser(ctx context.Context, roomID, userID uuid.UUID) error {
	key := fmt.Sprintf("mute:%s:%s", roomID.String(), userID.String())
	return r.rdb.Del(ctx, key).Err()
}

func (r *moderationRepository) IsUserMuted(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("mute:%s:%s", roomID.String(), userID.String())
	exists, err := r.rdb.Exists(ctx, key).Result()
	return exists > 0, err
}

func (r *moderationRepository) BanUser(ctx context.Context, roomID, userID uuid.UUID, reason string, duration time.Duration) error {
	key := fmt.Sprintf("ban:%s:%s", roomID.String(), userID.String())
	return r.rdb.Set(ctx, key, reason, duration).Err()
}

func (r *moderationRepository) UnbanUser(ctx context.Context, roomID, userID uuid.UUID) error {
	key := fmt.Sprintf("ban:%s:%s", roomID.String(), userID.String())
	return r.rdb.Del(ctx, key).Err()
}

func (r *moderationRepository) IsUserBanned(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("ban:%s:%s", roomID.String(), userID.String())
	exists, err := r.rdb.Exists(ctx, key).Result()
	return exists > 0, err
}

func (r *moderationRepository) SetSlowMode(ctx context.Context, roomID uuid.UUID, intervalSecs int) error {
	key := fmt.Sprintf("slowmode:%s", roomID.String())
	if intervalSecs <= 0 {
		return r.rdb.Del(ctx, key).Err()
	}
	return r.rdb.Set(ctx, key, intervalSecs, 0).Err()
}

func (r *moderationRepository) GetSlowMode(ctx context.Context, roomID uuid.UUID) (int, error) {
	key := fmt.Sprintf("slowmode:%s", roomID.String())
	val, err := r.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
