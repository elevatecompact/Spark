package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/events"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
)

type CreatorService struct {
	creatorRepo repository.CreatorRepository
	producer    *events.Producer
	cache       CacheProvider
}

type CacheProvider interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

func NewCreatorService(creatorRepo repository.CreatorRepository, producer *events.Producer, cache CacheProvider) *CreatorService {
	return &CreatorService{
		creatorRepo: creatorRepo,
		producer:    producer,
		cache:       cache,
	}
}

func (s *CreatorService) CreateProfile(ctx context.Context, userID uuid.UUID, req domain.CreateCreatorRequest) (*domain.Creator, error) {
	existing, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil && err != domain.ErrCreatorNotFound {
		return nil, fmt.Errorf("check existing creator: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrCreatorAlreadyExists
	}

	now := time.Now()
	creator := &domain.Creator{
		ID:          uuid.New(),
		UserID:      userID,
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		Language:    req.Language,
		Country:     req.Country,
		Timezone:    req.Timezone,
		Status:      domain.CreatorActive,
		Categories:  []string{},
		Tags:        []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.creatorRepo.Create(ctx, creator); err != nil {
		return nil, fmt.Errorf("create creator: %w", err)
	}

	if err := s.producer.CreatorCreated(ctx, creator.ID.String(), creator.UserID.String(), creator.DisplayName); err != nil {
		log.Warn().Err(err).Msg("failed to emit CreatorCreated event")
	}

	log.Info().Str("creator_id", creator.ID.String()).Str("user_id", userID.String()).Msg("Creator profile created")
	return creator, nil
}

func (s *CreatorService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.Creator, error) {
	creator, err := s.creatorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return creator, nil
}

func (s *CreatorService) UpdateProfile(ctx context.Context, id uuid.UUID, req domain.UpdateCreatorRequest) error {
	creator, err := s.creatorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	changes := make(map[string]interface{})

	if req.DisplayName != nil {
		if len(*req.DisplayName) < 2 || len(*req.DisplayName) > 100 {
			return domain.ErrInvalidInput
		}
		creator.DisplayName = *req.DisplayName
		changes["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		if len(*req.Bio) > 2000 {
			return domain.ErrInvalidInput
		}
		creator.Bio = *req.Bio
		changes["bio"] = *req.Bio
	}
	if req.AvatarURL != nil {
		creator.AvatarURL = *req.AvatarURL
		changes["avatar_url"] = *req.AvatarURL
	}
	if req.BannerURL != nil {
		creator.BannerURL = *req.BannerURL
		changes["banner_url"] = *req.BannerURL
	}
	if req.Categories != nil {
		creator.Categories = *req.Categories
		changes["categories"] = *req.Categories
	}
	if req.Tags != nil {
		creator.Tags = *req.Tags
		changes["tags"] = *req.Tags
	}
	if req.Language != nil {
		if len(*req.Language) != 2 {
			return domain.ErrInvalidInput
		}
		creator.Language = *req.Language
		changes["language"] = *req.Language
	}
	if req.Country != nil {
		if len(*req.Country) != 2 {
			return domain.ErrInvalidInput
		}
		creator.Country = *req.Country
		changes["country"] = *req.Country
	}
	if req.Timezone != nil {
		creator.Timezone = *req.Timezone
		changes["timezone"] = *req.Timezone
	}
	if req.SocialLinks != nil {
		creator.SocialLinks = *req.SocialLinks
		changes["social_links"] = *req.SocialLinks
	}

	creator.UpdatedAt = time.Now()

	if err := s.creatorRepo.Update(ctx, creator); err != nil {
		return fmt.Errorf("update creator: %w", err)
	}

	if len(changes) > 0 {
		if err := s.producer.CreatorUpdated(ctx, creator.ID.String(), changes); err != nil {
			log.Warn().Err(err).Msg("failed to emit CreatorUpdated event")
		}
	}

	cacheKey := fmt.Sprintf("creator:%s", id.String())
	_ = s.cache.Delete(ctx, cacheKey)

	log.Info().Str("creator_id", id.String()).Msg("Creator profile updated")
	return nil
}

func (s *CreatorService) SearchCreators(ctx context.Context, query string, categories, tags []string, language, country string, limit, offset int) ([]domain.Creator, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.creatorRepo.Search(ctx, query, categories, tags, language, country, limit, offset)
}

func (s *CreatorService) FollowCreator(ctx context.Context, followerID, creatorID uuid.UUID) error {
	if followerID == creatorID {
		return domain.ErrSelfFollow
	}

	creator, err := s.creatorRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}
	if creator.Status != domain.CreatorActive {
		return domain.ErrCreatorNotFound
	}

	already, err := s.creatorRepo.IsFollowing(ctx, followerID, creatorID)
	if err != nil {
		return fmt.Errorf("check following: %w", err)
	}
	if already {
		return nil
	}

	if err := s.creatorRepo.AddFollower(ctx, followerID, creatorID); err != nil {
		return fmt.Errorf("add follower: %w", err)
	}
	if err := s.creatorRepo.IncrementFollowers(ctx, creatorID, 1); err != nil {
		log.Warn().Err(err).Msg("failed to increment follower count")
	}

	if err := s.producer.CreatorFollowed(ctx, followerID.String(), creatorID.String()); err != nil {
		log.Warn().Err(err).Msg("failed to emit CreatorFollowed event")
	}
	return nil
}

func (s *CreatorService) UnfollowCreator(ctx context.Context, followerID, creatorID uuid.UUID) error {
	already, err := s.creatorRepo.IsFollowing(ctx, followerID, creatorID)
	if err != nil {
		return fmt.Errorf("check following: %w", err)
	}
	if !already {
		return nil
	}

	if err := s.creatorRepo.RemoveFollower(ctx, followerID, creatorID); err != nil {
		return fmt.Errorf("remove follower: %w", err)
	}
	if err := s.creatorRepo.IncrementFollowers(ctx, creatorID, -1); err != nil {
		log.Warn().Err(err).Msg("failed to decrement follower count")
	}
	return nil
}

func (s *CreatorService) GetFollowers(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]uuid.UUID, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.creatorRepo.GetFollowers(ctx, creatorID, limit, offset)
}

func (s *CreatorService) GetFollowing(ctx context.Context, followerID uuid.UUID, limit, offset int) ([]domain.Creator, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.creatorRepo.GetFollowing(ctx, followerID, limit, offset)
}

func (s *CreatorService) GetTrending(ctx context.Context, limit, offset int) ([]domain.Creator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.creatorRepo.GetTrending(ctx, limit, offset)
}

func (s *CreatorService) GetRecommended(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Creator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.creatorRepo.GetRecommended(ctx, userID, limit)
}

func (s *CreatorService) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Creator, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.creatorRepo.GetByCategoryID(ctx, categoryID, limit, offset)
}

func (s *CreatorService) IsFollowing(ctx context.Context, followerID, creatorID uuid.UUID) (bool, error) {
	return s.creatorRepo.IsFollowing(ctx, followerID, creatorID)
}
