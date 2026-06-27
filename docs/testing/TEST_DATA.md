# Test Data Management

## Principles

Test data is treated as a managed artifact — versioned, reproducible, and isolated per test run. No test depends on data created by another test or a persistent fixture.

## Synthetic Data Generation

All unit and integration tests use programmatic data generation. Go uses `gofakeit` and `faker` for realistic user profiles, content metadata, and timestamps. Python uses `factory_boy` with `Faker` for Django and FastAPI model factories. TypeScript uses `@faker-js/faker` for frontend test fixtures. Factories produce valid domain objects including users, videos, recommendations, and payments with randomized but coherent fields.

## Anonymized Production Data

For performance and E2E tests, Titan uses anonymized snapshots of production data. A weekly cron job extracts a 1% sample of production data. PII fields including email, name, IP address, and user IDs are replaced with synthetic equivalents. The anonymized dataset is validated to ensure no PII leakage before loading into staging environments.

## Test Data Fixtures

For integration tests requiring specific states, version-controlled SQL and NoSQL seed files are stored in `testdata/seeds/`. Test fixtures are preferably created via the service's own API. Database snapshots are used only for read-heavy performance tests.

## Cleanup & Versioning

Every test environment is ephemeral in CI or reset nightly in staging. Integration tests use Testcontainers that are destroyed after the suite completes. E2E tests clean up created resources via API teardown hooks. Test data schemas are versioned alongside application schema with CI breaking if factory data is incompatible with the current schema.