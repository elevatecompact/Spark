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
}

type embeddingRepo struct{ pool *pgxpool.Pool }
type interactionRepo struct{ pool *pgxpool.Pool }

func NewEmbeddingRepository(pool *pgxpool.Pool) EmbeddingRepository   { return &embeddingRepo{pool} }
func NewInteractionRepository(pool *pgxpool.Pool) InteractionRepository { return &interactionRepo{pool} }

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
	if limit <= 0 || limit > 100 {
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
