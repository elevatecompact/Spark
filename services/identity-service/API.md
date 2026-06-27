# identity-service — API Contract

## Authentication
- POST /v1/auth/register — Create account (email, password, displayName)
  Request: {"email":"string","password":"string","displayName":"string"}
  Response: {"userId":"uuid","email":"string"}
- POST /v1/auth/login — Authenticate and return JWT pair
  Response: {"accessToken":"jwt","refreshToken":"jwt","expiresIn":900}
- POST /v1/auth/refresh — Rotate access token using refresh token
- POST /v1/auth/logout — Invalidate current session and revoke tokens
- POST /v1/auth/mfa/setup — Initialize TOTP enrollment, returns secret and QR code URL
- POST /v1/auth/mfa/verify — Verify TOTP code during login flow

## User Management
- GET /v1/users/me — Current user profile
- PATCH /v1/users/me — Update display name, avatar, preferences
- DELETE /v1/users/me — Schedule account deletion with 30-day grace period
- GET /v1/users/{id} — Public profile (admin only for private fields)

## API Keys
- POST /v1/api-keys — Generate new API key with scoped permissions
- GET /v1/api-keys — List active keys (last 4 chars only)
- DELETE /v1/api-keys/{id} — Revoke and rotate key

## Health
- GET /health — DB + Redis connectivity
- GET /ready — Migrations applied, caches warm
- GET /metrics — Prometheus endpoint on :4101

All responses use the standard envelope. Pagination via cursor for list endpoints.
