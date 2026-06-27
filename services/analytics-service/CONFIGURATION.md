# analytics-service — Configuration
ANALYTICS_PORT=4013, ANALYTICS_CLICKHOUSE_URL, ANALYTICS_POSTGRES_URL, ANALYTICS_REDIS_URL, ANALYTICS_KAFKA_BROKERS, EVENT_RETENTION_DAYS=90, AGGREGATION_INTERVAL_SECONDS=60, DASHBOARD_CACHE_TTL=30
FF: realtime_dashboards=true, historical_analytics=true, funnel_analysis=true, report_scheduling=true, anomaly_detection=false
Kafka consumer groups: analytics-events-processor with 12 partitions.
