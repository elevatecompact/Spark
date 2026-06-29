package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
)

type GiftItemRepository interface {
	Create(ctx context.Context, item *domain.GiftItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftItem, error)
	ListActive(ctx context.Context) ([]*domain.GiftItem, error)
	ListAll(ctx context.Context) ([]*domain.GiftItem, error)
	Update(ctx context.Context, item *domain.GiftItem) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type giftItemRepository struct {
	pool *pgxpool.Pool
}

func NewGiftItemRepository(pool *pgxpool.Pool) GiftItemRepository {
	return &giftItemRepository{pool: pool}
}

func (r *giftItemRepository) Create(ctx context.Context, item *domain.GiftItem) error {
	query := `INSERT INTO gift_items (id, name, price_cents, image_url, category, is_active, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, item.ID, item.Name, item.PriceCents, item.ImageURL, item.Category, item.IsActive, item.SortOrder, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create gift item: %w", err)
	}
	return nil
}

func (r *giftItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftItem, error) {
	query := `SELECT id, name, price_cents, image_url, category, is_active, sort_order, created_at, updated_at
		FROM gift_items WHERE id = $1`
	item := &domain.GiftItem{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.PriceCents, &item.ImageURL, &item.Category, &item.IsActive, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGiftItemNotFound
		}
		return nil, fmt.Errorf("failed to get gift item: %w", err)
	}
	return item, nil
}

func (r *giftItemRepository) ListActive(ctx context.Context) ([]*domain.GiftItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, price_cents, image_url, category, is_active, sort_order, created_at, updated_at
		FROM gift_items WHERE is_active = true ORDER BY sort_order ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list gift items: %w", err)
	}
	defer rows.Close()
	return scanGiftItems(rows)
}

func (r *giftItemRepository) ListAll(ctx context.Context) ([]*domain.GiftItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, price_cents, image_url, category, is_active, sort_order, created_at, updated_at
		FROM gift_items ORDER BY sort_order ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all gift items: %w", err)
	}
	defer rows.Close()
	return scanGiftItems(rows)
}

func (r *giftItemRepository) Update(ctx context.Context, item *domain.GiftItem) error {
	query := `UPDATE gift_items SET name=$2, price_cents=$3, image_url=$4, category=$5, is_active=$6, sort_order=$7, updated_at=NOW() WHERE id=$1`
	tag, err := r.pool.Exec(ctx, query, item.ID, item.Name, item.PriceCents, item.ImageURL, item.Category, item.IsActive, item.SortOrder)
	if err != nil {
		return fmt.Errorf("failed to update gift item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrGiftItemNotFound
	}
	return nil
}

func (r *giftItemRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `UPDATE gift_items SET is_active=false, updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete gift item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrGiftItemNotFound
	}
	return nil
}

func scanGiftItems(rows pgx.Rows) ([]*domain.GiftItem, error) {
	var items []*domain.GiftItem
	for rows.Next() {
		item := &domain.GiftItem{}
		if err := rows.Scan(&item.ID, &item.Name, &item.PriceCents, &item.ImageURL, &item.Category, &item.IsActive, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan gift item: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []*domain.GiftItem{}
	}
	return items, nil
}
