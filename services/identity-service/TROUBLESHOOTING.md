# identity-service — Troubleshooting

## Users Cannot Log In
**Symptoms:** 401 errors on /v1/auth/login, spike in auth_failure metric
**Root Causes:**
1. Password hash format mismatch — bcrypt cost factor changed without migration
2. JWT signing key rotated without overlapping grace period (tokens signed with old key rejected)
3. Redis session cache unavailable — refresh token validation fails
4. Account locked due to excessive failed attempts

**Diagnostic Steps:**
1. Verify DB connectivity: ./identity check-db
2. Check Redis ping: edis-cli ping
3. Compare JWT secret in Vault vs running pod: kubectl exec deploy/identity -- env | grep JWT
4. Check login attempt counters: edis-cli GET "lockout:{email}"

## MFA Codes Not Working
**Symptoms:** Valid TOTP codes rejected during login
**Causes:** Clock drift on identity pods (NTP not synchronized)
**Fix:** Verify NTP: w32tm /query /status. Restart time service. The TOTP algorithm tolerates ±1 time step (30s) by default.

## High Latency on Auth
**Causes:** DB connection pool exhausted, bcrypt cost too high (13+), rate limiter Redis calls slow. Increase DB_MAX_CONNS, verify bcrypt cost <= 12.
