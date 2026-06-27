# ADR-0003: PostgreSQL as Primary Database

## Status

Accepted

## Context

Spark requires a primary relational database for storing user accounts, creator profiles, subscription state, content metadata, and transactional ledger data. The database must support ACID transactions, complex relational queries, row-level security for multi-tenant data isolation, and high availability with automated failover. Candidates included PostgreSQL, MySQL 8, Amazon Aurora, and CockroachDB. MySQL offered comparable features but weaker support for advanced data types and indexing. Amazon Aurora provided managed scaling but vendor lock-in. CockroachDB excelled at horizontal scaling but added latency for single-region workloads. PostgreSQL offered the best balance of feature completeness, mature replication, extension ecosystem, and operational flexibility.

## Decision

Use PostgreSQL 16 as the primary operational database. Deploy with Patroni for high-availability clustering and automated failover. Use pgBouncer for connection pooling. Replication uses streaming physical replication with synchronous commit for critical write paths and asynchronous replicas for read scaling. Extensions including PostGIS for geospatial queries, pg_partman for automated partitioning, and pg_stat_statements for query performance monitoring are standard. Schema migrations are managed by Sqitch with all changes reviewed and version-controlled.

## Consequences

### Positive
- Mature, battle-tested relational engine with strong ACID guarantees
- Rich extension ecosystem supporting geospatial, partitioning, and monitoring needs
- Row-level security enables multi-tenant data isolation at the database layer
- Strong replication semantics support high-availability requirements

### Negative
- Horizontal write scaling requires application-level sharding; not natively supported
- Connection overhead at scale demands careful pool management
- Vacuum and maintenance operations require ongoing DBA attention
