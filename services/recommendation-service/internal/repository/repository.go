package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/recommendation-service/internal/domain"
)

type EmbeddingRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserEmbedding, error)
	UpsertUser(ctx context.Context, e *domain.UserEmbedding) error
	GetContent(ctx context.Context, contentID uuid.UUID) (*domain.ContentEmbedding, error)
	UpsertContent(ctx context.Context, e *domain.ContentEmbedding) error
	ListActiveModels(ctx context.Context) ([]domain.ModelInfo, error)
}

type InteractionRepository interface {
	Insert(ctx context.Context, i *domain.UserContentInteraction) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.UserContentInteraction, error)
	CountByContent(ctx context.Context, contentID uuid.UUID, sinceMinutes int) (int64, error)
	CountAllSince(ctx context.Context, sinceMinutes int) (int64, error)
	TopContentSince(ctx context.Context, sinceMinutes, limit int) ([]uuid.UUID, error)
	ListUsersInteractedWith(ctx context.Context, userID uuid.UUID, limit int) ([]uuid.UUID, error)
	ListByContent(ctx context.Context, contentID uuid.UUID, limit int) ([]*domain.UserContentInteraction, error)
}

type ContentMeta struct {
	ID        uuid.UUID `json:"id"`
	CreatorID uuid.UUID `json:"creator_id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	Category  string    `json:"category"`
	CreatedAt int64     `json:"created_at"`
}

type ContentRepository interface {
	GetMeta(ctx context.Context, contentID uuid.UUID) (*ContentMeta, error)
	ListByCreator(ctx context.Context, creatorID uuid.UUID, limit int) ([]uuid.UUID, error)
	RandomSample(ctx context.Context, limit int) ([]uuid.UUID, error)
}

type embeddingRepo struct{ pool *pgxpool.Pool }
type interactionRepo struct{ pool *pgxpool.Pool }
type contentRepo struct{ pool *pgxpool.Pool }

func NewEmbeddingRepository(pool *pgxpool.Pool) EmbeddingRepository    { return &embeddingRepo{pool} }
func NewInteractionRepository(pool *pgxpool.Pool) InteractionRepository { return &interactionRepo{pool} }
func NewContentRepository(pool *pgxpool.Pool) ContentRepository         { return &contentRepo{pool} }

func (r *embeddingRepo) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserEmbedding, error) {
	e := &domain.UserEmbedding{}
	var emb []byte
	err := r.pool.QueryRow(ctx, `SELECT user_id, embedding, model_version, updated_at FROM user_embeddings WHERE user_id=$1`, userID).Scan(&e.UserID, &emb, &e.ModelVersion, &e.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNoEmbedding
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(emb, &e.Embedding); err != nil {
		return nil, err
	}
	return e, nil
}

func (r *embeddingRepo) UpsertUser(ctx context.Context, e *domain.UserEmbedding) error {
	emb, err := json.Marshal(e.Embedding)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO user_embeddings (user_id, embedding, model_version, updated_at) VALUES ($1,$2,$3,$4) ON CONFLICT (user_id) DO UPDATE SET embedding=$2, model_version=$3, updated_at=NOW()`, e.UserID, emb, e.ModelVersion)
	return err
}

func (r *embeddingRepo) GetContent(ctx context.Context, contentID uuid.UUID) (*domain.ContentEmbedding, error) {
	e := &domain.ContentEmbedding{}
	var emb []byte
	err := r.pool.QueryRow(ctx, `SELECT content_id, embedding, model_version, updated_at FROM content_embeddings WHERE content_id=$1`, contentID).Scan(&e.ContentID, &emb, &e.ModelVersion, &e.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNoEmbedding
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(emb, &e.Embedding); err != nil {
		return nil, err
	}
	return e, nil
}

func (r *embeddingRepo) UpsertContent(ctx context.Context, e *domain.ContentEmbedding) error {
	emb, err := json.Marshal(e.Embedding)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO content_embeddings (content_id, embedding, model_version, updated_at) VALUES ($1,$2,$3,$4) ON CONFLICT (content_id) DO UPDATE SET embedding=$2, model_version=$3, updated_at=NOW()`, e.ContentID, emb, e.ModelVersion)
	return err
}

func (r *embeddingRepo) ListActiveModels(ctx context.Context) ([]domain.ModelInfo, error) {
	rows, err := r.pool.Query(ctx, `SELECT version, deployed_at, metrics, is_active FROM model_versions WHERE is_active=true ORDER BY deployed_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []domain.ModelInfo
	for rows.Next() {
		m := domain.ModelInfo{}
		if err := rows.Scan(&m.Version, &m.DeployedAt, &m.Metrics, &m.IsActive); err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	if models == nil {
		models = []domain.ModelInfo{}
	}
	return models, nil
}

func (r *interactionRepo) Insert(ctx context.Context, i *domain.UserContentInteraction) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_content_interactions (user_id, content_id, interaction_type, weight, timestamp) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (user_id, content_id) DO UPDATE SET weight=$4, timestamp=$5`,
		i.UserID, i.ContentID, i.InteractionType, i.Weight, i.Timestamp)
	return err
}

func (r *interactionRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.UserContentInteraction, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT user_id, content_id, interaction_type, weight, timestamp FROM user_content_interactions WHERE user_id=$1 ORDER BY timestamp DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var is []*domain.UserContentInteraction
	for rows.Next() {
		i := &domain.UserContentInteraction{}
		if err := rows.Scan(&i.UserID, &i.ContentID, &i.InteractionType, &i.Weight, &i.Timestamp); err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	if is == nil {
		is = []*domain.UserContentInteraction{}
	}
	return is, nil
}

// CountByContent returns the number of weighted interactions for a content item
// within the last `sinceMinutes` minutes. Used to power the trending feed.
func (r *interactionRepo) CountByContent(ctx context.Context, contentID uuid.UUID, sinceMinutes int) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(weight), 0) FROM user_content_interactions WHERE content_id=$1 AND timestamp >= NOW() - ($2::int * INTERVAL '1 minute')`, contentID, sinceMinutes).Scan(&count)
	return count, err
}

// CountAllSince returns the total weight of all interactions in the window.
func (r *interactionRepo) CountAllSince(ctx context.Context, sinceMinutes int) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(weight), 0) FROM user_content_interactions WHERE timestamp >= NOW() - ($1::int * INTERVAL '1 minute')`, sinceMinutes).Scan(&count)
	return count, err
}

// TopContentSince returns the highest-weighted content IDs in the window.
func (r *interactionRepo) TopContentSince(ctx context.Context, sinceMinutes, limit int) ([]uuid.UUID, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT content_id, COALESCE(SUM(weight), 0) AS score FROM user_content_interactions WHERE timestamp >= NOW() - ($1::int * INTERVAL '1 minute') GROUP BY content_id ORDER BY score DESC LIMIT $2`, sinceMinutes, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		var score float64
		if err := rows.Scan(&id, &score); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ListUsersInteractedWith returns users who interacted with the same content
// as the supplied user. Used to seed collaborative filtering.
func (r *interactionRepo) ListUsersInteractedWith(ctx context.Context, userID uuid.UUID, limit int) ([]uuid.UUID, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT uci2.user_id
		FROM user_content_interactions uci1
		JOIN user_content_interactions uci2 ON uci1.content_id = uci2.content_id
		WHERE uci1.user_id = $1 AND uci2.user_id <> $1
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ListByContent returns recent interactions for a content item.
func (r *interactionRepo) ListByContent(ctx context.Context, contentID uuid.UUID, limit int) ([]*domain.UserContentInteraction, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT user_id, content_id, interaction_type, weight, timestamp FROM user_content_interactions WHERE content_id=$1 ORDER BY timestamp DESC LIMIT $2`, contentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var is []*domain.UserContentInteraction
	for rows.Next() {
		i := &domain.UserContentInteraction{}
		if err := rows.Scan(&i.UserID, &i.ContentID, &i.InteractionType, &i.Weight, &i.Timestamp); err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	if is == nil {
		is = []*domain.UserContentInteraction{}
	}
	return is, nil
}

func (r *contentRepo) GetMeta(ctx context.Context, contentID uuid.UUID) (*ContentMeta, error) {
	m := &ContentMeta{}
	err := r.pool.QueryRow(ctx, `SELECT id, creator_id, COALESCE(title, ''), COALESCE(tags, '{}'::text[]), COALESCE(category, ''), COALESCE(EXTRACT(EPOCH FROM created_at)::bigint, 0) FROM content_items WHERE id=$1`, contentID).
		Scan(&m.ID, &m.CreatorID, &m.Title, &m.Tags, &m.Category, &m.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *contentRepo) ListByCreator(ctx context.Context, creatorID uuid.UUID, limit int) ([]uuid.UUID, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id FROM content_items WHERE creator_id=$1 ORDER BY created_at DESC LIMIT $2`, creatorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *contentRepo) RandomSample(ctx context.Context, limit int) ([]uuid.UUID, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id FROM content_items ORDER BY RANDOM() LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
