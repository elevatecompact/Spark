# Database Architecture

SPARK's database architecture follows a polyglot persistence model, where each database system is purpose-built for its workload. This approach avoids the one-size-fits-all compromise and allows each system to be independently scaled, tuned, and operated.

## Core Principles

**Workload Isolation.** Transactional workloads (OLTP) run on PostgreSQL. Analytical workloads (OLAP) run on ClickHouse. Search workloads run on OpenSearch. Graph traversals run on Neo4j. Caching and real-time messaging run on Redis. This isolation prevents resource contention and allows independent scaling.

**Data Domains.** Each microservice owns its data domain. Services communicate through APIs, not direct database access. This enforces bounded contexts and prevents tight coupling between services.

**Consistency Boundaries.** Strong consistency is maintained within a single database instance. Eventual consistency is accepted across database boundaries, with compensating transactions and idempotency keys used where necessary.

## Data Flow Architecture

User requests flow through the API gateway to microservices. Services read and write to their respective databases. Events are emitted to the message bus, which feeds the analytics pipeline (ClickHouse) and search index (OpenSearch). The graph database (Neo4j) is updated asynchronously from relational data for social network analysis.

## High Availability

All production databases run in multi-AZ configurations with automated failover. Read replicas provide horizontal read scaling. Connection pooling is managed through PgBouncer for PostgreSQL and custom connection managers for other systems.
