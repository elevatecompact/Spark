# Backup Strategy

SPARK implements a comprehensive backup strategy to protect against data loss, corruption, and catastrophic failures. Backups are tested monthly through restoration drills.

## PostgreSQL Backups

Full database backups run daily using pg_dump with the custom format, enabling parallel restore and compression. WAL archival runs continuously through pg_receivewal, shipping WAL segments to object storage every 60 seconds. Point-in-time recovery allows restoration to any second within the retention window.

Backup retention is 30 days for daily full backups and 90 days for WAL archives. Monthly backups are retained for 12 months. All backups are encrypted at rest using AES-256 and encrypted in transit using TLS.

## ClickHouse Backups

ClickHouse data is backed up through clickhouse-backup, which creates snapshots of table data and metadata. Full backups run daily at midnight UTC. Incremental backups capture changed partitions every six hours. Backups are stored in object storage with 30-day retention.

## OpenSearch Backups

OpenSearch snapshots are taken daily using the snapshot API with an S3 repository plugin. Each snapshot captures all indices with their settings, mappings, and data. Snapshot retention is 14 days for daily snapshots and 60 days for weekly snapshots.

## Redis Backups

Redis AOF files are backed up hourly to object storage. RDB snapshots are taken every 6 hours for faster recovery scenarios. Backups are retained for 7 days.

## Restoration Testing

The first Saturday of every month, a full restoration drill is performed in the staging environment. Recovery time objectives are verified against the documented thresholds.
