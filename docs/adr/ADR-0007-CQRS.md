# ADR-0007: CQRS Pattern

## Status

Accepted

## Context

Spark's data access patterns exhibit a fundamental asymmetry: write operations tend to be transactional, narrow-scoped, and consistency-sensitive, while read operations are varied, aggregation-heavy, and latency-sensitive. A single-model approach using PostgreSQL for both commands and queries led to contention on hot tables, complex queries that mixed transactional and analytical concerns, and difficulty optimizing for disparate access patterns. The analytics and recommendation systems required denormalized views that did not fit the normalized write model. The team evaluated CQRS with shared database, CQRS with separate databases, and full event sourcing. Full event sourcing introduced unacceptable complexity for the immediate use case. Separating command and query models without full event sourcing struck the right balance.

## Decision

Adopt the CQRS pattern with separate command and query models. The command side uses PostgreSQL as the write-optimized store with normalized schemas and strong consistency guarantees. The query side uses dedicated read-optimized stores: ClickHouse for analytical queries, OpenSearch for full-text search, and Redis for low-latency cached reads. Materialized views in PostgreSQL serve hybrid use cases. Synchronization from the command store to query stores occurs through Kafka events published by the write model, consumed by projection builders that maintain the read models. Read models are eventually consistent with a target latency under one second for most projections.

## Consequences

### Positive
- Optimized data models for each access pattern improve query performance
- Reduced contention on write tables by offloading read traffic
- Independent scaling of read and write workloads
- Natural alignment with event-driven architecture and Kafka event backbone

### Negative
- Eventual consistency adds complexity for use cases requiring immediate consistency
- Projection builders must handle duplicate and out-of-order events
- Additional infrastructure and operational cost for multiple data stores
- Team must understand multiple query paradigms and data models
