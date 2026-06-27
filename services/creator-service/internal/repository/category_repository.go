package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	List(ctx context.Context, activeOnly bool) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{pool: pool}
}

func (r *categoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := INSERT INTO categories (id, name, slug, description, icon_url, color, parent_id, sort_order, active, created_at)
		VALUES (, , , , , , , , , )
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.Name, c.Slug, c.Description, c.IconURL, c.Color,
		c.ParentID, c.SortOrder, c.Active, c.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("category slug already exists")
		}
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	row := r.pool.QueryRow(ctx, SELECT * FROM categories WHERE id = , id)
	c, err := scanCategory(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	row := r.pool.QueryRow(ctx, SELECT * FROM categories WHERE slug = , slug)
	c, err := scanCategory(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) List(ctx context.Context, activeOnly bool) ([]domain.Category, error) {
	var rows pgx.Rows
	var err error
	if activeOnly {
		rows, err = r.pool.Query(ctx, SELECT * FROM categories WHERE active = true ORDER BY sort_order ASC, name ASC)
	} else {
		rows, err = r.pool.Query(ctx, SELECT * FROM categories ORDER BY sort_order ASC, name ASC)
	}
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		c, err := scanCategory(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *c)
	}
	if categories == nil {
		categories = []domain.Category{}
	}
	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := UPDATE categories SET name=, slug=, description=, icon_url=, color=, parent_id=, sort_order=, active= WHERE id=
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.Name, c.Slug, c.Description, c.IconURL, c.Color,
		c.ParentID, c.SortOrder, c.Active,
	)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, DELETE FROM categories WHERE id = , id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

func scanCategory(s interface{ Scan(dest ...interface{}) error }) (*domain.Category, error) {
	c := &domain.Category{}
	err := s.Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.IconURL, &c.Color,
		&c.ParentID, &c.SortOrder, &c.Active, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan category: %w", err)
	}
	return c, nil
}

func init() {
	var _ CategoryRepository = (*categoryRepository)(nil)
}
