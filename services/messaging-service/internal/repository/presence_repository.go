package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PresenceRepository interface {
	SetTyping(ctx context.Context, convID, userID uuid.UUID) error
	GetTyping(ctx context.Context, convID uuid.UUID) ([]uuid.UUID, error)
	SetOnline(ctx context.Context, userID uuid.UUID) error
	IsOnline(ctx context.Context, userID uuid.UUID) (bool, error)
}

type presenceRepository struct {
	rdb *redis.Client
}

func NewPresenceRepository(rdb *redis.Client) PresenceRepository {
	return &presenceRepository{rdb: rdb}
}

func (r *presenceRepository) SetTyping(ctx context.Context, convID, userID uuid.UUID) error {
	key := fmt.Sprintf("typing:%s", convID.String())
	err := r.rdb.Set(ctx, key+":"+userID.String(), "1", 10*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set typing indicator: %w", err)
	}
	return nil
}

func (r *presenceRepository) GetTyping(ctx context.Context, convID uuid.UUID) ([]uuid.UUID, error) {
	pattern := fmt.Sprintf("typing:%s:*", convID.String())
	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get typing indicators: %w", err)
	}

	var users []uuid.UUID
	for _, k := range keys {
		uid, err := uuid.Parse(k[len(fmt.Sprintf("typing:%s:", convID.String())):])
		if err != nil {
			continue
		}
		users = append(users, uid)
	}
	return users, nil
}

func (r *presenceRepository) SetOnline(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("online:%s", userID.String())
	return r.rdb.Set(ctx, key, "1", 5*time.Minute).Err()
}

func (r *presenceRepository) IsOnline(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("online:%s", userID.String())
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check online status: %w", err)
	}
	return exists > 0, nil
}
