package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
)

type DiscoverService struct {
	creatorRepo  repository.CreatorRepository
	categoryRepo repository.CategoryRepository
}

func NewDiscoverService(creatorRepo repository.CreatorRepository, categoryRepo repository.CategoryRepository) *DiscoverService {
	return &DiscoverService{
		creatorRepo:  creatorRepo,
		categoryRepo: categoryRepo,
	}
}

type SearchFilters struct {
	Query      string
	Category   string
	Language   string
	Country    string
	Tags       []string
	Limit      int
	Offset     int
}

func (s *DiscoverService) Search(ctx context.Context, filters SearchFilters) ([]domain.Creator, int, error) {
	var categoryIDs []string
	if filters.Category != "" {
		cat, err := s.categoryRepo.GetBySlug(ctx, filters.Category)
		if err != nil {
			catByID, err2 := s.categoryRepo.GetByID(ctx, uuid.MustParse(filters.Category))
			if err2 != nil {
				return nil, 0, domain.ErrCategoryNotFound
			}
			categoryIDs = []string{catByID.ID.String()}
		} else {
			categoryIDs = []string{cat.ID.String()}
		}
	}

	return s.creatorRepo.Search(ctx, filters.Query, categoryIDs, filters.Tags, filters.Language, filters.Country, filters.Limit, filters.Offset)
}

func (s *DiscoverService) GetTrending(ctx context.Context, limit, offset int) ([]domain.Creator, error) {
	return s.creatorRepo.GetTrending(ctx, limit, offset)
}

func (s *DiscoverService) GetRecommended(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Creator, error) {
	return s.creatorRepo.GetRecommended(ctx, userID, limit)
}

func (s *DiscoverService) GetNearby(ctx context.Context, country string, limit, offset int) ([]domain.Creator, error) {
	if country == "" {
		return s.creatorRepo.GetTrending(ctx, limit, offset)
	}

	creators, _, err := s.creatorRepo.Search(ctx, "", nil, nil, "", country, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get nearby: %w", err)
	}
	return creators, nil
}

func (s *DiscoverService) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Creator, int, error) {
	return s.creatorRepo.GetByCategoryID(ctx, categoryID, limit, offset)
}
