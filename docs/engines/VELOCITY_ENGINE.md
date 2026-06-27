# Velocity Engine — CDN Orchestration

## Purpose

Velocity is Titan's CDN orchestration engine. It manages content distribution across a multi-CDN strategy, optimizing for performance, cost, and reliability. Velocity handles origin selection, cache invalidation, and real-time traffic steering.

## Architecture

Velocity acts as a control plane over multiple CDN providers (CloudFront, Fastly, Cloudflare). It monitors origin health, CDN performance, and cost metrics to dynamically route traffic through the optimal path.

## Tech Stack

- **Language**: Go
- **CDN Providers**: AWS CloudFront, Fastly, Cloudflare
- **Origin Storage**: Nexus (S3-compatible)
- **DNS**: Route53 with latency-based routing and health checks
- **Cache**: Redis for real-time CDN performance data and routing decisions
- **Observability**: Prometheus metrics for per-CDN latency, error rate, cache hit ratio

## Key Features

- **Multi-CDN steering**: Dynamic traffic allocation based on real-time performance and cost
- **Origin shielding**: Shield tier reduces origin load by aggregating cache misses
- **Cache invalidation**: Bulk and selective cache purge with regional propagation tracking
- **Custom domain support**: Bring your own domain with automatic TLS certificate provisioning
- **Geo-restriction**: Region-based access control for licensing compliance
- **Token authentication**: Signed URLs and cookies for premium content access
- **Edge computation**: CloudFront Functions for lightweight request/response transformations at the edge
- **Failover**: Automatic CDN failover with < 30 second detection and switch

## Performance Targets

| Metric | Target |
|--------|--------|
| Global p99 latency (first byte) | < 50ms |
| Cache hit ratio | > 90% for popular content |
| CDN failover time | < 30 seconds |
| Invalidation propagation (global) | < 60 seconds |
| Uptime (combined multi-CDN) | 99.999% |