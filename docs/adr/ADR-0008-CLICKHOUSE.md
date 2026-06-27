# ADR-0008: ClickHouse for Analytics

## Status

Accepted

## Context

Spark generates massive volumes of time-series data including viewer engagement metrics, stream health telemetry, recommendation impression logs, payout calculations, and content performance analytics. These datasets grow at rates exceeding 10 TB per month. Query patterns involve high-cardinality aggregation over large time windows, often requiring sub-second response times for real-time dashboards. The evaluation compared ClickHouse, Druid, BigQuery, Snowflake, and TimescaleDB. Druid offered strong real-time ingestion but a complex deployment model. BigQuery and Snowflake provided fully managed analytics but introduced cloud-vendor dependency and growing costs at Spark's data scale. TimescaleDB extended PostgreSQL but could not match ClickHouse's columnar compression and vectorized execution for analytical workloads.

## Decision

Use ClickHouse as the primary analytical database. Deploy a multi-shard cluster with replication factor 2 across availability zones. Data is ingested via Kafka engine tables consuming from analytics-dedicated Kafka topics, with materialized views transforming raw events into aggregated tables. Tables use the ReplicatedMergeTree engine with partitioning by event date and ordering keys optimized for common query patterns. A separate Grafana datasource connects directly to ClickHouse for real-time dashboards. Long-term data retention follows a tiered strategy: hot data in ClickHouse local storage, warm data on object storage via the S3 engine, and cold data archived in Parquet format.

## Consequences

### Positive
- Columnar storage and vectorized execution deliver sub-second aggregations over billions of rows
- Exceptional compression ratios (5-10x) reduce storage costs
- Kafka engine tables enable real-time ingestion without separate pipelines
- SQL-compatible interface reduces team learning curve

### Negative
- Not suitable for point lookups or row-level OLTP workloads
- Mutations and deletes are expensive; data should be append-only
- Cluster topology changes require careful re-sharding procedures
- JOIN performance is limited compared to traditional OLTP databases
