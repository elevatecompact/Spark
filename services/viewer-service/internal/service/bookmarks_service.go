package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/repository"
)

type BookmarkService interface {
	Create(ctx context.Context, viewerID uuid.UUID, contentID uuid.UUID, note, folder string) (*domain.Bookmark, error)
	List(ctx context.Context, viewerID uuid.UUID, folder string, limit, offset int) ([]*domain.Bookmark, error)
	Delete(ctx context.Context, viewerID, bookmarkID uuid.UUID) error
}

type WatchLaterService interface {
	Add(ctx context.Context, viewerID uuid.UUID, contentID uuid.UUID) (*domain.WatchLaterItem, error)
	List(ctx context.Context, viewerID uuid.UUID) ([]*domain.WatchLaterItem, error)
	Remove(ctx context.Context, viewerID, itemID uuid.UUID) error
	Reorder(ctx context.Context, viewerID, itemID uuid.UUID, newPosition int) error
}

type bookmarkService struct {
	repo        repository.BookmarkRepository
	maxBookmarks int
}

type watchLaterService struct {
	repo          repository.WatchLaterRepository
	maxWatchLater int
}

func NewBookmarkService(repo repository.BookmarkRepository, maxBookmarks int) BookmarkService {
	return &bookmarkService{
		repo:         repo,
		maxBookmarks: maxBookmarks,
	}
}

func NewWatchLaterService(repo repository.WatchLaterRepository, maxWatchLater int) WatchLaterService {
	return &watchLaterService{
		repo:          repo,
		maxWatchLater: maxWatchLater,
	}
}

func (s *bookmarkService) Create(ctx context.Context, viewerID uuid.UUID, contentID uuid.UUID, note, folder string) (*domain.Bookmark, error) {
	count, err := s.repo.CountByViewer(ctx, viewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to count bookmarks: %w", err)
	}
	if count >= s.maxBookmarks {
		return nil, domain.ErrMaxBookmarksReached
	}

	bookmark := &domain.Bookmark{
		ID:        uuid.New(),
		ViewerID:  viewerID,
		ContentID: contentID,
		Note:      note,
		Folder:    folder,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, bookmark); err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (s *bookmarkService) List(ctx context.Context, viewerID uuid.UUID, folder string, limit, offset int) ([]*domain.Bookmark, error) {
	return s.repo.ListByViewer(ctx, viewerID, folder, limit, offset)
}

func (s *bookmarkService) Delete(ctx context.Context, viewerID, bookmarkID uuid.UUID) error {
	bookmark, err := s.repo.GetByID(ctx, bookmarkID)
	if err != nil {
		return err
	}
	if bookmark.ViewerID != viewerID {
		return domain.ErrForbidden
	}
	return s.repo.Delete(ctx, bookmarkID)
}

func (s *watchLaterService) Add(ctx context.Context, viewerID uuid.UUID, contentID uuid.UUID) (*domain.WatchLaterItem, error) {
	count, err := s.repo.CountByViewer(ctx, viewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to count watch later: %w", err)
	}
	if count >= s.maxWatchLater {
		return nil, domain.ErrMaxWatchLaterReached
	}

	item := &domain.WatchLaterItem{
		ID:        uuid.New(),
		ViewerID:  viewerID,
		ContentID: contentID,
		Position:  count + 1,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Add(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *watchLaterService) List(ctx context.Context, viewerID uuid.UUID) ([]*domain.WatchLaterItem, error) {
	return s.repo.ListByViewer(ctx, viewerID)
}

func (s *watchLaterService) Remove(ctx context.Context, viewerID, itemID uuid.UUID) error {
	item, err := s.repo.ListByViewer(ctx, viewerID)
	if err != nil {
		return err
	}

	found := false
	for _, i := range item {
		if i.ID == itemID {
			found = true
			break
		}
	}
	if !found {
		return domain.ErrNotFound
	}

	return s.repo.Remove(ctx, itemID)
}

func (s *watchLaterService) Reorder(ctx context.Context, viewerID, itemID uuid.UUID, newPosition int) error {
	items, err := s.repo.ListByViewer(ctx, viewerID)
	if err != nil {
		return err
	}

	found := false
	for _, i := range items {
		if i.ID == itemID {
			found = true
			break
		}
	}
	if !found {
		return domain.ErrNotFound
	}

	if newPosition < 1 {
		newPosition = 1
	}
	if newPosition > len(items) {
		newPosition = len(items)
	}

	return s.repo.Reorder(ctx, itemID, newPosition)
}
