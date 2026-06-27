# Observability Platform

## Unified Observability Model

Titan treats observability as a platform capability, not an afterthought. The observability stack is built on the Grafana LGTM (Loki, Grafana, Tempo, Mimir) ecosystem, providing a unified view across metrics, logs, traces, and profiles.

## Four Pillars

### Metrics (Prometheus + Mimir)
Every service exports standardized RED metrics (Rate, Errors, Duration). Custom business metrics track feature adoption, content relevance, and user engagement. Metrics retained for 28 days raw and 13 months aggregated.

### Logs (Loki)
Structured JSON logging with `log/slog` (Go) and tracing span IDs correlated. Labels include `service`, `environment`, `k8s_namespace`, `k8s_pod`, `engine`. Logs retained for 30 days for fast search and 1 year in cold storage.

### Traces (Tempo)
Distributed tracing across all engine boundaries via OpenTelemetry. Head-based sampling at 10% of all requests plus tail-based sampling at 100% of errors. Trace-to-logs and trace-to-metrics linking in Grafana.

### Profiles (Phlare/Pyroscope)
Continuous profiling for CPU, memory, and goroutine allocation. Integrated into CI as a budget: a PR that increases p99 CPU by more than 5% is flagged.

## Service Level Objectives

Every engine defines 3-5 SLOs with error budgets. SLO compliance is displayed on a shared dashboard and checked before production deployments.

## Instrumentation

All services use a shared OpenTelemetry SDK wrapper that enforces consistent span naming, attribute tagging, and context propagation. No manual instrumentation is required for basic RED metrics and trace context.