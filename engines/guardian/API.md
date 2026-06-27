# API

## Authentication
- POST /v1/auth/login - Authenticate with password.
- POST /v1/auth/refresh - Exchange refresh token.
- POST /v1/auth/logout - Revoke all tokens.
- POST /v1/auth/webauthn/register - Register WebAuthn credential.
- POST /v1/auth/webauthn/authenticate - Authenticate with WebAuthn.

## Token
- POST /v1/token/introspect - Validate token, return claims.
- GET /v1/token/public-key - JWKS for token verification.

## OAuth 2.0
- GET /v1/oauth/authorize, POST /v1/oauth/token, GET /v1/oauth/userinfo.

## Admin
- POST /v1/admin/users, POST /v1/admin/roles, GET /v1/admin/audit.
