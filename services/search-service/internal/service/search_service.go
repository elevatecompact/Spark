package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/search-service/internal/domain"
	"github.com/elevatecompact/spark/services/search-service/internal/events"
	"github.com/elevatecompact/spark/services/search-service/internal/repository"
)

type SearchService interface {
	Search(ctx context.Context, q *domain.SearchQuery) (*domain.SearchResult, error)
	Autocomplete(ctx context.Context, prefix string, size int) ([]domain.AutocompleteSuggestion, error)
	IndexDocument(ctx context.Context, contentType domain.ContentType, doc *domain.SearchDocument) error
	UpdateDocument(ctx context.Context, contentType domain.ContentType, id uuid.UUID, doc map[string]interface{}) error
	RemoveDocument(ctx context.Context, contentType domain.ContentType, id uuid.UUID) error
	Reindex(ctx context.Context, contentType domain.ContentType) error
	RecordSuggestionClick(ctx context.Context, userID uuid.UUID, suggestion string) error
	GetStats(ctx context.Context) ([]domain.IndexStats, error)
	PutSynonyms(ctx context.Context, set *domain.SynonymSet) error
	PutAnalyzers(ctx context.Context, config map[string]interface{}) error
	Health(ctx context.Context) (*domain.ESHealth, error)
}

type searchService struct {
	engine  repository.SearchEngine
	repo    repository.SearchRepository
	cache   repository.Cache
	eventPub events.EventProducer
}

func NewSearchService(engine repository.SearchEngine, repo repository.SearchRepository, cache repository.Cache, eventPub events.EventProducer) SearchService {
	return &searchService{engine: engine, repo: repo, cache: cache, eventPub: eventPub}
}

func (s *searchService) Search(ctx context.Context, q *domain.SearchQuery) (*domain.SearchResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Size < 1 || q.Size > 100 {
		q.Size = 20
	}
	if q.Sort == "" {
		q.Sort = domain.SortRelevance
	}

	start := time.Now()
	result, err := s.engine.Search(ctx, q)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return nil, err
	}

	topIDs := make([]string, 0, len(result.Results))
	for _, r := range result.Results {
		topIDs = append(topIDs, r.ID.String())
	}

	s.eventPub.PublishQueryExecuted(ctx, &events.SearchQueryExecutedEvent{
		QueryID:   uuid.New().String(),
		Query:     q.Query,
		Filters:   q.Filters,
		ResultCount: len(result.Results),
		TopResultIDs: topIDs,
		LatencyMs: latency,
		UserID:    q.UserID,
		Timestamp: time.Now().UTC(),
	})

	s.repo.LogSearch(ctx, &domain.SearchAnalytics{
		Query:     q.Query,
		ResultIDs: topIDs,
		LatencyMs: latency,
		UserID:    q.UserID,
		Timestamp: time.Now().UTC(),
	})

	return result, nil
}

func (s *searchService) Autocomplete(ctx context.Context, prefix string, size int) ([]domain.AutocompleteSuggestion, error) {
	if size <= 0 || size > 20 {
		size = 10
	}
	return s.engine.Autocomplete(ctx, prefix, size)
}

func (s *searchService) IndexDocument(ctx context.Context, contentType domain.ContentType, doc *domain.SearchDocument) error {
	if err := s.engine.Index(ctx, contentType, doc); err != nil {
		return err
	}
	return s.eventPub.PublishIndexUpdated(ctx, contentType, doc.ID)
}

func (s *searchService) UpdateDocument(ctx context.Context, contentType domain.ContentType, id uuid.UUID, doc map[string]interface{}) error {
	return s.engine.Update(ctx, contentType, id, doc)
}

func (s *searchService) RemoveDocument(ctx context.Context, contentType domain.ContentType, id uuid.UUID) error {
	return s.engine.Remove(ctx, contentType, id)
}

func (s *searchService) Reindex(ctx context.Context, contentType domain.ContentType) error {
	log.Info().Str("type", string(contentType)).Msg("reindex triggered (noop)")
	return nil
}

func (s *searchService) RecordSuggestionClick(ctx context.Context, userID uuid.UUID, suggestion string) error {
	return s.eventPub.PublishSuggestionClicked(ctx, userID, suggestion)
}

func (s *searchService) GetStats(ctx context.Context) ([]domain.IndexStats, error) {
	if cached, err := s.cache.Get(ctx, "search:stats"); err == nil && cached != "" {
		return []domain.IndexStats{}, nil
	}
	stats, err := s.engine.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *searchService) PutSynonyms(ctx context.Context, set *domain.SynonymSet) error {
	return s.engine.PutSynonyms(ctx, set)
}

func (s *searchService) PutAnalyzers(ctx context.Context, config map[string]interface{}) error {
	return s.engine.PutAnalyzers(ctx, config)
}

func (s *searchService) Health(ctx context.Context) (*domain.ESHealth, error) {
	return s.engine.Health(ctx)
}
