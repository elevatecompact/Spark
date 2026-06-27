# Vision Engine

**Purpose:** Real-time analytics and observability engine for the Titan platform.
**Tech Stack:** Go, ClickHouse, Kafka, Redis, Grafana, Prometheus, gRPC, Avro.

Vision ingests, processes, and stores time-series analytics data from all Titan engines. Provides dashboards, alerting, and ad-hoc querying with sub-second performance over petabytes of data. Built on ClickHouse for columnar storage and Kafka for ingestion.
