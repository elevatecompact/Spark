package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/discovery-service/internal/domain"
)

type DiscoveryRepository struct {
	pool *pgxpool.Pool
}

func NewDiscoveryRepository(pool *pgxpool.Pool) *DiscoveryRepository {
	return &DiscoveryRepository{pool: pool}
}

func (r *DiscoveryRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, slug, description, parent_id, icon_url, sort_order, is_active, content_count
		FROM categories WHERE is_active=true ORDER BY sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.IconURL, &c.SortOrder, &c.IsActive, &c.ContentCount); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, description, parent_id, icon_url, sort_order, is_active, content_count
		FROM categories WHERE slug=$1 AND is_active=true
	`, slug)
	c := &domain.Category{}
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.IconURL, &c.SortOrder, &c.IsActive, &c.ContentCount)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *DiscoveryRepository) GetSubcategories(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, slug, description, parent_id, icon_url, sort_order, is_active, content_count
		FROM categories WHERE parent_id=$1 AND is_active=true ORDER BY sort_order
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.IconURL, &c.SortOrder, &c.IsActive, &c.ContentCount); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetCategoryContentIDs(ctx context.Context, categorySlug string, limit, offset int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT content_id FROM category_contents cc
		JOIN categories c ON cc.category_id=c.id
		WHERE c.slug=$1 ORDER BY cc.added_at DESC LIMIT $2 OFFSET $3
	`, categorySlug, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *DiscoveryRepository) ListCollections(ctx context.Context, featured bool) ([]domain.Collection, error) {
	q := `SELECT id, title, description, type, cover_image_url, is_featured, start_at, end_at, curated_by, created_at FROM collections`
	args := []interface{}{}
	if featured {
		q += " WHERE is_featured=true"
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Collection
	for rows.Next() {
		var c domain.Collection
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Type, &c.CoverImageURL, &c.IsFeatured, &c.StartAt, &c.EndAt, &c.CuratedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetCollection(ctx context.Context, id uuid.UUID) (*domain.Collection, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, title, description, type, cover_image_url, is_featured, start_at, end_at, curated_by, created_at
		FROM collections WHERE id=$1
	`, id)
	c := &domain.Collection{}
	err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Type, &c.CoverImageURL, &c.IsFeatured, &c.StartAt, &c.EndAt, &c.CuratedBy, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	items, err := r.GetCollectionItems(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Items = items
	return c, nil
}

func (r *DiscoveryRepository) CreateCollection(ctx context.Context, c *domain.Collection) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO collections (id, title, description, type, cover_image_url, is_featured, start_at, end_at, curated_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, c.ID, c.Title, c.Description, string(c.Type), c.CoverImageURL, c.IsFeatured, c.StartAt, c.EndAt, c.CuratedBy, c.CreatedAt)
	return err
}

func (r *DiscoveryRepository) UpdateCollection(ctx context.Context, c *domain.Collection) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE collections SET title=$2, description=$3, type=$4, cover_image_url=$5, is_featured=$6, start_at=$7, end_at=$8, curated_by=$9
		WHERE id=$1
	`, c.ID, c.Title, c.Description, string(c.Type), c.CoverImageURL, c.IsFeatured, c.StartAt, c.EndAt, c.CuratedBy)
	return err
}

func (r *DiscoveryRepository) GetCollectionItems(ctx context.Context, collectionID uuid.UUID) ([]domain.CollectionItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT collection_id, content_id, sort_order, added_at
		FROM collection_items WHERE collection_id=$1 ORDER BY sort_order
	`, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.CollectionItem
	for rows.Next() {
		var ci domain.CollectionItem
		if err := rows.Scan(&ci.CollectionID, &ci.ContentID, &ci.SortOrder, &ci.AddedAt); err != nil {
			return nil, err
		}
		res = append(res, ci)
	}
	return res, nil
}

func (r *DiscoveryRepository) AddCollectionItem(ctx context.Context, ci *domain.CollectionItem) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO collection_items (collection_id, content_id, sort_order, added_at) VALUES ($1,$2,$3,$4)
	`, ci.CollectionID, ci.ContentID, ci.SortOrder, ci.AddedAt)
	return err
}

func (r *DiscoveryRepository) RemoveCollectionItem(ctx context.Context, collectionID, contentID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM collection_items WHERE collection_id=$1 AND content_id=$2`, collectionID, contentID)
	return err
}

func (r *DiscoveryRepository) GetEditorialPicks(ctx context.Context, pickType domain.PickType) ([]domain.EditorialPick, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT content_id, pick_type, label, reason, picked_by, start_at, end_at, sort_order
		FROM editorial_picks WHERE pick_type=$1 AND start_at<=NOW() AND end_at>=NOW() ORDER BY sort_order
	`, string(pickType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.EditorialPick
	for rows.Next() {
		var ep domain.EditorialPick
		if err := rows.Scan(&ep.ContentID, &ep.PickType, &ep.Label, &ep.Reason, &ep.PickedBy, &ep.StartAt, &ep.EndAt, &ep.SortOrder); err != nil {
			return nil, err
		}
		res = append(res, ep)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetTrendingContentIDs(ctx context.Context, limit int) ([]domain.TrendingItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT content_id, trending_score FROM trending_scores ORDER BY trending_score DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.TrendingItem
	for rows.Next() {
		var ti domain.TrendingItem
		if err := rows.Scan(&ti.ContentID, &ti.Score); err != nil {
			return nil, err
		}
		res = append(res, ti)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetTrendingContentIDsByCategory(ctx context.Context, categorySlug string, limit int) ([]domain.TrendingItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ts.content_id, ts.trending_score FROM trending_scores ts
		JOIN category_contents cc ON ts.content_id=cc.content_id
		JOIN categories c ON cc.category_id=c.id
		WHERE c.slug=$1 ORDER BY ts.trending_score DESC LIMIT $2
	`, categorySlug, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.TrendingItem
	for rows.Next() {
		var ti domain.TrendingItem
		if err := rows.Scan(&ti.ContentID, &ti.Score); err != nil {
			return nil, err
		}
		res = append(res, ti)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetTrendingCreatorIDs(ctx context.Context, limit int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT creator_id FROM trending_creators ORDER BY trending_score DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetHomeFeedContentIDs(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT content_id FROM home_feed ORDER BY score DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetNewContentIDs(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT content_id FROM new_contents ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *DiscoveryRepository) GetRelatedContentIDs(ctx context.Context, contentID uuid.UUID, limit int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT related_content_id FROM related_contents WHERE content_id=$1 ORDER BY score DESC LIMIT $2
	`, contentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}
