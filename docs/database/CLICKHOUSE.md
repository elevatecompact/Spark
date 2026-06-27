# ClickHouse

ClickHouse is SPARK's columnar analytics database, purpose-built for high-performance OLAP queries on event data. It ingests millions of events per second and returns aggregate queries in milliseconds.

## Schema Design

Event data is organized in MergeTree tables with order by clauses optimized for common query patterns. The primary sorting key typically includes toDate(timestamp), event_type, and a shard key. Partitioning is by month for most tables, allowing efficient time-range pruning. Sampling is enabled on large tables to provide approximate results when sub-second response is required over sub-second precision.

## Data Ingestion

The Kafka engine table engine ingests events directly from the message bus. Materialized views transform raw events into pre-aggregated tables optimized for dashboard queries. Data flows through a buffer table to absorb ingestion spikes, then flushes to the main MergeTree table every 5 seconds or 100,000 rows.

## Query Optimization

AggregatingMergeTree tables maintain partial aggregation states for common metrics. SummingMergeTree tables roll up counts by minute for time-series dashboards. ReplicatedMergeTree provides high availability across three nodes. Distributed tables enable cluster-wide queries with automatic shard routing.

## Retention

Raw event data is retained for 90 days. Pre-aggregated rollups are retained for 24 months. Older data is moved to object storage through the ClickHouse S3 table function for infrequent access queries.
