package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, tokenHash string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	Cleanup(ctx context.Context) error
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}

type sessionRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionRepository(client *redis.Client, ttl time.Duration) SessionRepository {
	return &sessionRepository{
		client: client,
		ttl:    ttl,
	}
}

func sessionKey(id uuid.UUID) string {
	return fmt.Sprintf("session:%s", id.String())
}

func userSessionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_sessions:%s", userID.String())
}

func tokenSessionKey(tokenHash string) string {
	return fmt.Sprintf("token_session:%s", tokenHash)
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	pipe := r.client.Pipeline()

	pipe.Set(ctx, sessionKey(session.ID), data, r.ttl)

	pipe.Set(ctx, tokenSessionKey(session.TokenHash), session.ID.String(), r.ttl)

	pipe.SAdd(ctx, userSessionsKey(session.UserID), session.ID.String())
	pipe.Expire(ctx, userSessionsKey(session.UserID), r.ttl)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session in redis: %w", err)
	}
	return nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, tokenHash string) (*domain.Session, error) {
	sessionID, err := r.client.Get(ctx, tokenSessionKey(tokenHash)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrSessionExpired
		}
		return nil, fmt.Errorf("failed to get session id by token hash: %w", err)
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session id: %w", err)
	}

	data, err := r.client.Get(ctx, sessionKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrSessionExpired
		}
		return nil, fmt.Errorf("failed to get session data: %w", err)
	}

	session := &domain.Session{}
	if err := json.Unmarshal(data, session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		r.Delete(ctx, session.ID)
		return nil, domain.ErrSessionExpired
	}

	return session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	sessionIDs, err := r.client.SMembers(ctx, userSessionsKey(userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session ids for user: %w", err)
	}

	var sessions []*domain.Session
	for _, sid := range sessionIDs {
		id, err := uuid.Parse(sid)
		if err != nil {
			continue
		}

		data, err := r.client.Get(ctx, sessionKey(id)).Bytes()
		if err != nil {
			continue
		}

		session := &domain.Session{}
		if err := json.Unmarshal(data, session); err != nil {
			continue
		}

		if time.Now().After(session.ExpiresAt) {
			r.Delete(ctx, session.ID)
			continue
		}

		sessions = append(sessions, session)
	}

	if sessions == nil {
		sessions = []*domain.Session{}
	}
	return sessions, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pipe := r.client.Pipeline()

	data, err := r.client.Get(ctx, sessionKey(id)).Bytes()
	if err == nil {
		session := &domain.Session{}
		if json.Unmarshal(data, session) == nil {
			pipe.SRem(ctx, userSessionsKey(session.UserID), id.String())
			pipe.Del(ctx, tokenSessionKey(session.TokenHash))
		}
	}

	pipe.Del(ctx, sessionKey(id))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	sessionID, err := r.client.Get(ctx, tokenSessionKey(tokenHash)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return fmt.Errorf("failed to get session id by token hash: %w", err)
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return fmt.Errorf("failed to parse session id: %w", err)
	}

	return r.Delete(ctx, id)
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	sessionIDs, err := r.client.SMembers(ctx, userSessionsKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get session ids for user: %w", err)
	}

	pipe := r.client.Pipeline()
	for _, sid := range sessionIDs {
		id, err := uuid.Parse(sid)
		if err != nil {
			continue
		}

		data, errData := r.client.Get(ctx, sessionKey(id)).Bytes()
		if errData == nil {
			session := &domain.Session{}
			if json.Unmarshal(data, session) == nil {
				pipe.Del(ctx, tokenSessionKey(session.TokenHash))
			}
		}

		pipe.Del(ctx, sessionKey(id))
	}
	pipe.Del(ctx, userSessionsKey(userID))

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete sessions for user: %w", err)
	}
	return nil
}

func (r *sessionRepository) Cleanup(ctx context.Context) error {
	var cursor uint64
	keysProcessed := 0

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, "session:*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan sessions: %w", err)
		}

		for _, key := range keys {
			data, err := r.client.Get(ctx, key).Bytes()
			if err != nil {
				continue
			}

			session := &domain.Session{}
			if json.Unmarshal(data, session) != nil {
				continue
			}

			if time.Now().After(session.ExpiresAt) {
				r.Delete(ctx, session.ID)
				keysProcessed++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
