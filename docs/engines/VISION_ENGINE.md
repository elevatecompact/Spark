# Vision Engine — Real-Time Analytics

## Purpose

Vision is Titan's real-time analytics engine. It ingests, processes, and serves analytics data from every engine and user interaction on the platform. Vision powers dashboards, reports, and automated decision-making with sub-second query latency.

## Architecture

Vision uses a column-oriented storage architecture with ClickHouse at its core. Events are streamed from every engine via Kafka, transformed in a streaming pipeline, and materialized into ClickHouse tables optimized for analytical queries.

## Tech Stack

- **Language**: Go
- **Analytics Database**: ClickHouse (distributed, replicated)
- **Stream Processing**: Kafka + custom Go consumers for ETL
- **Cache**: Redis for real-time counters and leaderboards
- **Visualization**: Grafana for dashboards, embedded charts via Chart.js
- **Scheduling**: Cron-based rollup jobs for aggregate tables

## Key Features

- **Real-time dashboards**: Sub-second query latency on streaming data (viewer counts, revenue, engagement)
- **Event pipeline**: Standardized event schema ingested from all engines with schema validation
- **Pre-aggregated materialized views**: Minute, hourly, daily rollups for common queries
- **Retention management**: Automated data lifecycle (raw: 7 days, hourly: 90 days, daily: 5 years)
- **User-facing analytics**: Embedded analytics for content creators (views, earnings, audience)
- **Anomaly detection**: Statistical outlier detection on key metrics with automated alerting
- **Funnel analysis**: Conversion funnel queries across the user journey
- **SQL interface**: ClickHouse SQL for ad-hoc queries with row-level security

## Performance Targets

| Metric | Target |
|--------|--------|
| Query latency (p99) | < 100ms for pre-aggregated queries |
| Ingestion throughput | 1M events/second per cluster |
| Event to query visibility | < 2 seconds |
| Concurrent queries | 500+ |
| Storage compression ratio | 8:1 |