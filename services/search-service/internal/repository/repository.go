package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/search-service/internal/domain"
)

type SearchEngine interface {
	Search(ctx context.Context, q *domain.SearchQuery) (*domain.SearchResult, error)
	Autocomplete(ctx context.Context, prefix string, size int) ([]domain.AutocompleteSuggestion, error)
	Index(ctx context.Context, contentType domain.ContentType, doc *domain.SearchDocument) error
	Update(ctx context.Context, contentType domain.ContentType, id uuid.UUID, doc map[string]interface{}) error
	Remove(ctx context.Context, contentType domain.ContentType, id uuid.UUID) error
	Reindex(ctx context.Context, contentType domain.ContentType, docs []domain.SearchDocument) error
	GetStats(ctx context.Context) ([]domain.IndexStats, error)
	Health(ctx context.Context) (*domain.ESHealth, error)
	PutSynonyms(ctx context.Context, set *domain.SynonymSet) error
	PutAnalyzers(ctx context.Context, config map[string]interface{}) error
}

type SearchRepository interface {
	LogSearch(ctx context.Context, a *domain.SearchAnalytics) error
}

type searchRepo struct{ pool *pgxpool.Pool }

func NewSearchRepository(pool *pgxpool.Pool) SearchRepository { return &searchRepo{pool} }

func (r *searchRepo) LogSearch(ctx context.Context, a *domain.SearchAnalytics) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO search_analytics (query, result_ids, latency_ms, user_id, timestamp) VALUES ($1,$2,$3,$4,$5)`,
		a.Query, a.ResultIDs, a.LatencyMs, a.UserID, a.Timestamp)
	return err
}

type noopSearchEngine struct{}

func NewNoopSearchEngine() SearchEngine { return &noopSearchEngine{} }

func (e *noopSearchEngine) Search(ctx context.Context, q *domain.SearchQuery) (*domain.SearchResult, error) {
	return &domain.SearchResult{
		Total:   0,
		Page:    q.Page,
		Size:    q.Size,
		Results: []domain.SearchDocument{},
	}, nil
}
func (e *noopSearchEngine) Autocomplete(ctx context.Context, prefix string, size int) ([]domain.AutocompleteSuggestion, error) {
	return []domain.AutocompleteSuggestion{}, nil
}
func (e *noopSearchEngine) Index(ctx context.Context, ct domain.ContentType, doc *domain.SearchDocument) error {
	return nil
}
func (e *noopSearchEngine) Update(ctx context.Context, ct domain.ContentType, id uuid.UUID, doc map[string]interface{}) error {
	return nil
}
func (e *noopSearchEngine) Remove(ctx context.Context, ct domain.ContentType, id uuid.UUID) error {
	return nil
}
func (e *noopSearchEngine) Reindex(ctx context.Context, ct domain.ContentType, docs []domain.SearchDocument) error {
	return nil
}
func (e *noopSearchEngine) GetStats(ctx context.Context) ([]domain.IndexStats, error) {
	return []domain.IndexStats{}, nil
}
func (e *noopSearchEngine) Health(ctx context.Context) (*domain.ESHealth, error) {
	return &domain.ESHealth{Status: "green (noop)", NodeCount: 0, ActiveShards: 0}, nil
}
func (e *noopSearchEngine) PutSynonyms(ctx context.Context, set *domain.SynonymSet) error {
	return nil
}
func (e *noopSearchEngine) PutAnalyzers(ctx context.Context, config map[string]interface{}) error {
	return nil
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type noopCache struct{}
func NewNoopCache() Cache { return &noopCache{} }
func (c *noopCache) Get(ctx context.Context, key string) (string, error) { return "", pgx.ErrNoRows }
func (c *noopCache) Set(ctx context.Context, key, value string, ttl time.Duration) error { return nil }
func (c *noopCache) Del(ctx context.Context, key string) error { return nil }
