# identity-service — Configuration

## Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| IDENTITY_PORT | 4001 | HTTP listener port |
| IDENTITY_DB_URL | — | PostgreSQL connection string (required) |
| IDENTITY_REDIS_URL | — | Redis connection string (required) |
| JWT_ACCESS_SECRET | — | RS256 private key path for access tokens |
| JWT_REFRESH_SECRET | — | HMAC secret for refresh tokens |
| JWT_ACCESS_TTL | 15m | Access token lifetime |
| JWT_REFRESH_TTL | 7d | Refresh token lifetime |
| MFA_ISSUER | Titan | TOTP issuer label for authenticator apps |
| OAUTH_GOOGLE_CLIENT_ID | — | Google OAuth 2.0 client ID |
| OAUTH_GITHUB_CLIENT_ID | — | GitHub OAuth app client ID |
| OAUTH_DISCORD_CLIENT_ID | — | Discord application client ID |
| BCRYPT_COST | 12 | Password hashing cost factor |
| MAX_LOGIN_ATTEMPTS | 5 | Before temporary lockout |

## Feature Flags
| Flag | Default | Purpose |
|------|---------|---------|
| mfa_required | false | Enforce MFA for all users |
| oauth_enabled | true | Enable social login flows |
| registration_open | true | Allow new account creation |
| rate_limiting_enabled | true | Apply rate limits to auth endpoints |

All secrets managed via HashiCorp Vault. Rotation policy: JWT secrets every 90 days.
