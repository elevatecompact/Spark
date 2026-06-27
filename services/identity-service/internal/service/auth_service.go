package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/events"
	"github.com/elevatecompact/spark/services/identity-service/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, email, username, password, displayName string) (*domain.User, *domain.Session, error)
	Login(ctx context.Context, email, password, ip, userAgent string) (*domain.Session, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error
	ValidateToken(ctx context.Context, tokenString string) (*domain.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tokenSvc    TokenService
	eventPub    events.EventProducer
	sessionTTL  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tokenSvc TokenService,
	eventPub events.EventProducer,
	sessionTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenSvc:    tokenSvc,
		eventPub:    eventPub,
		sessionTTL:  sessionTTL,
	}
}

func (s *authService) Register(ctx context.Context, email, username, password, displayName string) (*domain.User, *domain.Session, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)
	displayName = strings.TrimSpace(displayName)

	if err := validateEmail(email); err != nil {
		return nil, nil, domain.NewDomainErrorMsg(domain.ErrValidation, err.Error(), 400)
	}
	if err := validateUsername(username); err != nil {
		return nil, nil, domain.NewDomainErrorMsg(domain.ErrValidation, err.Error(), 400)
	}
	if err := validatePassword(password); err != nil {
		return nil, nil, domain.NewDomainErrorMsg(domain.ErrValidation, err.Error(), 400)
	}
	if displayName == "" {
		displayName = username
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, nil, domain.ErrEmailTaken
	}

	existingUser, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return nil, nil, domain.ErrUsernameTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: string(passwordHash),
		Bio:          "",
		AvatarURL:    "",
		BannerURL:    "",
		Verified:     false,
		Role:         domain.RoleViewer,
		Status:       domain.StatusActive,
		Categories:   []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		TokenHash:    s.tokenSvc.HashToken(accessToken),
		RefreshHash:  s.tokenSvc.HashToken(refreshToken),
		ExpiresAt:    time.Now().UTC().Add(s.sessionTTL),
		CreatedAt:    now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := s.eventPub.PublishUserCreated(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return user, session, nil
}

func (s *authService) Login(ctx context.Context, email, password, ip, userAgent string) (*domain.Session, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Status == domain.StatusSuspended {
		return nil, domain.ErrUserSuspended
	}
	if user.Status == domain.StatusBanned {
		return nil, domain.ErrUserBanned
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		TokenHash:    s.tokenSvc.HashToken(accessToken),
		RefreshHash:  s.tokenSvc.HashToken(refreshToken),
		IPAddress:    ip,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().UTC().Add(s.sessionTTL),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := s.eventPub.PublishUserLoggedIn(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return session, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	refreshHash := s.tokenSvc.HashToken(refreshToken)

	session, err := s.sessionRepo.GetByToken(ctx, refreshHash)
	if err != nil {
		if err == domain.ErrSessionExpired {
			return nil, domain.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		s.sessionRepo.Delete(ctx, session.ID)
		return nil, domain.ErrExpiredToken
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Status != domain.StatusActive {
		s.sessionRepo.DeleteByUserID(ctx, user.ID)
		return nil, domain.ErrForbidden
	}

	s.sessionRepo.Delete(ctx, session.ID)

	accessToken, err := s.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newSession := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		TokenHash:    s.tokenSvc.HashToken(accessToken),
		RefreshHash:  s.tokenSvc.HashToken(newRefreshToken),
		IPAddress:    session.IPAddress,
		UserAgent:    session.UserAgent,
		ExpiresAt:    time.Now().UTC().Add(s.sessionTTL),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return newSession, nil
}

func (s *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if err := s.eventPub.PublishUserLoggedOut(ctx, sessionID.String()); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *authService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	if err := s.eventPub.PublishUserLoggedOut(ctx, userID.String()); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := s.tokenSvc.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Status != domain.StatusActive {
		return nil, domain.ErrForbidden
	}

	return user, nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return domain.ErrPasswordTooWeak
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	return nil
}
