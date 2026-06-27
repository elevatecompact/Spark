# analytics-service — Deployment Guide
Components: Stream processor (Kafka Streams, 6 replicas), Batch processor (Spark, 4 executors), API server (3 replicas), ClickHouse cluster (3 shards + 2 replicas).
K8s: k8s/analytics-service/ — api (3x 1GB), stream-processor (6x 2GB), spark-driver (1x).
Event ingestion: 50K events/s sustained. ClickHouse merges parts every 5min.
Deploy: kubectl apply -f k8s/analytics-service/, verify consumer lag.
