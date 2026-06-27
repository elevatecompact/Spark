# analytics-service — Runbook
## Alerts: EventProcessingLag > 60s, DashboardLoadLatency > 5s, ClickHouseQueryErrors > 1%, ReportGenerationFailure
## Flush stream: Restart stream processor to clear buffer: kubectl rollout restart deploy/analytics-stream
## Rebuild metric: POST /v1/admin/metrics/rebuild {metricName, timeRange}
## Check ClickHouse: SELECT * FROM system.mutations WHERE is_done=0
## Reprocess events: Republish failed events from dead-letter queue.
