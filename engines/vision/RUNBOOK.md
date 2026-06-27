# Runbook

## Startup
1. Start ClickHouse: clickhouse-server.
2. Create tables: ./vision migrate.
3. Start Kafka consumers: ./vision-ingest.
4. Start query service: ./vision-query.
5. Verify: POST /v1/query -d '{"sql":"SELECT count() FROM events"}'.

## Monitoring
- Dashboard: ingestion rate, query latency, Kafka lag, storage.
- Alerts: Kafka lag > 5min, disk > 80%, query error rate > 5%.
