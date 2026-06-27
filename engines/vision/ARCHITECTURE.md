# Architecture

Vision uses a log-centric architecture. Events flow into Kafka topics partitioned by event type. Go consumers transform and write data to ClickHouse using Avro-encoded batches. ClickHouse uses MergeTree engines with materialised views for pre-aggregated rollups. Query service (gRPC) serves dashboard and alert queries with rewriting to hit correct granularity table. Row-level security enforces tenant isolation. Prometheus metrics scraped via remote-write protocol.
