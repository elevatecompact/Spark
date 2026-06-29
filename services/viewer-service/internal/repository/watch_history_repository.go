package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
)

type WatchHistoryRepository interface {
	Create(ctx context.Context, entry *domain.WatchHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WatchHistory, error)
	ListByViewer(ctx context.Context, viewerID uuid.UUID, contentType string, days int, limit, offset int) ([]*domain.WatchHistory, error)
	UpdateProgress(ctx context.Context, id uuid.UUID, progress float64, duration int, completed bool) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByViewer(ctx context.Context, viewerID uuid.UUID) error
	DeleteOlderThan(ctx context.Context, days int) error
	GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.WatchHistory, error)
}

type watchHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewWatchHistoryRepository(pool *pgxpool.Pool) WatchHistoryRepository {
	return &watchHistoryRepository{pool: pool}
}

func (r *watchHistoryRepository) Create(ctx context.Context, entry *domain.WatchHistory) error {
	query := `
		INSERT INTO watch_history (id, viewer_id, content_id, content_type, progress, watch_duration_seconds, completed, watched_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		entry.ID,
		entry.ViewerID,
		entry.ContentID,
		entry.ContentType,
		entry.Progress,
		entry.WatchDurationSeconds,
		entry.Completed,
		entry.WatchedAt,
		entry.CreatedAt,
	).Scan(&entry.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create watch history: %w", err)
	}
	return nil
}

func (r *watchHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WatchHistory, error) {
	query := `
		SELECT id, viewer_id, content_id, content_type, progress, watch_duration_seconds, completed, watched_at, created_at
		FROM watch_history WHERE id = $1`

	entry := &domain.WatchHistory{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&entry.ID, &entry.ViewerID, &entry.ContentID, &entry.ContentType,
		&entry.Progress, &entry.WatchDurationSeconds, &entry.Completed,
		&entry.WatchedAt, &entry.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get watch history: %w", err)
	}
	return entry, nil
}

func (r *watchHistoryRepository) ListByViewer(ctx context.Context, viewerID uuid.UUID, contentType string, days int, limit, offset int) ([]*domain.WatchHistory, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var rows pgx.Rows
	var err error

	since := time.Now().AddDate(0, 0, -days)

	if contentType != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, viewer_id, content_id, content_type, progress, watch_duration_seconds, completed, watched_at, created_at
			FROM watch_history
			WHERE viewer_id = $1 AND content_type = $2 AND watched_at >= $3
			ORDER BY watched_at DESC LIMIT $4 OFFSET $5`,
			viewerID, contentType, since, limit, offset)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, viewer_id, content_id, content_type, progress, watch_duration_seconds, completed, watched_at, created_at
			FROM watch_history
			WHERE viewer_id = $1 AND watched_at >= $2
			ORDER BY watched_at DESC LIMIT $3 OFFSET $4`,
			viewerID, since, limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list watch history: %w", err)
	}
	defer rows.Close()

	var entries []*domain.WatchHistory
	for rows.Next() {
		entry := &domain.WatchHistory{}
		err := rows.Scan(
			&entry.ID, &entry.ViewerID, &entry.ContentID, &entry.ContentType,
			&entry.Progress, &entry.WatchDurationSeconds, &entry.Completed,
			&entry.WatchedAt, &entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watch history: %w", err)
		}
		entries = append(entries, entry)
	}
	if entries == nil {
		entries = []*domain.WatchHistory{}
	}
	return entries, nil
}

func (r *watchHistoryRepository) UpdateProgress(ctx context.Context, id uuid.UUID, progress float64, duration int, completed bool) error {
	query := `UPDATE watch_history SET progress = $2, watch_duration_seconds = $3, completed = $4 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, progress, duration, completed)
	if err != nil {
		return fmt.Errorf("failed to update watch progress: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM watch_history WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete watch history: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchHistoryRepository) DeleteByViewer(ctx context.Context, viewerID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM watch_history WHERE viewer_id = $1`, viewerID)
	if err != nil {
		return fmt.Errorf("failed to delete viewer watch history: %w", err)
	}
	return nil
}

func (r *watchHistoryRepository) DeleteOlderThan(ctx context.Context, days int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM watch_history WHERE watched_at < NOW() - INTERVAL '1 day' * $1`, days)
	if err != nil {
		return fmt.Errorf("failed to delete old watch history: %w", err)
	}
	return nil
}

func (r *watchHistoryRepository) GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.WatchHistory, error) {
	query := `
		SELECT id, viewer_id, content_id, content_type, progress, watch_duration_seconds, completed, watched_at, created_at
		FROM watch_history
		WHERE viewer_id = $1 AND content_id = $2
		ORDER BY watched_at DESC LIMIT 1`

	entry := &domain.WatchHistory{}
	err := r.pool.QueryRow(ctx, query, viewerID, contentID).Scan(
		&entry.ID, &entry.ViewerID, &entry.ContentID, &entry.ContentType,
		&entry.Progress, &entry.WatchDurationSeconds, &entry.Completed,
		&entry.WatchedAt, &entry.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get watch history by viewer and content: %w", err)
	}
	return entry, nil
}
