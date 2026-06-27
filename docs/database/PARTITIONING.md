# Partitioning Strategy

SPARK uses database partitioning to manage large tables, improve query performance, and simplify data lifecycle management. Partitioning strategies vary by database system.

## PostgreSQL Partitioning

Large tables are partitioned using declarative partitioning. The content table is range-partitioned by created_at on monthly boundaries. The events table is list-partitioned by event_type, with separate partitions for high-volume event types. The notifications table is partitioned by user_id hash into 16 partitions to distribute write load.

Partition pruning is verified through EXPLAIN plans for all common query patterns. New partitions are created automatically by a scheduled job that runs weekly, creating partitions for the next three months. Old partitions are detached and archived to cold storage rather than deleted.

## ClickHouse Partitioning

ClickHouse tables are partitioned by toYYYYMM(timestamp) by default, providing efficient time-range pruning. Large tables use a composite partition key combining month and event_type for more granular control over data lifecycle.

## OpenSearch Sharding

OpenSearch indices are time-based with daily rollovers. Each index is sized to approximately 50GB with five primary shards. The hot-warm-cold architecture moves indices through tiers based on age, with older indices having fewer replicas and eventually being closed for read-only access.

## Partition Maintenance

A centralized partition manager service handles partition creation, detachment, archiving, and cleanup. It runs as a Kubernetes CronJob and logs all operations to the audit system.
