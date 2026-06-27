package domain

import (
	"time"

	"github.com/google/uuid"
)

type OAuthClient struct {
	ID           string    `json:"id"`
	ClientSecret string    `json:"client_secret,omitempty"`
	RedirectURIs []string  `json:"redirect_uris"`
	GrantTypes   []string  `json:"grant_types"`
	Scope        string    `json:"scope"`
	Name         string    `json:"name"`
	LogoURL      string    `json:"logo_url,omitempty"`
	Trusted      bool      `json:"trusted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthorizationCode struct {
	Code        string    `json:"code"`
	ClientID    string    `json:"client_id"`
	UserID      uuid.UUID `json:"user_id"`
	RedirectURI string    `json:"redirect_uri"`
	Scope       string    `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type OAuthToken struct {
	ID               uuid.UUID `json:"id"`
	ClientID         string    `json:"client_id"`
	UserID           uuid.UUID `json:"user_id"`
	AccessToken      string    `json:"access_token,omitempty"`
	AccessTokenHash  string    `json:"-"`
	RefreshToken     string    `json:"refresh_token,omitempty"`
	RefreshTokenHash string    `json:"-"`
	Scope            string    `json:"scope"`
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
}

type OIDCDiscovery struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	IntrospectionEndpoint            string   `json:"introspection_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	GrantTypesSupported              []string `json:"grant_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethods         []string `json:"token_endpoint_auth_methods"`
	ClaimsSupported                  []string `json:"claims_supported"`
}

type TokenIntrospection struct {
	Active   bool      `json:"active"`
	ClientID string    `json:"client_id,omitempty"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Scope    string    `json:"scope,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
