# Canary Deployments

## Purpose

Canary deployments reduce deployment risk by exposing a new version to a small subset of users before full rollout. Titan uses canary releases for all production deployments after the initial smoke test phase.

## Canary Phases

| Phase | Traffic Share | Duration | Validation |
|-------|--------------|----------|------------|
| Phase 1 | 1% | 15 min | Error rate, p99 latency, 5xx count |
| Phase 2 | 5% | 15 min | + Business metrics (conversion, signups) |
| Phase 3 | 25% | 15 min | + Resource utilization (CPU, memory, network) |
| Phase 4 | 50% | 15 min | Full metric suite |
| Phase 5 | 100% | — | Rollout complete |

## Implementation

Kubernetes with Istio VirtualService using weighted destination rules between `stable` and `canary` subsets. Traffic splitting is based on HTTP headers (internal users always go to canary) and percentage-based for external traffic. Progressive delivery is automated via Flagger or Argo Rollouts.

## Metric Thresholds

A canary phase fails if HTTP 5xx rate exceeds baseline + 1%, p99 latency exceeds baseline + 20%, error budget burn rate exceeds 2, or CPU/memory usage exceeds baseline + 30%. Evaluation windows are 1-5 minutes per metric.

## Automated Actions

If a phase passes, the pipeline proceeds to the next phase automatically. If a phase fails, the canary is scaled to zero, traffic returns entirely to stable, and an incident is automatically created. Inconclusive results notify a human for manual evaluation. A Grafana dashboard shows real-time canary progress with annotations at each decision point.