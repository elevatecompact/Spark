# analytics-service — Database Schema
## ClickHouse — Time-series analytics
## events materialized table: event_date Date (partition), event_timestamp DateTime64(3), event_name String, user_id UUID, session_id UUID, properties JSON, context JSON
## metric_aggregates: metric_name String, time_bucket DateTime (1min granularity), dimension{}. Aggregate functions: count, sum, avg, p50, p95, p99
## PostgreSQL — Metadata store for dashboards, reports, funnel definitions, scheduled report configs, user preferences
