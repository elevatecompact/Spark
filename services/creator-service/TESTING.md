# creator-service — Testing Guide

## Test Categories

### Unit Tests
Located in internal/ packages with _test.go suffix. Cover:
- Core business logic and validation rules
- Input sanitization and error handling
- State machine transitions
- All external dependencies mocked

### Integration Tests
Located in 	ests/integration/. Cover:
- Full CRUD lifecycle for each resource
- Database migration and rollback
- Kafka producer/consumer behavior
- Redis caching behavior

### Contract Tests
Pact framework for consumer-driven contracts with downstream services.
Verified in CI pipeline on every PR.

## Running Tests
`ash
# Unit tests with coverage
go test ./internal/... -v -cover -short

# Integration tests (requires local dependencies)
go test ./tests/integration/... -v -tags=integration

# Full test suite
go test ./... -v -count=1 -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out -o coverage.html
`

## Test Fixtures
Use 	estify/suite for test organization. Factories in 	ests/factories/ generate randomized valid entities. WireMock stubs for external services in 	ests/mocks/.
