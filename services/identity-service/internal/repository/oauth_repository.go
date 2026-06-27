package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
)

type OAuthRepository interface {
	CreateClient(ctx context.Context, client *domain.OAuthClient) error
	GetClientByID(ctx context.Context, id string) (*domain.OAuthClient, error)
	CreateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode) error
	GetAuthorizationCode(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	DeleteAuthorizationCode(ctx context.Context, code string) error
	CreateAccessToken(ctx context.Context, token *domain.OAuthToken) error
	GetAccessToken(ctx context.Context, tokenHash string) (*domain.OAuthToken, error)
	GetAccessTokenByRefreshHash(ctx context.Context, refreshHash string) (*domain.OAuthToken, error)
	DeleteAccessToken(ctx context.Context, id uuid.UUID) error
}

type oauthRepository struct {
	pool *pgxpool.Pool
}

func NewOAuthRepository(pool *pgxpool.Pool) OAuthRepository {
	return &oauthRepository{pool: pool}
}

func (r *oauthRepository) CreateClient(ctx context.Context, client *domain.OAuthClient) error {
	query := `
		INSERT INTO oauth_clients (id, client_secret, redirect_uris, grant_types, scope, name, logo_url, trusted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		client.ID, client.ClientSecret, client.RedirectURIs, client.GrantTypes,
		client.Scope, client.Name, client.LogoURL, client.Trusted,
		client.CreatedAt, client.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create oauth client: %w", err)
	}
	return nil
}

func (r *oauthRepository) GetClientByID(ctx context.Context, id string) (*domain.OAuthClient, error) {
	query := `
		SELECT id, client_secret, redirect_uris, grant_types, scope, name, logo_url, trusted, created_at, updated_at
		FROM oauth_clients WHERE id = $1`

	client := &domain.OAuthClient{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&client.ID, &client.ClientSecret, &client.RedirectURIs, &client.GrantTypes,
		&client.Scope, &client.Name, &client.LogoURL, &client.Trusted,
		&client.CreatedAt, &client.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvalidClientID
		}
		return nil, fmt.Errorf("failed to get oauth client: %w", err)
	}
	return client, nil
}

func (r *oauthRepository) CreateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode) error {
	query := `
		INSERT INTO oauth_authorization_codes (code, client_id, user_id, redirect_uri, scope, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		code.Code, code.ClientID, code.UserID, code.RedirectURI,
		code.Scope, code.ExpiresAt, code.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create authorization code: %w", err)
	}
	return nil
}

func (r *oauthRepository) GetAuthorizationCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	query := `
		SELECT code, client_id, user_id, redirect_uri, scope, expires_at, created_at
		FROM oauth_authorization_codes WHERE code = $1`

	authCode := &domain.AuthorizationCode{}
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&authCode.Code, &authCode.ClientID, &authCode.UserID,
		&authCode.RedirectURI, &authCode.Scope, &authCode.ExpiresAt, &authCode.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get authorization code: %w", err)
	}

	if time.Now().After(authCode.ExpiresAt) {
		r.DeleteAuthorizationCode(ctx, code)
		return nil, domain.ErrExpiredToken
	}

	return authCode, nil
}

func (r *oauthRepository) DeleteAuthorizationCode(ctx context.Context, code string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM oauth_authorization_codes WHERE code = $1`, code)
	if err != nil {
		return fmt.Errorf("failed to delete authorization code: %w", err)
	}
	return nil
}

func (r *oauthRepository) CreateAccessToken(ctx context.Context, token *domain.OAuthToken) error {
	query := `
		INSERT INTO oauth_access_tokens (id, client_id, user_id, access_token_hash, refresh_token_hash, scope, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at`

	err := r.pool.QueryRow(ctx, query,
		token.ID, token.ClientID, token.UserID,
		token.AccessTokenHash, token.RefreshTokenHash,
		token.Scope, token.ExpiresAt, token.CreatedAt,
	).Scan(&token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create access token: %w", err)
	}
	return nil
}

func (r *oauthRepository) GetAccessToken(ctx context.Context, tokenHash string) (*domain.OAuthToken, error) {
	query := `
		SELECT id, client_id, user_id, access_token_hash, refresh_token_hash, scope, expires_at, created_at
		FROM oauth_access_tokens WHERE access_token_hash = $1`

	token := &domain.OAuthToken{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID, &token.ClientID, &token.UserID,
		&token.AccessTokenHash, &token.RefreshTokenHash,
		&token.Scope, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		r.DeleteAccessToken(ctx, token.ID)
		return nil, domain.ErrExpiredToken
	}

	return token, nil
}

func (r *oauthRepository) GetAccessTokenByRefreshHash(ctx context.Context, refreshHash string) (*domain.OAuthToken, error) {
	query := `
		SELECT id, client_id, user_id, access_token_hash, refresh_token_hash, scope, expires_at, created_at
		FROM oauth_access_tokens WHERE refresh_token_hash = $1`

	token := &domain.OAuthToken{}
	err := r.pool.QueryRow(ctx, query, refreshHash).Scan(
		&token.ID, &token.ClientID, &token.UserID,
		&token.AccessTokenHash, &token.RefreshTokenHash,
		&token.Scope, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("failed to get access token by refresh hash: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		r.DeleteAccessToken(ctx, token.ID)
		return nil, domain.ErrExpiredToken
	}

	return token, nil
}

func (r *oauthRepository) DeleteAccessToken(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM oauth_access_tokens WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}
	return nil
}
