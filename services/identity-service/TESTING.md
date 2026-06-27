# identity-service — Testing Guide

## Unit Tests
- Password hashing with bcrypt cost factor validation
- JWT signing and verification (RS256 + HMAC)
- MFA TOTP code generation and time-step verification
- Email/password format validation and normalization
- OAuth state parameter generation and verification
- Token refresh rotation logic (no double-use)

## Integration Tests
- Full register → login → refresh → logout lifecycle
- MFA enrollment and second-factor challenge flow
- OAuth redirect handling with mock providers
- Session revocation and token blacklist enforcement
- Rate limiter correctness under concurrent requests
- Account suspension and reactivation workflow

## Running Tests
`ash
go test ./internal/... -v -cover -short          # Unit: ~10s
go test ./tests/integration/... -v -tags=integration  # Integration: ~60s
go test ./tests/e2e/... -v -tags=e2e             # E2E: ~5min
`

## Fixtures
Test fixtures in 	ests/fixtures/ provide pre-seeded users with known credentials. Use actory.NewUser() for randomized test data with controlled password hashes.
