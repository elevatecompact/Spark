# analytics-service — Troubleshooting
## Dashboards returning stale data: Stream processor consumer lag, ClickHouse merge stalled, dashboard cache not invalidated. Check Kafka consumer lag, ClickHouse system.mutations, flush cache.
## Events not appearing: Kafka topic ACL issue, schema registry incompatibility, serializer mismatch. Check schema registry compatibility, verify event format.
## Report export fails: Too much data for single export, ClickHouse query timeout. Break into smaller time windows, increase ClickHouse max_execution_time.
