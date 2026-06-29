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

type WatchHistoryService interface {
	RecordWatch(ctx context.Context, viewerID uuid.UUID, update domain.WatchProgressUpdate) (*domain.WatchHistory, error)
	GetHistory(ctx context.Context, viewerID uuid.UUID, contentType string, days int, limit, offset int) ([]*domain.WatchHistory, error)
	DeleteEntry(ctx context.Context, viewerID, entryID uuid.UUID) error
	ClearHistory(ctx context.Context, viewerID uuid.UUID) error
}

type watchHistoryService struct {
	repo     repository.WatchHistoryRepository
	eventPub events.EventProducer
}

func NewWatchHistoryService(repo repository.WatchHistoryRepository, eventPub events.EventProducer) WatchHistoryService {
	return &watchHistoryService{
		repo:     repo,
		eventPub: eventPub,
	}
}

func (s *watchHistoryService) RecordWatch(ctx context.Context, viewerID uuid.UUID, update domain.WatchProgressUpdate) (*domain.WatchHistory, error) {
	existing, err := s.repo.GetByViewerAndContent(ctx, viewerID, update.ContentID)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing entry: %w", err)
	}

	if err == nil {
		newProgress := update.Progress
		if update.Progress < existing.Progress {
			newProgress = existing.Progress
		}
		newDuration := existing.WatchDurationSeconds + update.WatchDurationSeconds
		completed := existing.Completed || update.Completed

		if err := s.repo.UpdateProgress(ctx, existing.ID, newProgress, newDuration, completed); err != nil {
			return nil, fmt.Errorf("failed to update progress: %w", err)
		}

		existing.Progress = newProgress
		existing.WatchDurationSeconds = newDuration
		existing.Completed = completed

		if err := s.eventPub.PublishWatchProgress(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to publish watch event: %w", err)
		}

		return existing, nil
	}

	entry := &domain.WatchHistory{
		ID:                   uuid.New(),
		ViewerID:             viewerID,
		ContentID:            update.ContentID,
		ContentType:          update.ContentType,
		Progress:             update.Progress,
		WatchDurationSeconds: update.WatchDurationSeconds,
		Completed:            update.Completed,
		WatchedAt:            time.Now().UTC(),
		CreatedAt:            time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to create watch entry: %w", err)
	}

	if err := s.eventPub.PublishWatchStarted(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to publish watch event: %w", err)
	}

	if entry.Completed {
		if err := s.eventPub.PublishWatchCompleted(ctx, entry); err != nil {
			return nil, fmt.Errorf("failed to publish watch event: %w", err)
		}
	}

	return entry, nil
}

func (s *watchHistoryService) GetHistory(ctx context.Context, viewerID uuid.UUID, contentType string, days int, limit, offset int) ([]*domain.WatchHistory, error) {
	if days <= 0 {
		days = 90
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListByViewer(ctx, viewerID, contentType, days, limit, offset)
}

func (s *watchHistoryService) DeleteEntry(ctx context.Context, viewerID, entryID uuid.UUID) error {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return err
	}
	if entry.ViewerID != viewerID {
		return domain.ErrForbidden
	}
	return s.repo.Delete(ctx, entryID)
}

func (s *watchHistoryService) ClearHistory(ctx context.Context, viewerID uuid.UUID) error {
	return s.repo.DeleteByViewer(ctx, viewerID)
}
