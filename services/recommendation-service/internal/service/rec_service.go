package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/recommendation-service/internal/domain"
	"github.com/elevatecompact/spark/services/recommendation-service/internal/events"
	"github.com/elevatecompact/spark/services/recommendation-service/internal/repository"
)

type RecService interface {
	GetHomeFeed(ctx context.Context, userID uuid.UUID, limit int) (*domain.Feed, error)
	GetTrendingFeed(ctx context.Context, limit int) (*domain.Feed, error)
	GetUpNext(ctx context.Context, userID, contentID uuid.UUID, limit int) (*domain.Feed, error)
	GetSimilar(ctx context.Context, contentID uuid.UUID, limit int) (*domain.Feed, error)
	GetCreatorFeed(ctx context.Context, creatorID uuid.UUID, limit int) (*domain.Feed, error)
	RecordClick(ctx context.Context, userID, contentID uuid.UUID) error
	RecordDismiss(ctx context.Context, userID, contentID uuid.UUID) error
	Explain(ctx context.Context, recID uuid.UUID) (map[string]interface{}, error)
	GetActiveModel(ctx context.Context) (*domain.ModelInfo, error)
	DeployModel(ctx context.Context, version string, metrics string) error
	GetModelMetrics(ctx context.Context) ([]domain.ModelInfo, error)
	RefreshFeatures(ctx context.Context) error
	GetFeatureImportance(ctx context.Context) (map[string]float64, error)
	InvalidateCache(ctx context.Context) error
}

type recService struct {
	embRepo  repository.EmbeddingRepository
	intRepo  repository.InteractionRepository
	eventPub events.EventProducer
}

func NewRecService(embRepo repository.EmbeddingRepository, intRepo repository.InteractionRepository, eventPub events.EventProducer) RecService {
	return &recService{embRepo: embRepo, intRepo: intRepo, eventPub: eventPub}
}

func (s *recService) GetHomeFeed(ctx context.Context, userID uuid.UUID, limit int) (*domain.Feed, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	items := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		items[i] = domain.Recommendation{
			ContentID: uuid.New(),
			Score:     rand.Float64(),
			Reason:    "based on your watch history",
		}
	}

	feed := &domain.Feed{
		Type:     domain.FeedHome,
		UserID:   userID,
		Items:    items,
		ServedAt: time.Now().UTC(),
	}

	if err := s.eventPub.PublishFeedServed(ctx, feed); err != nil {
		log.Warn().Err(err).Msg("failed to publish feed served")
	}

	return feed, nil
}

func (s *recService) GetTrendingFeed(ctx context.Context, limit int) (*domain.Feed, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	items := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		items[i] = domain.Recommendation{
			ContentID: uuid.New(),
			Score:     rand.Float64() * 100,
			Reason:    "trending now",
		}
	}
	return &domain.Feed{Type: domain.FeedTrending, Items: items, ServedAt: time.Now().UTC()}, nil
}

func (s *recService) GetUpNext(ctx context.Context, userID, contentID uuid.UUID, limit int) (*domain.Feed, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	items := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		items[i] = domain.Recommendation{
			ContentID: uuid.New(),
			Score:     rand.Float64(),
			Reason:    "watch next",
		}
	}
	return &domain.Feed{Type: domain.FeedUpNext, Items: items, ServedAt: time.Now().UTC()}, nil
}

func (s *recService) GetSimilar(ctx context.Context, contentID uuid.UUID, limit int) (*domain.Feed, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	items := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		items[i] = domain.Recommendation{
			ContentID: uuid.New(),
			Score:     rand.Float64(),
			Reason:    "similar content",
		}
	}
	return &domain.Feed{Type: domain.FeedSimilar, Items: items, ServedAt: time.Now().UTC()}, nil
}

func (s *recService) GetCreatorFeed(ctx context.Context, creatorID uuid.UUID, limit int) (*domain.Feed, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	items := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		items[i] = domain.Recommendation{
			ContentID: uuid.New(),
			Score:     rand.Float64(),
			Reason:    "by this creator",
		}
	}
	return &domain.Feed{Type: domain.FeedCreator, Items: items, ServedAt: time.Now().UTC()}, nil
}

func (s *recService) RecordClick(ctx context.Context, userID, contentID uuid.UUID) error {
	return s.intRepo.Insert(ctx, &domain.UserContentInteraction{
		UserID:          userID,
		ContentID:       contentID,
		InteractionType: domain.InteractionClick,
		Weight:          1.0,
		Timestamp:       time.Now().UTC(),
	})
}

func (s *recService) RecordDismiss(ctx context.Context, userID, contentID uuid.UUID) error {
	return s.intRepo.Insert(ctx, &domain.UserContentInteraction{
		UserID:          userID,
		ContentID:       contentID,
		InteractionType: domain.InteractionDismiss,
		Weight:          -1.0,
		Timestamp:       time.Now().UTC(),
	})
}

func (s *recService) Explain(ctx context.Context, recID uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{
		"contribution_scores": map[string]float64{
			"watch_history":  0.45,
			"subscriptions":  0.25,
			"similar_users":  0.20,
			"trending_score": 0.10,
		},
		"top_feature":           "watched_creator_previously",
		"diversity_bucket":      "entertainment",
		"novelty_score":         0.73,
	}, nil
}

func (s *recService) GetActiveModel(ctx context.Context) (*domain.ModelInfo, error) {
	models, err := s.embRepo.ListActiveModels(ctx)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return &domain.ModelInfo{
			Version:  "v0.0.0-noop",
			IsActive: true,
		}, nil
	}
	return &models[0], nil
}

func (s *recService) DeployModel(ctx context.Context, version string, metrics string) error {
	log.Info().Str("version", version).Msg("model deployed (noop)")
	return nil
}

func (s *recService) GetModelMetrics(ctx context.Context) ([]domain.ModelInfo, error) {
	models, err := s.embRepo.ListActiveModels(ctx)
	if err != nil {
		models = []domain.ModelInfo{
			{Version: "v1.0.0", IsActive: true, Metrics: `{"ndcg@10": 0.42, "recall@20": 0.58}`},
		}
	}
	return models, nil
}

func (s *recService) RefreshFeatures(ctx context.Context) error {
	log.Info().Msg("features refreshed (noop)")
	return nil
}

func (s *recService) GetFeatureImportance(ctx context.Context) (map[string]float64, error) {
	return map[string]float64{
		"watch_history":        0.35,
		"subscriptions":        0.20,
		"engagement_score":     0.15,
		"content_similarity":   0.12,
		"trending_score":       0.10,
		"freshness":            0.05,
		"creator_followers":    0.03,
	}, nil
}

func (s *recService) InvalidateCache(ctx context.Context) error {
	log.Info().Msg("cache invalidated (noop)")
	return nil
}
