package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/events"
	"github.com/elevatecompact/spark/services/identity-service/internal/repository"
)

type UpdateUserProfile struct {
	Username    *string  `json:"username,omitempty"`
	DisplayName *string  `json:"display_name,omitempty"`
	Bio         *string  `json:"bio,omitempty"`
	AvatarURL   *string  `json:"avatar_url,omitempty"`
	BannerURL   *string  `json:"banner_url,omitempty"`
	Categories  *[]string `json:"categories,omitempty"`
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetPublicProfile(ctx context.Context, userID uuid.UUID) (*domain.PublicUser, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates UpdateUserProfile) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	SuspendUser(ctx context.Context, userID uuid.UUID, reason string) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role domain.UserRole) error
	SearchUsers(ctx context.Context, query string, role string, limit, offset int) ([]*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	eventPub events.EventProducer
}

func NewUserService(userRepo repository.UserRepository, eventPub events.EventProducer) UserService {
	return &userService{
		userRepo: userRepo,
		eventPub: eventPub,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *userService) GetPublicProfile(ctx context.Context, userID uuid.UUID) (*domain.PublicUser, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.ToPublic(), nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates UpdateUserProfile) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if updates.Username != nil {
		trimmed := strings.TrimSpace(*updates.Username)
		if trimmed != user.Username {
			if err := validateUsername(trimmed); err != nil {
				return nil, domain.NewDomainErrorMsg(domain.ErrValidation, err.Error(), 400)
			}
			existing, err := s.userRepo.GetByUsername(ctx, trimmed)
			if err != nil && err != domain.ErrUserNotFound {
				return nil, fmt.Errorf("failed to check username: %w", err)
			}
			if existing != nil {
				return nil, domain.ErrUsernameTaken
			}
			user.Username = trimmed
		}
	}
	if updates.DisplayName != nil {
		user.DisplayName = strings.TrimSpace(*updates.DisplayName)
	}
	if updates.Bio != nil {
		user.Bio = *updates.Bio
	}
	if updates.AvatarURL != nil {
		user.AvatarURL = *updates.AvatarURL
	}
	if updates.BannerURL != nil {
		user.BannerURL = *updates.BannerURL
	}
	if updates.Categories != nil {
		user.Categories = *updates.Categories
	}

	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if err := s.eventPub.PublishUserUpdated(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(newHash)
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.HardDelete(ctx, userID); err != nil {
		return err
	}

	if err := s.eventPub.PublishUserDeleted(ctx, userID.String()); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *userService) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Verified = true
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	if err := s.eventPub.PublishUserUpdated(ctx, user); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *userService) SuspendUser(ctx context.Context, userID uuid.UUID, reason string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = domain.StatusSuspended
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to suspend user: %w", err)
	}

	if err := s.eventPub.PublishUserUpdated(ctx, user); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *userService) UpdateRole(ctx context.Context, userID uuid.UUID, role domain.UserRole) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Role = role
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if err := s.eventPub.PublishUserUpdated(ctx, user); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *userService) SearchUsers(ctx context.Context, query string, role string, limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.Search(ctx, query, role, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}
