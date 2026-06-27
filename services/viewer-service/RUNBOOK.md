# viewer-service — Runbook
## Alerts: WatchEventLatency > 500ms (Kafka slow), HistoryWriteErrorRate > 1% (DB capacity), PreferenceCacheHitRate < 70%.
## Procedures: Clear viewer history via admin API, warm preference cache (./viewer cache warm --type preferences --count 10000), force history cleanup (./viewer history cleanup --before 2026-03-27).
## Monitor dead-letter queue for failed watch events. Check Kafka consumer lag for progress topic.
