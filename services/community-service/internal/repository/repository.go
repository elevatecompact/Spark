package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/community-service/internal/domain"
)

type CommunityRepository interface {
	Create(ctx context.Context, c *domain.Community) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Community, error)
	Update(ctx context.Context, c *domain.Community) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, category string, offset, limit int) ([]domain.Community, error)

	AddMember(ctx context.Context, m *domain.CommunityMember) error
	RemoveMember(ctx context.Context, communityID, userID uuid.UUID) error
	GetMember(ctx context.Context, communityID, userID uuid.UUID) (*domain.CommunityMember, error)
	ListMembers(ctx context.Context, communityID uuid.UUID, offset, limit int) ([]domain.CommunityMember, error)
	UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role domain.MemberRole) error

	CreatePost(ctx context.Context, p *domain.CommunityPost) error
	GetPost(ctx context.Context, id uuid.UUID) (*domain.CommunityPost, error)
	UpdatePost(ctx context.Context, p *domain.CommunityPost) error
	SoftDeletePost(ctx context.Context, id uuid.UUID) error
	ListPosts(ctx context.Context, communityID uuid.UUID, offset, limit int) ([]domain.CommunityPost, error)
	PinPost(ctx context.Context, id uuid.UUID, pinned bool) error

	CreateComment(ctx context.Context, c *domain.PostComment) error
	ListComments(ctx context.Context, postID uuid.UUID) ([]domain.PostComment, error)
	SoftDeleteComment(ctx context.Context, id uuid.UUID) error

	UpsertReaction(ctx context.Context, r *domain.PostReaction) error
	DeleteReaction(ctx context.Context, postID, userID uuid.UUID, emoji string) error

	GetAdminStats(ctx context.Context) (*domain.CommunityAdminStats, error)
}

type commRepo struct{ pool *pgxpool.Pool }

func NewCommunityRepository(pool *pgxpool.Pool) CommunityRepository { return &commRepo{pool} }

func (r *commRepo) Create(ctx context.Context, c *domain.Community) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO communities (id, name, description, creator_id, type, category, avatar_url, banner_url, rules, member_count, post_count, is_active, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,0,0,true,NOW())`,
		c.ID, c.Name, c.Description, c.CreatorID, c.Type, c.Category, c.AvatarURL, c.BannerURL, c.Rules)
	return err
}

func (r *commRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Community, error) {
	c := &domain.Community{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, creator_id, type, category, avatar_url, banner_url, rules, member_count, post_count, is_active, created_at FROM communities WHERE id=$1 AND is_active=true`, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatorID, &c.Type, &c.Category, &c.AvatarURL, &c.BannerURL, &c.Rules, &c.MemberCount, &c.PostCount, &c.IsActive, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return c, err
}

func (r *commRepo) Update(ctx context.Context, c *domain.Community) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE communities SET name=$2, description=$3, type=$4, category=$5, avatar_url=$6, banner_url=$7, rules=$8 WHERE id=$1`,
		c.ID, c.Name, c.Description, c.Type, c.Category, c.AvatarURL, c.BannerURL, c.Rules)
	return err
}

func (r *commRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE communities SET is_active=false WHERE id=$1`, id)
	return err
}

func (r *commRepo) List(ctx context.Context, category string, offset, limit int) ([]domain.Community, error) {
	if limit <= 0 || limit > 50 {
		limit = 25
	}
	var rows pgx.Rows
	var err error
	if category != "" {
		rows, err = r.pool.Query(ctx, `SELECT id, name, description, creator_id, type, category, avatar_url, banner_url, rules, member_count, post_count, is_active, created_at FROM communities WHERE is_active=true AND category=$1 ORDER BY member_count DESC OFFSET $2 LIMIT $3`, category, offset, limit)
	} else {
		rows, err = r.pool.Query(ctx, `SELECT id, name, description, creator_id, type, category, avatar_url, banner_url, rules, member_count, post_count, is_active, created_at FROM communities WHERE is_active=true ORDER BY member_count DESC OFFSET $1 LIMIT $2`, offset, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var communities []domain.Community
	for rows.Next() {
		var c domain.Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatorID, &c.Type, &c.Category, &c.AvatarURL, &c.BannerURL, &c.Rules, &c.MemberCount, &c.PostCount, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	if communities == nil {
		communities = []domain.Community{}
	}
	return communities, nil
}

func (r *commRepo) AddMember(ctx context.Context, m *domain.CommunityMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO community_members (community_id, user_id, role, joined_at, last_active_at) VALUES ($1,$2,$3,NOW(),NOW()) ON CONFLICT DO NOTHING`,
		m.CommunityID, m.UserID, m.Role)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE communities SET member_count = (SELECT COUNT(*) FROM community_members WHERE community_id=$1) WHERE id=$1`, m.CommunityID)
	return err
}

func (r *commRepo) RemoveMember(ctx context.Context, communityID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM community_members WHERE community_id=$1 AND user_id=$2`, communityID, userID)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE communities SET member_count = (SELECT COUNT(*) FROM community_members WHERE community_id=$1) WHERE id=$1`, communityID)
	return err
}

func (r *commRepo) GetMember(ctx context.Context, communityID, userID uuid.UUID) (*domain.CommunityMember, error) {
	m := &domain.CommunityMember{}
	err := r.pool.QueryRow(ctx,
		`SELECT community_id, user_id, role, joined_at, last_active_at FROM community_members WHERE community_id=$1 AND user_id=$2`,
		communityID, userID).Scan(&m.CommunityID, &m.UserID, &m.Role, &m.JoinedAt, &m.LastActiveAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return m, err
}

func (r *commRepo) ListMembers(ctx context.Context, communityID uuid.UUID, offset, limit int) ([]domain.CommunityMember, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx,
		`SELECT community_id, user_id, role, joined_at, last_active_at FROM community_members WHERE community_id=$1 ORDER BY joined_at ASC OFFSET $2 LIMIT $3`,
		communityID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []domain.CommunityMember
	for rows.Next() {
		var m domain.CommunityMember
		if err := rows.Scan(&m.CommunityID, &m.UserID, &m.Role, &m.JoinedAt, &m.LastActiveAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	if members == nil {
		members = []domain.CommunityMember{}
	}
	return members, nil
}

func (r *commRepo) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role domain.MemberRole) error {
	_, err := r.pool.Exec(ctx, `UPDATE community_members SET role=$3 WHERE community_id=$1 AND user_id=$2`, communityID, userID, role)
	return err
}

func (r *commRepo) CreatePost(ctx context.Context, p *domain.CommunityPost) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO community_posts (id, community_id, author_id, title, content, is_pinned, is_announcement, reaction_counts, comment_count, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,'{}'::jsonb,0,NOW())`,
		p.ID, p.CommunityID, p.AuthorID, p.Title, p.Content, p.IsPinned, p.IsAnnouncement)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE communities SET post_count = (SELECT COUNT(*) FROM community_posts WHERE community_id=$1 AND deleted_at IS NULL) WHERE id=$1`, p.CommunityID)
	return err
}

func (r *commRepo) GetPost(ctx context.Context, id uuid.UUID) (*domain.CommunityPost, error) {
	p := &domain.CommunityPost{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, community_id, author_id, title, content, is_pinned, is_announcement, reaction_counts, comment_count, deleted_at, created_at FROM community_posts WHERE id=$1 AND deleted_at IS NULL`, id).
		Scan(&p.ID, &p.CommunityID, &p.AuthorID, &p.Title, &p.Content, &p.IsPinned, &p.IsAnnouncement, &p.ReactionCounts, &p.CommentCount, &p.DeletedAt, &p.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *commRepo) UpdatePost(ctx context.Context, p *domain.CommunityPost) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE community_posts SET title=$2, content=$3, is_announcement=$4 WHERE id=$1 AND deleted_at IS NULL`,
		p.ID, p.Title, p.Content, p.IsAnnouncement)
	return err
}

func (r *commRepo) SoftDeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE community_posts SET deleted_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *commRepo) ListPosts(ctx context.Context, communityID uuid.UUID, offset, limit int) ([]domain.CommunityPost, error) {
	if limit <= 0 || limit > 50 {
		limit = 25
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id, community_id, author_id, title, content, is_pinned, is_announcement, reaction_counts, comment_count, deleted_at, created_at FROM community_posts WHERE community_id=$1 AND deleted_at IS NULL ORDER BY is_pinned DESC, created_at DESC OFFSET $2 LIMIT $3`,
		communityID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []domain.CommunityPost
	for rows.Next() {
		var p domain.CommunityPost
		if err := rows.Scan(&p.ID, &p.CommunityID, &p.AuthorID, &p.Title, &p.Content, &p.IsPinned, &p.IsAnnouncement, &p.ReactionCounts, &p.CommentCount, &p.DeletedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if posts == nil {
		posts = []domain.CommunityPost{}
	}
	return posts, nil
}

func (r *commRepo) PinPost(ctx context.Context, id uuid.UUID, pinned bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE community_posts SET is_pinned=$2 WHERE id=$1`, id, pinned)
	return err
}

func (r *commRepo) CreateComment(ctx context.Context, c *domain.PostComment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO post_comments (id, post_id, author_id, parent_id, content, reaction_counts, created_at) VALUES ($1,$2,$3,$4,$5,'{}'::jsonb,NOW())`,
		c.ID, c.PostID, c.AuthorID, c.ParentID, c.Content)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE community_posts SET comment_count = (SELECT COUNT(*) FROM post_comments WHERE post_id=$1 AND deleted_at IS NULL) WHERE id=$1`, c.PostID)
	return err
}

func (r *commRepo) ListComments(ctx context.Context, postID uuid.UUID) ([]domain.PostComment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, post_id, author_id, parent_id, content, reaction_counts, deleted_at, created_at FROM post_comments WHERE post_id=$1 AND deleted_at IS NULL ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []domain.PostComment
	for rows.Next() {
		var c domain.PostComment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.ParentID, &c.Content, &c.ReactionCounts, &c.DeletedAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []domain.PostComment{}
	}
	return comments, nil
}

func (r *commRepo) SoftDeleteComment(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE post_comments SET deleted_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *commRepo) UpsertReaction(ctx context.Context, re *domain.PostReaction) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO post_reactions (post_id, comment_id, user_id, emoji) VALUES ($1,$2,$3,$4) ON CONFLICT (post_id, comment_id, user_id, emoji) DO NOTHING`,
		re.PostID, re.CommentID, re.UserID, re.Emoji)
	return err
}

func (r *commRepo) DeleteReaction(ctx context.Context, postID, userID uuid.UUID, emoji string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM post_reactions WHERE post_id=$1 AND user_id=$2 AND emoji=$3 AND comment_id IS NULL`, postID, userID, emoji)
	return err
}

func (r *commRepo) GetAdminStats(ctx context.Context) (*domain.CommunityAdminStats, error) {
	s := &domain.CommunityAdminStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
		 (SELECT COUNT(*) FROM communities WHERE is_active=true) as total_communities,
		 (SELECT COALESCE(SUM(member_count), 0) FROM communities WHERE is_active=true) as total_members,
		 (SELECT COUNT(*) FROM community_posts WHERE deleted_at IS NULL) as total_posts,
		 (SELECT COUNT(*) FROM communities WHERE is_active=true AND member_count > 0) as active_communities`).
		Scan(&s.TotalCommunities, &s.TotalMembers, &s.TotalPosts, &s.ActiveCommunities)
	if err != nil {
		return nil, err
	}
	s.GrowthRate = 0.05
	return s, nil
}
