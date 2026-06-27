# PostgreSQL

PostgreSQL is SPARK's primary relational database, serving as the system of record for all transactional data. It stores user accounts, content metadata, subscriptions, transactions, and all other core business entities.

## Configuration

Production instances run PostgreSQL 16 with the following tuning parameters: shared_buffers set to 25% of available RAM, effective_cache_size set to 75% of available RAM, work_mem set to 64MB per operation, and maintenance_work_mem set to 1GB for vacuum operations. max_connections is capped at 200 with PgBouncer handling connection pooling at the application layer.

## Extensions

We use the following PostgreSQL extensions: pg_stat_statements for query performance monitoring, pgcrypto for column-level encryption, uuid-ossp for UUID generation, hstore for semi-structured metadata, pg_trgm for trigram-based text search fallback, and postgis for geospatial queries in location-based features.

## Key Tables

Core tables include users (profiles, authentication credentials), content (videos, metadata, status), subscriptions (plans, billing cycles, status), transactions (payments, refunds, adjustments), and notifications (delivery records, preferences). Each table includes created_at, updated_at, and a soft-delete column.

## Monitoring

Query performance is monitored through pg_stat_statements with alerts for queries exceeding 100ms average execution time. Connection pool saturation, replication lag, and vacuum age are tracked through Prometheus exporters with Grafana dashboards.
