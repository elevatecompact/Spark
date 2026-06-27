# Guardian Engine — Authentication

## Purpose

Guardian is Titan's authentication and authorization engine. It owns identity verification, session management, and access control across all Titan services. Guardian supports modern auth standards including OAuth 2.0, OIDC, WebAuthn passkeys, and SSO.

## Architecture

Guardian issues signed JWTs with short expiration. Sessions are validated via token introspection without database lookups. The engine integrates with external identity providers through OIDC discovery.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL (user identities, credentials, sessions)
- **Cache**: Redis (rate limiting, refresh token blacklist, OTP cache)
- **Keys**: AWS KMS for JWK signing key management with automatic rotation
- **Protocol**: OAuth 2.0, OIDC, WebAuthn (passkeys via FIDO2)

## Key Features

- **Multi-factor authentication**: TOTP, SMS OTP, email OTP, passkeys (WebAuthn)
- **OAuth 2.0 / OIDC provider**: First-party and third-party application authorization
- **Social login**: Google, Apple, GitHub, Twitter OIDC integration
- **Passkeys**: Passwordless authentication via WebAuthn (biometric, platform authenticator)
- **Session management**: View and revoke active sessions from any device
- **Role-based access control**: RBAC with hierarchical roles and permission sets
- **API key management**: Scoped API keys for programmatic access with per-key rate limits
- **Account recovery**: Graceful password reset flow with QR-code based device transfer

## Performance Targets

| Metric | Target |
|--------|--------|
| Token issuance latency | < 10ms (p99) |
| Token verification latency | < 2ms (p99) |
| Throughput per node | 20,000 auth requests/second |
| Uptime | 99.999% (auth is a critical path) |
| Key rotation | Every 90 days with 7-day grace period |