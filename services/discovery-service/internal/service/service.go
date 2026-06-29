package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/discovery-service/internal/domain"
	"github.com/elevatecompact/spark/services/discovery-service/internal/repository"
)

type DiscoveryService struct {
	repo *repository.DiscoveryRepository
	evt  domain.EventProducer
}

func NewDiscoveryService(repo *repository.DiscoveryRepository, evt domain.EventProducer) *DiscoveryService {
	return &DiscoveryService{repo: repo, evt: evt}
}

func (s *DiscoveryService) GetHomeFeed(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	ids, err := s.repo.GetHomeFeedContentIDs(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	feedType := "anonymous"
	if userID != nil {
		feedType = "personalized"
	}
	s.evt.Publish(ctx, "discovery.feed.served", map[string]interface{}{
		"feedType": "home", "source": feedType, "contentIds": ids,
	})
	return ids, nil
}

func (s *DiscoveryService) GetTrendingFeed(ctx context.Context, timeframe string, limit int) ([]domain.TrendingItem, error) {
	return s.repo.GetTrendingContentIDs(ctx, limit)
}

func (s *DiscoveryService) GetCategoryFeed(ctx context.Context, slug string, limit, offset int) ([]uuid.UUID, error) {
	return s.repo.GetCategoryContentIDs(ctx, slug, limit, offset)
}

func (s *DiscoveryService) GetNewFeed(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	return s.repo.GetNewContentIDs(ctx, limit, offset)
}

func (s *DiscoveryService) GetRelatedFeed(ctx context.Context, contentID uuid.UUID, limit int) ([]uuid.UUID, error) {
	return s.repo.GetRelatedContentIDs(ctx, contentID, limit)
}

func (s *DiscoveryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *DiscoveryService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return s.repo.GetCategoryBySlug(ctx, slug)
}

func (s *DiscoveryService) GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	return s.repo.GetSubcategories(ctx, parentID)
}

func (s *DiscoveryService) ListCollections(ctx context.Context, featured bool) ([]domain.Collection, error) {
	return s.repo.ListCollections(ctx, featured)
}

func (s *DiscoveryService) GetCollection(ctx context.Context, id uuid.UUID) (*domain.Collection, error) {
	return s.repo.GetCollection(ctx, id)
}

func (s *DiscoveryService) CreateCollection(ctx context.Context, c *domain.Collection) (*domain.Collection, error) {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	if err := s.repo.CreateCollection(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *DiscoveryService) UpdateCollection(ctx context.Context, c *domain.Collection) error {
	return s.repo.UpdateCollection(ctx, c)
}

func (s *DiscoveryService) AddCollectionItem(ctx context.Context, collectionID, contentID uuid.UUID, sortOrder int) error {
	item := &domain.CollectionItem{
		CollectionID: collectionID,
		ContentID:    contentID,
		SortOrder:    sortOrder,
		AddedAt:      time.Now(),
	}
	return s.repo.AddCollectionItem(ctx, item)
}

func (s *DiscoveryService) RemoveCollectionItem(ctx context.Context, collectionID, contentID uuid.UUID) error {
	return s.repo.RemoveCollectionItem(ctx, collectionID, contentID)
}

func (s *DiscoveryService) GetEditorialPicks(ctx context.Context, pickType domain.PickType) ([]domain.EditorialPick, error) {
	return s.repo.GetEditorialPicks(ctx, pickType)
}

func (s *DiscoveryService) GetTrending(ctx context.Context, limit int) ([]domain.TrendingItem, error) {
	return s.repo.GetTrendingContentIDs(ctx, limit)
}

func (s *DiscoveryService) GetTrendingByCategory(ctx context.Context, slug string, limit int) ([]domain.TrendingItem, error) {
	return s.repo.GetTrendingContentIDsByCategory(ctx, slug, limit)
}

func (s *DiscoveryService) GetTrendingCreators(ctx context.Context, limit int) ([]uuid.UUID, error) {
	return s.repo.GetTrendingCreatorIDs(ctx, limit)
}

func (s *DiscoveryService) WarmFeedCache(ctx context.Context, feedTypes []string) error {
	log.Info().Strs("feedTypes", feedTypes).Msg("feed cache warmed (noop)")
	return nil
}

func (s *DiscoveryService) RefreshTrending(ctx context.Context) error {
	log.Info().Msg("trending refreshed (noop)")
	return nil
}

func (s *DiscoveryService) ReorderCategories(ctx context.Context, order []uuid.UUID) error {
	log.Info().Msg("categories reordered (noop)")
	return nil
}

func (s *DiscoveryService) GetSpotlight(ctx context.Context) ([]domain.EditorialPick, error) {
	return s.repo.GetEditorialPicks(ctx, domain.PickSpotlight)
}

func (s *DiscoveryService) GetHolidayPicks(ctx context.Context, campaign string) ([]domain.EditorialPick, error) {
	return s.repo.GetEditorialPicks(ctx, domain.PickHoliday)
}

func (s *DiscoveryService) GetStaffPicks(ctx context.Context) ([]domain.EditorialPick, error) {
	return s.repo.GetEditorialPicks(ctx, domain.PickStaffPick)
}
