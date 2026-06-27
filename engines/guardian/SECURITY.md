# Security

- Password hashing: Argon2id with configurable cost.
- Token binding with cnf claim for holder-of-key JWTs.
- Brute force protection with adaptive rate limiting.
- Breached password detection via HIBP k-anonymity API.
- Session management with refresh token rotation.
- Audit logging to immutable ClickHouse store.
- WebAuthn attestation verification.
