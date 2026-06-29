package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/community-service/internal/domain"
	"github.com/elevatecompact/spark/services/community-service/internal/events"
	"github.com/elevatecompact/spark/services/community-service/internal/repository"
)

type CommunityService interface {
	Create(ctx context.Context, c *domain.Community) (*domain.Community, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Community, error)
	Update(ctx context.Context, c *domain.Community) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, category string, page, size int) ([]domain.Community, error)

	Join(ctx context.Context, communityID, userID uuid.UUID) error
	Leave(ctx context.Context, communityID, userID uuid.UUID) error
	ListMembers(ctx context.Context, communityID uuid.UUID, page, size int) ([]domain.CommunityMember, error)
	UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role domain.MemberRole) error

	CreatePost(ctx context.Context, p *domain.CommunityPost) (*domain.CommunityPost, error)
	GetPost(ctx context.Context, id uuid.UUID) (*domain.CommunityPost, error)
	UpdatePost(ctx context.Context, p *domain.CommunityPost) error
	DeletePost(ctx context.Context, id uuid.UUID) error
	ListPosts(ctx context.Context, communityID uuid.UUID, page, size int) ([]domain.CommunityPost, error)
	PinPost(ctx context.Context, id uuid.UUID) error

	CreateComment(ctx context.Context, c *domain.PostComment) (*domain.PostComment, error)
	ListComments(ctx context.Context, postID uuid.UUID) ([]domain.PostComment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error

	ReactToPost(ctx context.Context, postID, userID uuid.UUID, emoji string) error
	ReactToComment(ctx context.Context, commentID, userID uuid.UUID, emoji string) error

	FeatureCommunity(ctx context.Context, id uuid.UUID) error
	SuspendCommunity(ctx context.Context, id uuid.UUID) error
	GetAdminStats(ctx context.Context) (*domain.CommunityAdminStats, error)
}

type commService struct {
	repo     repository.CommunityRepository
	eventPub events.EventProducer
}

func NewCommunityService(repo repository.CommunityRepository, eventPub events.EventProducer) CommunityService {
	return &commService{repo: repo, eventPub: eventPub}
}

func (s *commService) Create(ctx context.Context, c *domain.Community) (*domain.Community, error) {
	c.ID = uuid.New()
	c.IsActive = true
	c.MemberCount = 0
	c.PostCount = 0
	c.CreatedAt = time.Now().UTC()
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.repo.AddMember(ctx, &domain.CommunityMember{
		CommunityID: c.ID,
		UserID:      c.CreatorID,
		Role:        domain.RoleAdmin,
	})
	s.eventPub.PublishCommunityCreated(ctx, c.ID)
	return c, nil
}

func (s *commService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Community, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *commService) Update(ctx context.Context, c *domain.Community) error {
	return s.repo.Update(ctx, c)
}

func (s *commService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *commService) List(ctx context.Context, category string, page, size int) ([]domain.Community, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 25
	}
	offset := (page - 1) * size
	return s.repo.List(ctx, category, offset, size)
}

func (s *commService) Join(ctx context.Context, communityID, userID uuid.UUID) error {
	existing, err := s.repo.GetMember(ctx, communityID, userID)
	if err == nil && existing != nil {
		return domain.ErrAlreadyMember
	}
	return s.repo.AddMember(ctx, &domain.CommunityMember{
		CommunityID: communityID,
		UserID:      userID,
		Role:        domain.RoleMember,
	})
}

func (s *commService) Leave(ctx context.Context, communityID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, communityID, userID)
}

func (s *commService) ListMembers(ctx context.Context, communityID uuid.UUID, page, size int) ([]domain.CommunityMember, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 50
	}
	offset := (page - 1) * size
	return s.repo.ListMembers(ctx, communityID, offset, size)
}

func (s *commService) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role domain.MemberRole) error {
	if err := s.repo.UpdateMemberRole(ctx, communityID, userID, role); err != nil {
		return err
	}
	return s.eventPub.PublishRoleChanged(ctx, communityID, userID, string(role))
}

func (s *commService) CreatePost(ctx context.Context, p *domain.CommunityPost) (*domain.CommunityPost, error) {
	p.ID = uuid.New()
	p.CommentCount = 0
	p.ReactionCounts = map[string]int{}
	p.CreatedAt = time.Now().UTC()
	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	s.eventPub.PublishPostCreated(ctx, p.ID)
	return p, nil
}

func (s *commService) GetPost(ctx context.Context, id uuid.UUID) (*domain.CommunityPost, error) {
	return s.repo.GetPost(ctx, id)
}

func (s *commService) UpdatePost(ctx context.Context, p *domain.CommunityPost) error {
	return s.repo.UpdatePost(ctx, p)
}

func (s *commService) DeletePost(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDeletePost(ctx, id)
}

func (s *commService) ListPosts(ctx context.Context, communityID uuid.UUID, page, size int) ([]domain.CommunityPost, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 25
	}
	offset := (page - 1) * size
	return s.repo.ListPosts(ctx, communityID, offset, size)
}

func (s *commService) PinPost(ctx context.Context, id uuid.UUID) error {
	return s.repo.PinPost(ctx, id, true)
}

func (s *commService) CreateComment(ctx context.Context, c *domain.PostComment) (*domain.PostComment, error) {
	c.ID = uuid.New()
	c.ReactionCounts = map[string]int{}
	c.CreatedAt = time.Now().UTC()
	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *commService) ListComments(ctx context.Context, postID uuid.UUID) ([]domain.PostComment, error) {
	return s.repo.ListComments(ctx, postID)
}

func (s *commService) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDeleteComment(ctx, id)
}

func (s *commService) ReactToPost(ctx context.Context, postID, userID uuid.UUID, emoji string) error {
	return s.repo.UpsertReaction(ctx, &domain.PostReaction{
		PostID: postID,
		UserID: userID,
		Emoji:  emoji,
	})
}

func (s *commService) ReactToComment(ctx context.Context, commentID, userID uuid.UUID, emoji string) error {
	commentIDCopy := commentID
	return s.repo.UpsertReaction(ctx, &domain.PostReaction{
		CommentID: &commentIDCopy,
		UserID:    userID,
		Emoji:     emoji,
	})
}

func (s *commService) FeatureCommunity(ctx context.Context, id uuid.UUID) error {
	log.Info().Str("id", id.String()).Msg("community featured (noop)")
	return nil
}

func (s *commService) SuspendCommunity(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *commService) GetAdminStats(ctx context.Context) (*domain.CommunityAdminStats, error) {
	return s.repo.GetAdminStats(ctx)
}
