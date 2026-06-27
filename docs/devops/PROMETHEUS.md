# Prometheus Metrics

## Architecture

Titan runs Prometheus via the Prometheus Operator on each Kubernetes cluster. Metrics are scraped from all engine pods, Kubernetes components, and infrastructure (RDS, ElastiCache, ALB) via exporters.

## Scrape Configuration

Services advertise scrape endpoints via `prometheus.io/scrape: "true"` and `prometheus.io/port: "8080"` annotations. ServiceMonitor custom resources provide fine-grained scrape targets with relabeling configs. Thanos sidecar or Mimir provides a global view across clusters.

## Metric Categories

### RED Metrics (required for every service)
- `titan_requests_total` — labeled by `service`, `method`, `path`, `status_code`
- `titan_request_duration_seconds` — histogram buckets [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
- `titan_errors_total` — labeled by `service`, `error_type`

### Resource Metrics
Container CPU usage, memory working set bytes from kubelet/cAdvisor, plus custom metrics for queue depth, goroutines, and database connections.

### Business Metrics
Engine-specific metrics such as `titan_recommendations_served_total` (Oracle), `titan_streams_active` (Pulse), and `titan_videos_transcoded_total` (Forge).

## Alerting Rules

Alerting rules are defined per service and cluster. Examples include `TitanHighErrorRate` (error rate > 5% for 5m), `TitanLatencySpike` (p99 > 500ms for 5m), and `TitanPodCrashLooping` (restart count > 3 in 10m).

## Retention

Prometheus retains 15 days of local storage, Mimir retains 28 days with downsampling after 7 days, and Thanos bucket retains 13 months for long-term trend analysis.