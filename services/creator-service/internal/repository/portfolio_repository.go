package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
)

type PortfolioRepository interface {
	Create(ctx context.Context, item *domain.PortfolioItem) error
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID) ([]domain.PortfolioItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PortfolioItem, error)
	Update(ctx context.Context, item *domain.PortfolioItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetFeatured(ctx context.Context, id, creatorID uuid.UUID) error
}

type portfolioRepository struct {
	pool *pgxpool.Pool
}

func NewPortfolioRepository(pool *pgxpool.Pool) PortfolioRepository {
	return &portfolioRepository{pool: pool}
}

func (r *portfolioRepository) Create(ctx context.Context, item *domain.PortfolioItem) error {
	query := INSERT INTO portfolio_items (id, creator_id, title, description, media_url, media_type, thumbnail_url, featured, sort_order, created_at)
		VALUES (, , , , , , , , , )
	_, err := r.pool.Exec(ctx, query,
		item.ID, item.CreatorID, item.Title, item.Description, item.MediaURL,
		item.MediaType, item.ThumbnailURL, item.Featured, item.SortOrder, item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert portfolio item: %w", err)
	}
	return nil
}

func (r *portfolioRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PortfolioItem, error) {
	row := r.pool.QueryRow(ctx, SELECT * FROM portfolio_items WHERE id = , id)
	return scanPortfolioItem(row)
}

func (r *portfolioRepository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID) ([]domain.PortfolioItem, error) {
	rows, err := r.pool.Query(ctx, SELECT * FROM portfolio_items WHERE creator_id =  ORDER BY featured DESC, sort_order ASC, created_at DESC, creatorID)
	if err != nil {
		return nil, fmt.Errorf("get portfolio items: %w", err)
	}
	defer rows.Close()

	var items []domain.PortfolioItem
	for rows.Next() {
		item, err := scanPortfolioItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if items == nil {
		items = []domain.PortfolioItem{}
	}
	return items, nil
}

func (r *portfolioRepository) Update(ctx context.Context, item *domain.PortfolioItem) error {
	query := UPDATE portfolio_items SET title=, description=, media_url=, media_type=, thumbnail_url=, featured=, sort_order= WHERE id=
	_, err := r.pool.Exec(ctx, query,
		item.ID, item.Title, item.Description, item.MediaURL, item.MediaType,
		item.ThumbnailURL, item.Featured, item.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("update portfolio item: %w", err)
	}
	return nil
}

func (r *portfolioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, DELETE FROM portfolio_items WHERE id = , id)
	if err != nil {
		return fmt.Errorf("delete portfolio item: %w", err)
	}
	return nil
}

func (r *portfolioRepository) SetFeatured(ctx context.Context, id, creatorID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, UPDATE portfolio_items SET featured = false WHERE creator_id = , creatorID); err != nil {
		return fmt.Errorf("unset featured: %w", err)
	}
	if _, err := tx.Exec(ctx, UPDATE portfolio_items SET featured = true WHERE id =  AND creator_id = , id, creatorID); err != nil {
		return fmt.Errorf("set featured: %w", err)
	}

	return tx.Commit(ctx)
}

func scanPortfolioItem(s interface{ Scan(dest ...interface{}) error }) (*domain.PortfolioItem, error) {
	item := &domain.PortfolioItem{}
	err := s.Scan(
		&item.ID, &item.CreatorID, &item.Title, &item.Description,
		&item.MediaURL, &item.MediaType, &item.ThumbnailURL,
		&item.Featured, &item.SortOrder, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPortfolioNotFound
		}
		return nil, fmt.Errorf("scan portfolio item: %w", err)
	}
	return item, nil
}

func init() {
	var _ PortfolioRepository = (*portfolioRepository)(nil)
}
