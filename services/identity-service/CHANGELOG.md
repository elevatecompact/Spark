# identity-service — Changelog

## [1.4.0] — 2026-06-15
### Added
- WebAuthn passkey support for passwordless authentication on supported devices
- Session management UI API (list active sessions per user, remote logout)
- Per-API-key rate limiting instead of per-user

### Changed
- Migrated session store from Redis to PostgreSQL with read replicas for durability
- JWT kid header now includes key version for rotation transparency

## [1.3.0] — 2026-05-01
### Added
- OAuth scopes for granular API key permissions
- Account recovery via trusted devices (email + device verification)

### Fixed
- Race condition in MFA enrollment during concurrent requests

## [1.2.0] — 2026-03-10
### Added
- Multi-factor authentication (TOTP) with backup codes
- Email verification flow with one-time links
- Admin user suspension API

## [1.1.0] — 2026-01-20
### Added
- Social login (Google, GitHub, Discord) with account linking
- Refresh token rotation on each use
- API key management with role-scoped permissions

## [1.0.0] — 2025-11-01
### Initial Release — Email/password auth, JWT tokens, RBAC with user/admin roles
