# Guardian Engine

**Purpose:** Centralised authentication and authorisation engine for the Titan platform.
**Tech Stack:** Go, OAuth 2.0, OIDC, JWT, WebAuthn, Redis, PostgreSQL, mTLS.

Guardian handles user identity, session management, OAuth flows, WebAuthn passkeys, API key management, and fine-grained RBAC. Single auth boundary for all Titan engines providing token issuance, validation, and revocation with sub-millisecond verification.
