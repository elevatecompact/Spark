package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
)

type BookmarkRepository interface {
	Create(ctx context.Context, bookmark *domain.Bookmark) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Bookmark, error)
	ListByViewer(ctx context.Context, viewerID uuid.UUID, folder string, limit, offset int) ([]*domain.Bookmark, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) error
	CountByViewer(ctx context.Context, viewerID uuid.UUID) (int, error)
}

type WatchLaterRepository interface {
	Add(ctx context.Context, item *domain.WatchLaterItem) error
	ListByViewer(ctx context.Context, viewerID uuid.UUID) ([]*domain.WatchLaterItem, error)
	Remove(ctx context.Context, id uuid.UUID) error
	RemoveByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) error
	CountByViewer(ctx context.Context, viewerID uuid.UUID) (int, error)
	Reorder(ctx context.Context, id uuid.UUID, newPosition int) error
}

type bookmarkRepository struct {
	pool *pgxpool.Pool
}

func NewBookmarkRepository(pool *pgxpool.Pool) BookmarkRepository {
	return &bookmarkRepository{pool: pool}
}

func (r *bookmarkRepository) Create(ctx context.Context, bookmark *domain.Bookmark) error {
	query := `
		INSERT INTO bookmarks (id, viewer_id, content_id, note, folder, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		bookmark.ID, bookmark.ViewerID, bookmark.ContentID,
		bookmark.Note, bookmark.Folder, bookmark.CreatedAt,
	).Scan(&bookmark.CreatedAt)
	if err != nil {
		if isPGUniqueViolation(err) {
			return domain.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to create bookmark: %w", err)
	}
	return nil
}

func (r *bookmarkRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bookmark, error) {
	query := `SELECT id, viewer_id, content_id, note, folder, created_at FROM bookmarks WHERE id = $1`
	bookmark := &domain.Bookmark{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&bookmark.ID, &bookmark.ViewerID, &bookmark.ContentID,
		&bookmark.Note, &bookmark.Folder, &bookmark.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}
	return bookmark, nil
}

func (r *bookmarkRepository) ListByViewer(ctx context.Context, viewerID uuid.UUID, folder string, limit, offset int) ([]*domain.Bookmark, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var rows pgx.Rows
	var err error

	if folder != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, viewer_id, content_id, note, folder, created_at
			FROM bookmarks
			WHERE viewer_id = $1 AND folder = $2
			ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
			viewerID, folder, limit, offset)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, viewer_id, content_id, note, folder, created_at
			FROM bookmarks
			WHERE viewer_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			viewerID, limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []*domain.Bookmark
	for rows.Next() {
		b := &domain.Bookmark{}
		err := rows.Scan(&b.ID, &b.ViewerID, &b.ContentID, &b.Note, &b.Folder, &b.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, b)
	}
	if bookmarks == nil {
		bookmarks = []*domain.Bookmark{}
	}
	return bookmarks, nil
}

func (r *bookmarkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM bookmarks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *bookmarkRepository) DeleteByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM bookmarks WHERE viewer_id = $1 AND content_id = $2`, viewerID, contentID)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *bookmarkRepository) CountByViewer(ctx context.Context, viewerID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM bookmarks WHERE viewer_id = $1`, viewerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count bookmarks: %w", err)
	}
	return count, nil
}

type watchLaterRepository struct {
	pool *pgxpool.Pool
}

func NewWatchLaterRepository(pool *pgxpool.Pool) WatchLaterRepository {
	return &watchLaterRepository{pool: pool}
}

func (r *watchLaterRepository) Add(ctx context.Context, item *domain.WatchLaterItem) error {
	query := `
		INSERT INTO watch_later (id, viewer_id, content_id, position, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		item.ID, item.ViewerID, item.ContentID, item.Position, item.CreatedAt,
	).Scan(&item.CreatedAt)
	if err != nil {
		if isPGUniqueViolation(err) {
			return domain.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to add watch later: %w", err)
	}
	return nil
}

func (r *watchLaterRepository) ListByViewer(ctx context.Context, viewerID uuid.UUID) ([]*domain.WatchLaterItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, viewer_id, content_id, position, created_at
		FROM watch_later
		WHERE viewer_id = $1
		ORDER BY position ASC, created_at ASC`, viewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list watch later: %w", err)
	}
	defer rows.Close()

	var items []*domain.WatchLaterItem
	for rows.Next() {
		item := &domain.WatchLaterItem{}
		err := rows.Scan(&item.ID, &item.ViewerID, &item.ContentID, &item.Position, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watch later item: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []*domain.WatchLaterItem{}
	}
	return items, nil
}

func (r *watchLaterRepository) Remove(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM watch_later WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove watch later: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchLaterRepository) RemoveByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM watch_later WHERE viewer_id = $1 AND content_id = $2`, viewerID, contentID)
	if err != nil {
		return fmt.Errorf("failed to remove watch later: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchLaterRepository) CountByViewer(ctx context.Context, viewerID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM watch_later WHERE viewer_id = $1`, viewerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count watch later: %w", err)
	}
	return count, nil
}

func (r *watchLaterRepository) Reorder(ctx context.Context, id uuid.UUID, newPosition int) error {
	_, err := r.pool.Exec(ctx, `UPDATE watch_later SET position = $2 WHERE id = $1`, id, newPosition)
	if err != nil {
		return fmt.Errorf("failed to reorder watch later: %w", err)
	}
	return nil
}

func isPGUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key value violates unique constraint") || contains(err.Error(), "23505"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
