package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
)

type CreatorRepository interface {
	Create(ctx context.Context, creator *domain.Creator) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Creator, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Creator, error)
	Update(ctx context.Context, creator *domain.Creator) error
	Search(ctx context.Context, query string, categories, tags []string, language, country string, limit, offset int) ([]domain.Creator, int, error)
	GetTrending(ctx context.Context, limit, offset int) ([]domain.Creator, error)
	GetRecommended(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Creator, error)
	IncrementFollowers(ctx context.Context, id uuid.UUID, delta int) error
	IncrementViews(ctx context.Context, id uuid.UUID) error
	IncrementStreams(ctx context.Context, id uuid.UUID) error
	UpdateRank(ctx context.Context, id uuid.UUID, rank int) error
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Creator, int, error)
	GetFollowers(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]uuid.UUID, int, error)
	GetFollowing(ctx context.Context, followerID uuid.UUID, limit, offset int) ([]domain.Creator, int, error)
	AddFollower(ctx context.Context, followerID, creatorID uuid.UUID) error
	RemoveFollower(ctx context.Context, followerID, creatorID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, creatorID uuid.UUID) (bool, error)
}

type creatorRepository struct {
	pool *pgxpool.Pool
}

func NewCreatorRepository(pool *pgxpool.Pool) CreatorRepository {
	return &creatorRepository{pool: pool}
}

func (r *creatorRepository) Create(ctx context.Context, c *domain.Creator) error {
	categoriesJSON, err := json.Marshal(c.Categories)
	if err != nil {
		return fmt.Errorf("marshal categories: %w", err)
	}
	tagsJSON, err := json.Marshal(c.Tags)
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}
	socialJSON, err := json.Marshal(c.SocialLinks)
	if err != nil {
		return fmt.Errorf("marshal social links: %w", err)
	}

	query := INSERT INTO creators (id, user_id, display_name, bio, avatar_url, banner_url, categories, tags, language, country, timezone, social_links, verified, status, follower_count, subscriber_count, total_views, total_streams, level, rank, created_at, updated_at)
		VALUES (, , , , , , , , , , , , , , , , , , , , , )

	_, err = r.pool.Exec(ctx, query,
		c.ID, c.UserID, c.DisplayName, c.Bio, c.AvatarURL, c.BannerURL,
		categoriesJSON, tagsJSON, c.Language, c.Country, c.Timezone, socialJSON,
		c.Verified, c.Status, c.FollowerCount, c.SubscriberCount,
		c.TotalViews, c.TotalStreams, c.Level, c.Rank, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrCreatorAlreadyExists
		}
		return fmt.Errorf("insert creator: %w", err)
	}
	return nil
}

func (r *creatorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Creator, error) {
	return r.scanOne(ctx, r.pool.QueryRow(ctx, SELECT * FROM creators WHERE id = , id))
}

func (r *creatorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Creator, error) {
	return r.scanOne(ctx, r.pool.QueryRow(ctx, SELECT * FROM creators WHERE user_id = , userID))
}

func (r *creatorRepository) Update(ctx context.Context, c *domain.Creator) error {
	categoriesJSON, err := json.Marshal(c.Categories)
	if err != nil {
		return fmt.Errorf("marshal categories: %w", err)
	}
	tagsJSON, err := json.Marshal(c.Tags)
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}
	socialJSON, err := json.Marshal(c.SocialLinks)
	if err != nil {
		return fmt.Errorf("marshal social links: %w", err)
	}

	query := UPDATE creators SET display_name=, bio=, avatar_url=, banner_url=, categories=, tags=, language=, country=, timezone=, social_links=, verified=, verified_at=, status=, follower_count=, subscriber_count=, total_views=, total_streams=, level=, rank=, updated_at= WHERE id=

	_, err = r.pool.Exec(ctx, query,
		c.ID, c.DisplayName, c.Bio, c.AvatarURL, c.BannerURL,
		categoriesJSON, tagsJSON, c.Language, c.Country, c.Timezone, socialJSON,
		c.Verified, c.VerifiedAt, c.Status, c.FollowerCount, c.SubscriberCount,
		c.TotalViews, c.TotalStreams, c.Level, c.Rank, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update creator: %w", err)
	}
	return nil
}

func (r *creatorRepository) Search(ctx context.Context, query string, categories, tags []string, language, country string, limit, offset int) ([]domain.Creator, int, error) {
	args := []interface{}{}
	conditions := []string{"c.status = "}
	args = append(args, domain.CreatorActive)
	argIdx := 2

	if query != "" {
		conditions = append(conditions, fmt.Sprintf(	o_tsvector('english', c.display_name || ' ' || coalesce(c.bio, '')) @@ plainto_tsquery('english', $%d), argIdx))
		args = append(args, query)
		argIdx++
	}

	if len(categories) > 0 {
		catJSON, _ := json.Marshal(categories)
		conditions = append(conditions, fmt.Sprintf(c.categories @> $%d::jsonb, argIdx))
		args = append(args, string(catJSON))
		argIdx++
	}

	if language != "" {
		conditions = append(conditions, fmt.Sprintf(c.language = $%d, argIdx))
		args = append(args, language)
		argIdx++
	}

	if country != "" {
		conditions = append(conditions, fmt.Sprintf(c.country = $%d, argIdx))
		args = append(args, country)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(SELECT COUNT(*) FROM creators c WHERE %s, where)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count creators: %w", err)
	}

	dataQuery := fmt.Sprintf(SELECT c.* FROM creators c WHERE %s ORDER BY c.follower_count DESC LIMIT $%d OFFSET $%d, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("search creators: %w", err)
	}
	defer rows.Close()

	var creators []domain.Creator
	for rows.Next() {
		c, err := scanCreator(rows)
		if err != nil {
			return nil, 0, err
		}
		creators = append(creators, *c)
	}
	if creators == nil {
		creators = []domain.Creator{}
	}
	return creators, total, nil
}

func (r *creatorRepository) GetTrending(ctx context.Context, limit, offset int) ([]domain.Creator, error) {
	query := SELECT c.* FROM creators c WHERE c.status =  ORDER BY (c.follower_count * 2 + c.total_views / 100 + c.total_streams * 5) DESC, c.rank ASC LIMIT  OFFSET 
	rows, err := r.pool.Query(ctx, query, domain.CreatorActive, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get trending: %w", err)
	}
	defer rows.Close()

	var creators []domain.Creator
	for rows.Next() {
		c, err := scanCreator(rows)
		if err != nil {
			return nil, err
		}
		creators = append(creators, *c)
	}
	if creators == nil {
		creators = []domain.Creator{}
	}
	return creators, nil
}

func (r *creatorRepository) GetRecommended(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Creator, error) {
	query := SELECT c.* FROM creators c 
		WHERE c.status =  AND c.id NOT IN (SELECT creator_id FROM creator_followers WHERE follower_id = )
		ORDER BY c.follower_count DESC, c.total_views DESC LIMIT 
	rows, err := r.pool.Query(ctx, query, domain.CreatorActive, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recommended: %w", err)
	}
	defer rows.Close()

	var creators []domain.Creator
	for rows.Next() {
		c, err := scanCreator(rows)
		if err != nil {
			return nil, err
		}
		creators = append(creators, *c)
	}
	if creators == nil {
		creators = []domain.Creator{}
	}
	return creators, nil
}

func (r *creatorRepository) IncrementFollowers(ctx context.Context, id uuid.UUID, delta int) error {
	_, err := r.pool.Exec(ctx, UPDATE creators SET follower_count = follower_count + , updated_at =  WHERE id = , id, delta, time.Now())
	return err
}

func (r *creatorRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, UPDATE creators SET total_views = total_views + 1, updated_at =  WHERE id = , id, time.Now())
	return err
}

func (r *creatorRepository) IncrementStreams(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, UPDATE creators SET total_streams = total_streams + 1, updated_at =  WHERE id = , id, time.Now())
	return err
}

func (r *creatorRepository) UpdateRank(ctx context.Context, id uuid.UUID, rank int) error {
	_, err := r.pool.Exec(ctx, UPDATE creators SET rank = , updated_at =  WHERE id = , id, rank, time.Now())
	return err
}

func (r *creatorRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Creator, int, error) {
	countQuery := SELECT COUNT(*) FROM creators c WHERE c.status =  AND c.categories @> ::jsonb
	catID := categoryID.String()
	catJSON := fmt.Sprintf(["%s"], catID)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, domain.CreatorActive, catJSON).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count by category: %w", err)
	}

	query := SELECT c.* FROM creators c WHERE c.status =  AND c.categories @> ::jsonb ORDER BY c.follower_count DESC LIMIT  OFFSET 
	rows, err := r.pool.Query(ctx, query, domain.CreatorActive, catJSON, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get by category: %w", err)
	}
	defer rows.Close()

	var creators []domain.Creator
	for rows.Next() {
		c, err := scanCreator(rows)
		if err != nil {
			return nil, 0, err
		}
		creators = append(creators, *c)
	}
	if creators == nil {
		creators = []domain.Creator{}
	}
	return creators, total, nil
}

func (r *creatorRepository) GetFollowers(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]uuid.UUID, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, SELECT COUNT(*) FROM creator_followers WHERE creator_id = , creatorID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count followers: %w", err)
	}

	rows, err := r.pool.Query(ctx, SELECT follower_id FROM creator_followers WHERE creator_id =  ORDER BY followed_at DESC LIMIT  OFFSET , creatorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get followers: %w", err)
	}
	defer rows.Close()

	var followerIDs []uuid.UUID
	for rows.Next() {
		var fid uuid.UUID
		if err := rows.Scan(&fid); err != nil {
			return nil, 0, err
		}
		followerIDs = append(followerIDs, fid)
	}
	if followerIDs == nil {
		followerIDs = []uuid.UUID{}
	}
	return followerIDs, total, nil
}

func (r *creatorRepository) GetFollowing(ctx context.Context, followerID uuid.UUID, limit, offset int) ([]domain.Creator, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, SELECT COUNT(*) FROM creator_followers WHERE follower_id = , followerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count following: %w", err)
	}

	query := SELECT c.* FROM creators c INNER JOIN creator_followers cf ON cf.creator_id = c.id WHERE cf.follower_id =  ORDER BY cf.followed_at DESC LIMIT  OFFSET 
	rows, err := r.pool.Query(ctx, query, followerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get following: %w", err)
	}
	defer rows.Close()

	var creators []domain.Creator
	for rows.Next() {
		c, err := scanCreator(rows)
		if err != nil {
			return nil, 0, err
		}
		creators = append(creators, *c)
	}
	if creators == nil {
		creators = []domain.Creator{}
	}
	return creators, total, nil
}

func (r *creatorRepository) AddFollower(ctx context.Context, followerID, creatorID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, INSERT INTO creator_followers (follower_id, creator_id, followed_at) VALUES (, , ) ON CONFLICT DO NOTHING, followerID, creatorID, time.Now())
	return err
}

func (r *creatorRepository) RemoveFollower(ctx context.Context, followerID, creatorID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, DELETE FROM creator_followers WHERE follower_id =  AND creator_id = , followerID, creatorID)
	return err
}

func (r *creatorRepository) IsFollowing(ctx context.Context, followerID, creatorID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, SELECT EXISTS(SELECT 1 FROM creator_followers WHERE follower_id =  AND creator_id = ), followerID, creatorID).Scan(&exists)
	return exists, err
}

func (r *creatorRepository) scanOne(ctx context.Context, row pgx.Row) (*domain.Creator, error) {
	c, err := scanCreatorRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCreatorNotFound
		}
		return nil, err
	}
	return c, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanCreator(s scanner) (*domain.Creator, error) {
	return scanCreatorRow(s)
}

func scanCreatorRow(s scanner) (*domain.Creator, error) {
	var (
		categoriesJSON []byte
		tagsJSON       []byte
		socialJSON     []byte
		verifiedAt     *time.Time
	)
	c := &domain.Creator{}

	err := s.Scan(
		&c.ID, &c.UserID, &c.DisplayName, &c.Bio, &c.AvatarURL, &c.BannerURL,
		&categoriesJSON, &tagsJSON, &c.Language, &c.Country, &c.Timezone,
		&socialJSON, &c.Verified, &verifiedAt, &c.Status,
		&c.FollowerCount, &c.SubscriberCount, &c.TotalViews, &c.TotalStreams,
		&c.Level, &c.Rank, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan creator row: %w", err)
	}

	c.VerifiedAt = verifiedAt

	if len(categoriesJSON) > 0 {
		if err := json.Unmarshal(categoriesJSON, &c.Categories); err != nil {
			log.Warn().Err(err).Msg("unmarshal categories")
		}
	}
	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &c.Tags); err != nil {
			log.Warn().Err(err).Msg("unmarshal tags")
		}
	}
	if len(socialJSON) > 0 {
		if err := json.Unmarshal(socialJSON, &c.SocialLinks); err != nil {
			log.Warn().Err(err).Msg("unmarshal social links")
		}
	}
	if c.Categories == nil {
		c.Categories = []string{}
	}
	if c.Tags == nil {
		c.Tags = []string{}
	}

	return c, nil
}

func init() {
	var _ CreatorRepository = (*creatorRepository)(nil)
}
