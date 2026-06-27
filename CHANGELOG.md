# Changelog

All notable changes to the Spark platform will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- New features to be listed here.

### Changed
- Behavioral changes to existing features listed here.

### Deprecated
- Features to be removed in future releases listed here.

### Removed
- Features removed in this release listed here.

### Fixed
- Bug fixes listed here.

### Security
- Security patches listed here.

---

## [0.1.0] — 2024-06-01

### Added
- Initial monorepo scaffold with pnpm workspaces
- Next.js web application with React 18 and TypeScript
- API gateway with JWT authentication and rate limiting
- User service (Go) for registration, profiles, and session management
- Video service (Rust) for transcoding and content storage
- PostgreSQL schema for users, content, and analytics
- Redis caching layer for session and rate limit data
- CI/CD pipeline with GitHub Actions
- Docker Compose setup for local development

### Changed
- N/A — initial release

[Unreleased]: https://github.com/spark/spark/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/spark/spark/releases/tag/v0.1.0
