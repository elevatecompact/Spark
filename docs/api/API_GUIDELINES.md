# API Design Guidelines

SPARK's APIs follow a consistent set of design principles to ensure predictability, discoverability, and ease of integration.

## Design Principles

**Consistency.** All APIs use the same authentication mechanism, error format, pagination pattern, and rate limiting headers regardless of the protocol. This reduces the learning curve for developers integrating with SPARK.

**Backward Compatibility.** API changes must never break existing clients. Additive changes (new fields, new endpoints) are always safe. Breaking changes require a new version and a deprecation period of at least 90 days.

**Idempotency.** Mutable operations support idempotency keys to allow safe retries. The server deduplicates requests with the same idempotency key within a 24-hour window.

**Semantic Naming.** Resource names are plural nouns. Endpoints follow the pattern /v{version}/{resource}/{id}. Nested resources use the pattern /{parent}/{parentId}/{child}. Actions on resources use HTTP methods: GET for retrieval, POST for creation, PATCH for partial updates, DELETE for removal.

**Pagination.** All list endpoints return paginated results using cursor-based pagination by default, with offset pagination available for simple use cases.

**Documentation.** Every endpoint is documented in the OpenAPI specification (REST) or GraphQL schema (GraphQL). Documentation is generated from code annotations and updated on every deployment.

## Security

All API traffic must use TLS 1.2 or higher. Authentication is required unless the endpoint is explicitly marked as public. Input validation follows a whitelist approach — known good inputs are accepted, everything else is rejected.
