# Loki Logging

## Architecture

Titan uses Grafana Loki for log aggregation. Logs are collected by Fluent Bit running as a DaemonSet on each Kubernetes node, shipped to Loki via the gRPC push endpoint, and stored in object storage (S3).

## Log Format

All services emit structured JSON logs via a standardized logging library. Every log line includes `timestamp` (RFC3339 with nanoseconds), `level` (debug, info, warn, error, fatal), `service`, `engine`, `environment`, `trace_id`, `span_id` (correlated with Tempo), `message` (human-readable), and `error` (structured error details).

## Labeling Strategy

Loki labels are kept to a small, cardinality-bound set: `service` (logical service name), `environment` (`dev`, `staging`, `production`), and `cluster` (cluster identifier). High-cardinality fields (pod name, request ID, user ID) remain as structured metadata.

## Querying

Common Loki LogQL queries include `{service="pulse"} |= "error"` for all Pulse engine errors, `{service="oracle"} | json | latency > 1.0` for slow Oracle requests, and `rate({environment="production"} |= "panic" [5m])` for panic rate across all services.

## Retention

Hot tier retains 7 days on S3 with SSD-backed caching. Warm tier retains 30 days on S3 Standard. Cold tier retains 365 days on S3 Glacier. Log access is scoped by environment with production logs restricted to SRE and on-call engineers via Grafana RBAC.