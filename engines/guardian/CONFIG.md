# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| jwt.signing_algorithm | Ed25519 | JWT signing algorithm |
| jwt.access_token_ttl | 3600 | Access token TTL (seconds) |
| jwt.refresh_token_ttl | 2592000 | Refresh token TTL (30 days) |
| oauth.providers | [] | OIDC provider configs |
| session.max_concurrent | 10 | Max sessions per user |
| rate_limit.login | 5/m | Login attempts per minute |
| webauthn.rp_name | Titan | WebAuthn relying party |
| postgres.conn_pool_size | 50 | Connection pool size |
