package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/repository"
)

type OAuthService interface {
	Authorize(ctx context.Context, clientID, redirectURI, scope, state string, userID uuid.UUID) (*domain.AuthorizationCode, error)
	ExchangeAuthorizationCode(ctx context.Context, code, clientID, clientSecret, redirectURI string) (*domain.OAuthToken, error)
	ExchangePasswordCredentials(ctx context.Context, clientID, clientSecret, email, password, scope string) (*domain.OAuthToken, error)
	ExchangeRefreshToken(ctx context.Context, clientID, refreshToken string) (*domain.OAuthToken, error)
	ValidateAccessToken(ctx context.Context, token string) (*domain.User, error)
	IntrospectToken(ctx context.Context, token string) (*domain.TokenIntrospection, error)
	GetDiscoveryConfig(ctx context.Context, baseURL string) *domain.OIDCDiscovery
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type oauthService struct {
	oauthRepo  repository.OAuthRepository
	userRepo   repository.UserRepository
	tokenSvc   TokenService
	authSvc    AuthService
	cfg        OAuthConfig
}

type OAuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AuthCodeTTL     time.Duration
}

func NewOAuthService(
	oauthRepo repository.OAuthRepository,
	userRepo repository.UserRepository,
	tokenSvc TokenService,
	authSvc AuthService,
	cfg OAuthConfig,
) OAuthService {
	return &oauthService{
		oauthRepo: oauthRepo,
		userRepo:  userRepo,
		tokenSvc:  tokenSvc,
		authSvc:   authSvc,
		cfg:       cfg,
	}
}

func (s *oauthService) Authorize(ctx context.Context, clientID, redirectURI, scope, state string, userID uuid.UUID) (*domain.AuthorizationCode, error) {
	client, err := s.oauthRepo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !containsURI(client.RedirectURIs, redirectURI) {
		return nil, domain.ErrInvalidRedirectURI
	}

	if scope != "" && !isScopeValid(scope, client.Scope) {
		return nil, domain.ErrInvalidScope
	}

	codeBytes := make([]byte, 32)
	if _, err := rand.Read(codeBytes); err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}
	code := hex.EncodeToString(codeBytes)

	authCode := &domain.AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().UTC().Add(s.cfg.AuthCodeTTL),
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.oauthRepo.CreateAuthorizationCode(ctx, authCode); err != nil {
		return nil, err
	}

	return authCode, nil
}

func (s *oauthService) ExchangeAuthorizationCode(ctx context.Context, code, clientID, clientSecret, redirectURI string) (*domain.OAuthToken, error) {
	client, err := s.oauthRepo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if client.ClientSecret != clientSecret {
		return nil, domain.ErrInvalidClientSecret
	}

	authCode, err := s.oauthRepo.GetAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if authCode.ClientID != clientID {
		return nil, domain.ErrInvalidClientID
	}

	if authCode.RedirectURI != redirectURI {
		return nil, domain.ErrInvalidRedirectURI
	}

	if time.Now().UTC().After(authCode.ExpiresAt) {
		s.oauthRepo.DeleteAuthorizationCode(ctx, code)
		return nil, domain.ErrExpiredToken
	}

	s.oauthRepo.DeleteAuthorizationCode(ctx, code)

	now := time.Now().UTC()
	accessToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	oauthToken := &domain.OAuthToken{
		ID:               uuid.New(),
		ClientID:         clientID,
		UserID:           authCode.UserID,
		AccessToken:      accessToken,
		AccessTokenHash:  s.tokenSvc.HashToken(accessToken),
		RefreshToken:     refreshToken,
		RefreshTokenHash: s.tokenSvc.HashToken(refreshToken),
		Scope:            authCode.Scope,
		ExpiresAt:        now.Add(s.cfg.AccessTokenTTL),
		CreatedAt:        now,
	}

	if err := s.oauthRepo.CreateAccessToken(ctx, oauthToken); err != nil {
		return nil, err
	}

	return oauthToken, nil
}

func (s *oauthService) ExchangePasswordCredentials(ctx context.Context, clientID, clientSecret, email, password, scope string) (*domain.OAuthToken, error) {
	client, err := s.oauthRepo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if client.ClientSecret != clientSecret {
		return nil, domain.ErrInvalidClientSecret
	}

	session, err := s.authSvc.Login(ctx, email, password, "", "oauth")
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	accessToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	oauthToken := &domain.OAuthToken{
		ID:               uuid.New(),
		ClientID:         clientID,
		UserID:           session.UserID,
		AccessToken:      accessToken,
		AccessTokenHash:  s.tokenSvc.HashToken(accessToken),
		RefreshToken:     refreshToken,
		RefreshTokenHash: s.tokenSvc.HashToken(refreshToken),
		Scope:            scope,
		ExpiresAt:        now.Add(s.cfg.AccessTokenTTL),
		CreatedAt:        now,
	}

	if err := s.oauthRepo.CreateAccessToken(ctx, oauthToken); err != nil {
		return nil, err
	}

	return oauthToken, nil
}

func (s *oauthService) ExchangeRefreshToken(ctx context.Context, clientID, refreshToken string) (*domain.OAuthToken, error) {
	refreshHash := s.tokenSvc.HashToken(refreshToken)

	existingToken, err := s.oauthRepo.GetAccessTokenByRefreshHash(ctx, refreshHash)
	if err != nil {
		return nil, err
	}

	if existingToken.ClientID != clientID {
		return nil, domain.ErrInvalidClientID
	}

	s.oauthRepo.DeleteAccessToken(ctx, existingToken.ID)

	now := time.Now().UTC()
	newAccessToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken, _, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newToken := &domain.OAuthToken{
		ID:               uuid.New(),
		ClientID:         clientID,
		UserID:           existingToken.UserID,
		AccessToken:      newAccessToken,
		AccessTokenHash:  s.tokenSvc.HashToken(newAccessToken),
		RefreshToken:     newRefreshToken,
		RefreshTokenHash: s.tokenSvc.HashToken(newRefreshToken),
		Scope:            existingToken.Scope,
		ExpiresAt:        now.Add(s.cfg.AccessTokenTTL),
		CreatedAt:        now,
	}

	if err := s.oauthRepo.CreateAccessToken(ctx, newToken); err != nil {
		return nil, err
	}

	return newToken, nil
}

func (s *oauthService) ValidateAccessToken(ctx context.Context, token string) (*domain.User, error) {
	tokenHash := hashTokenSHA256(token)

	oauthToken, err := s.oauthRepo.GetAccessToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, oauthToken.UserID)
	if err != nil {
		return nil, err
	}

	if user.Status != domain.StatusActive {
		return nil, domain.ErrForbidden
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *oauthService) IntrospectToken(ctx context.Context, token string) (*domain.TokenIntrospection, error) {
	tokenHash := hashTokenSHA256(token)

	oauthToken, err := s.oauthRepo.GetAccessToken(ctx, tokenHash)
	if err != nil {
		return &domain.TokenIntrospection{Active: false}, nil
	}

	return &domain.TokenIntrospection{
		Active:    !time.Now().UTC().After(oauthToken.ExpiresAt),
		ClientID:  oauthToken.ClientID,
		UserID:    oauthToken.UserID,
		Scope:     oauthToken.Scope,
		ExpiresAt: oauthToken.ExpiresAt,
	}, nil
}

func (s *oauthService) GetDiscoveryConfig(ctx context.Context, baseURL string) *domain.OIDCDiscovery {
	return &domain.OIDCDiscovery{
		Issuer:                     baseURL,
		AuthorizationEndpoint:      baseURL + "/oauth/authorize",
		TokenEndpoint:              baseURL + "/oauth/token",
		UserinfoEndpoint:           baseURL + "/oauth/userinfo",
		IntrospectionEndpoint:      baseURL + "/oauth/introspect",
		JwksURI:                    baseURL + "/oauth/.well-known/jwks.json",
		ScopesSupported:            []string{"openid", "profile", "email", "offline_access"},
		ResponseTypesSupported:     []string{"code", "token", "id_token"},
		GrantTypesSupported:        []string{"authorization_code", "password", "client_credentials", "refresh_token"},
		SubjectTypesSupported:      []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256", "HS256"},
		TokenEndpointAuthMethods:   []string{"client_secret_basic", "client_secret_post"},
		ClaimsSupported:            []string{"sub", "name", "preferred_username", "email", "email_verified", "picture"},
	}
}

func (s *oauthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

func containsURI(uris []string, target string) bool {
	for _, uri := range uris {
		if uri == target {
			return true
		}
	}
	return false
}

func isScopeValid(requested, allowed string) bool {
	if allowed == "" {
		return true
	}
	requestedScopes := strings.Fields(requested)
	allowedScopes := strings.Fields(allowed)
	allowedSet := make(map[string]bool, len(allowedScopes))
	for _, s := range allowedScopes {
		allowedSet[s] = true
	}
	for _, s := range requestedScopes {
		if !allowedSet[s] {
			return false
		}
	}
	return true
}

func hashTokenSHA256(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
