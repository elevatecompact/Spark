package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleViewer  UserRole = "viewer"
	RoleCreator UserRole = "creator"
	RoleAdmin   UserRole = "admin"
	RoleMod     UserRole = "moderator"
)

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusBanned    UserStatus = "banned"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"display_name"`
	PasswordHash string     `json:"-"`
	Bio          string     `json:"bio"`
	AvatarURL    string     `json:"avatar_url"`
	BannerURL    string     `json:"banner_url"`
	Verified     bool       `json:"verified"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	Categories   []string   `json:"categories"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PublicUser struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	BannerURL   string    `json:"banner_url"`
	Verified    bool      `json:"verified"`
	Role        UserRole  `json:"role"`
	Categories  []string  `json:"categories"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		BannerURL:   u.BannerURL,
		Verified:    u.Verified,
		Role:        u.Role,
		Categories:  u.Categories,
		CreatedAt:   u.CreatedAt,
	}
}
