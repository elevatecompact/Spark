# creator-service — Testing Guide

## Unit Tests
- Channel name uniqueness and slug generation
- Tier pricing validation (minimums, currency formatting)
- Verification document processing pipeline
- Payout eligibility and preference serialization

## Integration Tests
- Full onboarding: register → channel → verify → create tier → payout setup
- Subscription tier CRUD with pricing validation
- Verification workflow: submit → pending → approve/reject
- Metrics materialized view refresh correctness

## Load Tests
Simulate 1000 concurrent channel creations. Measure metrics view refresh with 100K creator records. k6 script in 	ests/load/creator-onboarding.js.

## Running
`ash
go test ./internal/... -v -cover          # Unit
go test ./tests/integration/... -tags=integration  # Integration
k6 run tests/load/creator-onboarding.js   # Load
`
"@ | Set-Content (Join-Path C:\Users\Dell\Downloads\SPARK\services\creator-service "TESTING.md") -Encoding UTF8

@"
# creator-service — Runbook

## Alerts
- VerificationQueueDepth > 100 — Manual review backlog, page Trust & Safety
- ChannelCreateErrorRate > 2% — Potential DB issue, check connection pool
- PayoutFailureRate > 1% — Wallet integration problem

## Procedures
### Manual Verification Override
`ash
curl -X POST http://creator:4002/v1/admin/verification/{id}/approve -H "Authorization: Bearer "
`

### Rebuild Metrics View
`ash
./creator refresh-metrics
`

### Force Payout
`ash
curl -X POST http://creator:4002/v1/admin/creators/{id}/payout -H "Authorization: Bearer "
`
## On-Call
Verify dashboard, check verification queue depth, review recent tier changes.
