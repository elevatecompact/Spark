# High Availability Design

Spark's architecture is designed for 99.99% availability for critical services and 99.95% for standard services. High availability is achieved through redundancy, graceful degradation, and automated recovery at every layer.

## Availability Targets

| Service Tier | Target | Measured (Q2 2026) | Monitoring |
|-------------|--------|--------------------|------------|
| Streaming (ingest + delivery) | 99.99% | 99.992% | End-to-end synthetic checks |
| API (gateway + services) | 99.99% | 99.988% | Request success rate |
| Financial transactions | 99.999% | 99.995% | Transaction success rate |
| Discovery & search | 99.95% | 99.97% | Query success rate |
| Administrative interfaces | 99.9% | 99.95% | Dashboard availability |

## HA Patterns

### Stateless Services
All stateless services run with minimum 3 replicas spread across availability zones. Horizontal Pod Autoscaler maintains 150% of required capacity. Pod Disruption Budgets ensure at least 50% of replicas remain available during rolling updates.

### Stateful Services
PostgreSQL uses streaming replication with automatic failover via Patroni. Redis clusters use Sentinel for automatic master election. Kafka uses rack-aware partition assignment with min-insync-replicas set to 2.

### Gateway and Load Balancing
Cloudflare global load balancer provides DNS-level failover. Envoy gateways run in each zone with health-aware routing. Multiple gateways per zone tolerate instance failures.

## Graceful Degradation

When components fail, Spark degrades gracefully rather than failing entirely:

| Failure | Degraded Behavior | User Impact |
|---------|-------------------|-------------|
| Recommendation service offline | Serve cached popular content | Less personalized discovery |
| Transcoding pipeline degraded | Serve single bitrate only | Lower video quality |
| Chat service degraded | Read-only mode | Cannot send messages |
| Moderation service offline | Queue actions for review | Delayed enforcement |

## SLA Measurement

Availability is measured using SLI probes that run every 30 seconds from multiple geographic locations. SLO burn rate alerts trigger when error budget is consumed faster than expected. Monthly availability reports are published to stakeholders with detailed incident analysis.
