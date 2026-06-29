package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/events"
	"github.com/elevatecompact/spark/services/viewer-service/internal/repository"
)

type EngagementService interface {
	RateContent(ctx context.Context, viewerID, contentID uuid.UUID, score int) (*domain.Rating, error)
	ToggleReaction(ctx context.Context, viewerID, contentID uuid.UUID, reactionType domain.ReactionType) (*domain.Reaction, error)
	ReportContent(ctx context.Context, viewerID, contentID uuid.UUID, reportType domain.ReportType, description string) (*domain.Report, error)
}

type engagementService struct {
	ratingRepo   repository.RatingRepository
	reactionRepo repository.ReactionRepository
	reportRepo   repository.ReportRepository
	eventPub     events.EventProducer
}

func NewEngagementService(
	ratingRepo repository.RatingRepository,
	reactionRepo repository.ReactionRepository,
	reportRepo repository.ReportRepository,
	eventPub events.EventProducer,
) EngagementService {
	return &engagementService{
		ratingRepo:   ratingRepo,
		reactionRepo: reactionRepo,
		reportRepo:   reportRepo,
		eventPub:     eventPub,
	}
}

func (s *engagementService) RateContent(ctx context.Context, viewerID, contentID uuid.UUID, score int) (*domain.Rating, error) {
	if score < 1 || score > 5 {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "score must be between 1 and 5", 400)
	}

	now := time.Now().UTC()
	rating := &domain.Rating{
		ID:        uuid.New(),
		ViewerID:  viewerID,
		ContentID: contentID,
		Score:     score,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.ratingRepo.Upsert(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to save rating: %w", err)
	}

	if err := s.eventPub.PublishRatingSubmitted(ctx, rating); err != nil {
		return nil, fmt.Errorf("failed to publish rating event: %w", err)
	}

	return rating, nil
}

func (s *engagementService) ToggleReaction(ctx context.Context, viewerID, contentID uuid.UUID, reactionType domain.ReactionType) (*domain.Reaction, error) {
	if reactionType != domain.ReactionLike && reactionType != domain.ReactionDislike {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "reaction must be 'like' or 'dislike'", 400)
	}

	now := time.Now().UTC()
	reaction := &domain.Reaction{
		ID:        uuid.New(),
		ViewerID:  viewerID,
		ContentID: contentID,
		Type:      reactionType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := s.reactionRepo.Toggle(ctx, reaction)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle reaction: %w", err)
	}

	if result != nil {
		if err := s.eventPub.PublishReactionAdded(ctx, result); err != nil {
			return nil, fmt.Errorf("failed to publish reaction event: %w", err)
		}
	}

	return result, nil
}

func (s *engagementService) ReportContent(ctx context.Context, viewerID, contentID uuid.UUID, reportType domain.ReportType, description string) (*domain.Report, error) {
	switch reportType {
	case domain.ReportSpam, domain.ReportHarassment, domain.ReportCopyright, domain.ReportOther:
	default:
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, "invalid report type", 400)
	}

	report := &domain.Report{
		ID:          uuid.New(),
		ViewerID:    viewerID,
		ContentID:   contentID,
		Type:        reportType,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}
