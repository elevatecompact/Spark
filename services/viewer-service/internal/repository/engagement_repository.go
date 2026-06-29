package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
)

type RatingRepository interface {
	Upsert(ctx context.Context, rating *domain.Rating) error
	GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.Rating, error)
	GetAverageByContent(ctx context.Context, contentID uuid.UUID) (float64, int, error)
}

type ReactionRepository interface {
	Toggle(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.Reaction, error)
	CountByContent(ctx context.Context, contentID uuid.UUID) (likes int, dislikes int, err error)
}

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
}

type ratingRepository struct {
	pool *pgxpool.Pool
}

func NewRatingRepository(pool *pgxpool.Pool) RatingRepository {
	return &ratingRepository{pool: pool}
}

func (r *ratingRepository) Upsert(ctx context.Context, rating *domain.Rating) error {
	query := `
		INSERT INTO ratings (id, viewer_id, content_id, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (viewer_id, content_id) DO UPDATE SET
			score = $4, updated_at = $6
		RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		rating.ID, rating.ViewerID, rating.ContentID,
		rating.Score, rating.CreatedAt, rating.UpdatedAt,
	).Scan(&rating.CreatedAt, &rating.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert rating: %w", err)
	}
	return nil
}

func (r *ratingRepository) GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.Rating, error) {
	query := `SELECT id, viewer_id, content_id, score, created_at, updated_at FROM ratings WHERE viewer_id = $1 AND content_id = $2`
	rating := &domain.Rating{}
	err := r.pool.QueryRow(ctx, query, viewerID, contentID).Scan(
		&rating.ID, &rating.ViewerID, &rating.ContentID,
		&rating.Score, &rating.CreatedAt, &rating.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}
	return rating, nil
}

func (r *ratingRepository) GetAverageByContent(ctx context.Context, contentID uuid.UUID) (float64, int, error) {
	var avg float64
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(score), 0), COUNT(*) FROM ratings WHERE content_id = $1`,
		contentID,
	).Scan(&avg, &count)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}
	return avg, count, nil
}

type reactionRepository struct {
	pool *pgxpool.Pool
}

func NewReactionRepository(pool *pgxpool.Pool) ReactionRepository {
	return &reactionRepository{pool: pool}
}

func (r *reactionRepository) Toggle(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error) {
	existing, err := r.GetByViewerAndContent(ctx, reaction.ViewerID, reaction.ContentID)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing reaction: %w", err)
	}

	if err == nil {
		if existing.Type == reaction.Type {
			_, err := r.pool.Exec(ctx, `DELETE FROM reactions WHERE id = $1`, existing.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to remove reaction: %w", err)
			}
			return nil, nil
		}
		query := `UPDATE reactions SET type = $3, updated_at = NOW() WHERE id = $1 RETURNING id, viewer_id, content_id, type, created_at, updated_at`
		err = r.pool.QueryRow(ctx, query, existing.ID, reaction.ViewerID, reaction.Type).Scan(
			&existing.ID, &existing.ViewerID, &existing.ContentID,
			&existing.Type, &existing.CreatedAt, &existing.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update reaction: %w", err)
		}
		return existing, nil
	}

	query := `
		INSERT INTO reactions (id, viewer_id, content_id, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err = r.pool.QueryRow(ctx, query,
		reaction.ID, reaction.ViewerID, reaction.ContentID,
		reaction.Type, reaction.CreatedAt, reaction.UpdatedAt,
	).Scan(&reaction.CreatedAt, &reaction.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create reaction: %w", err)
	}
	return reaction, nil
}

func (r *reactionRepository) GetByViewerAndContent(ctx context.Context, viewerID, contentID uuid.UUID) (*domain.Reaction, error) {
	query := `SELECT id, viewer_id, content_id, type, created_at, updated_at FROM reactions WHERE viewer_id = $1 AND content_id = $2`
	reaction := &domain.Reaction{}
	err := r.pool.QueryRow(ctx, query, viewerID, contentID).Scan(
		&reaction.ID, &reaction.ViewerID, &reaction.ContentID,
		&reaction.Type, &reaction.CreatedAt, &reaction.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}
	return reaction, nil
}

func (r *reactionRepository) CountByContent(ctx context.Context, contentID uuid.UUID) (int, int, error) {
	var likes, dislikes int
	err := r.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM reactions WHERE content_id = $1`,
		contentID,
	).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count reactions: %w", err)
	}
	return likes, dislikes, nil
}

type reportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) ReportRepository {
	return &reportRepository{pool: pool}
}

func (r *reportRepository) Create(ctx context.Context, report *domain.Report) error {
	query := `
		INSERT INTO reports (id, viewer_id, content_id, type, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		report.ID, report.ViewerID, report.ContentID,
		report.Type, report.Description, report.CreatedAt,
	).Scan(&report.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}
	return nil
}
