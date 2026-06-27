# End-to-End Testing

## Scope

End-to-end tests validate complete user journeys across all Titan engines and touchpoints — web, mobile API, CDN delivery, and backend processing.

## Test Frameworks

Playwright is used for browser-based E2E tests covering the web UI. k6 handles API-based E2E scenarios for mobile and headless clients. A custom Go test harness covers backend-only workflows such as content ingestion through transcoding to delivery.

## Journey Catalog

User journeys include anonymous browse through sign up, content discovery, stream video playback, and like/comment interactions. Content creator journeys cover upload through transcode, publish, share, and analytics viewing. Subscription journeys test free tier upgrade to Pro, payment processing, and premium content access. Platform journeys validate multi-region failover, rate limit behavior, and concurrent content ingestion pipelines.

## Environment

E2E tests run against the staging environment, which mirrors production infrastructure and data scale. A dedicated test user pool with known credentials and content libraries is maintained.

## Test Design

Tests are data-independent, creating their own test data via API calls. They are idempotent with retry up to 2 times with exponential backoff for transient failures. Playwright captures screenshots on failure for post-mortem analysis.

## CI Integration

E2E tests run on every staging deployment and are gated — the deployment is considered green only if the full E2E suite passes. A smoke test subset runs on every production canary. The full suite completes within 30 minutes with individual scenarios targeting under 60 seconds.