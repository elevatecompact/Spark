# Integration Testing

## Purpose

Integration tests verify that Titan services interact correctly with their dependencies — databases, message queues, caches, object storage, and other engines. These tests validate real I/O paths without mocking.

## Approach

Every integration test spins up real dependency instances using Testcontainers (or equivalent). Tests run against ephemeral containers that are destroyed after the test suite completes.

## Dependency Containers

PostgreSQL 16 Alpine initialized with Flyway migrations, Redis 7 Alpine with default configuration, MinIO for S3-compatible storage emulation, RabbitMQ 3 Management for message queue testing, and LocalStack for AWS service emulation.

## Test Categories

### Database Tests
Schema migrations apply correctly, CRUD operations against real tables, transaction isolation and rollback behavior, and stored procedures and triggers.

### API Integration Tests
Request/response contract matches OpenAPI spec, authentication middleware (JWT, OAuth, API keys), rate limiting and throttling behavior, and payload validation error messages.

### Cross-Engine Tests
Oracle to Pulse recommendation flow, Forge to Nexus transcoding and storage pipeline, and Guardian to Vault authentication and payment flow.

## Execution

Integration tests run in parallel CI jobs, each with its own Docker Compose environment. The test suite should complete within 15 minutes. Any test taking longer than 30 seconds is flagged as a performance concern. Test data is seeded via programmatic factories with each test case getting a fresh, isolated data set. Teardown truncates all tables and removes object storage content.