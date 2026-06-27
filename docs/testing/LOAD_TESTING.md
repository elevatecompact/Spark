# Load Testing

## Strategy

Titan conducts load testing at every stage of the delivery pipeline to ensure services can handle peak traffic with adequate headroom. Load tests are automated in CI and include baseline comparison against the previous release.

## Tools

k6 is the primary load testing tool for HTTP and gRPC endpoints, with scripts written in JavaScript and stored alongside service code. Vegeta is used for constant-rate HTTP load testing and latency distribution analysis. Custom simulators handle WebSocket (Pulse engine) and real-time streaming scenarios.

## Load Profiles

| Profile | Traffic Pattern | Duration | Purpose |
|---------|----------------|----------|---------|
| Baseline | 50 req/s constant | 10 min | Measure steady-state performance |
| Spike | 0 → 1000 req/s in 10s | 2 min | Test auto-scaling reaction |
| Stress | Increment 10 req/s every 30s until failure | Variable | Find breaking point |
| Soak | 200 req/s constant | 4 hours | Detect memory leaks, slow degradation |

## Metrics Collected

Request latency at p50, p95, p99, and p99.9 percentiles. Error rate and error types including timeout, 5xx, and connection refused. Throughput in requests per second. Resource utilization on service pods including CPU, memory, goroutines, and open connections. Downstream dependency latency for database, cache, and queue.

## Thresholds

Load tests fail the CI pipeline if p99 latency exceeds 500ms for API endpoints, error rate exceeds 0.1%, throughput degrades by more than 10% compared to baseline, or any pod is OOMKilled or CPU throttled for more than 30 seconds. Results are published as a PR comment with comparison against the last known good run.