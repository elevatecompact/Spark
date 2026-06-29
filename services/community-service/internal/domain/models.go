package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommunityType string

const (
	CommPublic     CommunityType = "public"
	CommRestricted CommunityType = "restricted"
	CommPrivate    CommunityType = "private"
)

type MemberRole string

const (
	RoleAdmin     MemberRole = "admin"
	RoleModerator MemberRole = "moderator"
	RoleMember    MemberRole = "member"
)

type Community struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatorID   uuid.UUID     `json:"creator_id"`
	Type        CommunityType `json:"type"`
	Category    string        `json:"category"`
	AvatarURL   string        `json:"avatar_url,omitempty"`
	BannerURL   string        `json:"banner_url,omitempty"`
	Rules       []string      `json:"rules,omitempty"`
	MemberCount int           `json:"member_count"`
	PostCount   int           `json:"post_count"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CommunityMember struct {
	CommunityID  uuid.UUID  `json:"community_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Role         MemberRole `json:"role"`
	JoinedAt     time.Time  `json:"joined_at"`
	LastActiveAt time.Time  `json:"last_active_at"`
}

type CommunityPost struct {
	ID              uuid.UUID              `json:"id"`
	CommunityID     uuid.UUID              `json:"community_id"`
	AuthorID        uuid.UUID              `json:"author_id"`
	Title           string                 `json:"title"`
	Content         string                 `json:"content"`
	IsPinned        bool                   `json:"is_pinned"`
	IsAnnouncement  bool                   `json:"is_announcement"`
	ReactionCounts  map[string]int         `json:"reaction_counts"`
	CommentCount    int                    `json:"comment_count"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

type PostComment struct {
	ID             uuid.UUID      `json:"id"`
	PostID         uuid.UUID      `json:"post_id"`
	AuthorID       uuid.UUID      `json:"author_id"`
	ParentID       *uuid.UUID     `json:"parent_id,omitempty"`
	Content        string         `json:"content"`
	ReactionCounts map[string]int `json:"reaction_counts"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

type PostReaction struct {
	PostID    uuid.UUID `json:"post_id"`
	CommentID *uuid.UUID `json:"comment_id,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	Emoji     string    `json:"emoji"`
}

type CommunityAdminStats struct {
	TotalCommunities  int     `json:"total_communities"`
	TotalMembers      int     `json:"total_members"`
	TotalPosts        int     `json:"total_posts"`
	ActiveCommunities int     `json:"active_communities"`
	GrowthRate        float64 `json:"growth_rate"`
}
