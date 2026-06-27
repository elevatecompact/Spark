# Tempo Tracing

## Trace Architecture

Titan uses Grafana Tempo for distributed tracing. All services are instrumented with OpenTelemetry SDKs that export traces via OTLP to the Tempo gateway. Tempo stores traces in S3 with a Parquet-backed backend for efficient querying.

## Instrumentation

HTTP/gRPC middleware, database clients, and message queue publishers all include OpenTelemetry instrumentation via shared libraries. Critical code paths add custom spans with domain-specific attributes such as `recommendation.candidate_count` and `stream.codec`. W3C TraceContext headers flow across engine boundaries automatically.

## Sampling

Head-based sampling covers 10% of all traffic for general observability. Tail-based sampling captures 100% of error spans. Probabilistic sampling at 1% is used for high-volume paths to reduce storage.

## Trace-to-Logs and Trace-to-Metrics

Every span carries its `trace_id` and `span_id`. Fluent Bit enriches log entries with these IDs, enabling direct navigation from trace waterfall to matching log lines. Prometheus exemplars attach trace IDs to metric samples for drill-down analysis.

## Search & Discovery

Tempo's TraceQL query language allows searching by service name, span name, duration range, resource attributes (`deployment.environment`, `service.version`), and span events with status codes.

## Retention & Performance

Traces are retained for 7 days in hot storage and 30 days in cold storage. High-value traces (errors, high latency) are tagged for extended retention of 90 days. The Tempo ingester handles 100k spans/second per cluster with 3x replication for durability.