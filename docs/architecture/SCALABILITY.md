# Scalability Strategy

Spark's scalability strategy addresses growth across multiple dimensions: concurrent viewers, stream count, geographic distribution, and feature complexity. The architecture scales horizontally at every layer with minimal operational overhead.

## Scaling Dimensions

### Horizontal Scaling (Services)
Stateless microservices scale horizontally via Kubernetes HPA based on CPU, memory, and custom metrics (requests per second, queue depth). Service meshes distribute traffic evenly across all replicas.

### Vertical Scaling (Data)
PostgreSQL scales up within a shard while shard count scales horizontally. Auto-sharding redistributes data when shard size exceeds thresholds. Read replicas absorb query traffic.

### Ingress Scaling
Cloudflare absorbs DDoS traffic and provides global anycast ingress. Envoy gateways scale with traffic volume. Each gateway handles 50,000 concurrent connections.

## Auto-Scaling Policies

| Component | Metric | Scale Out | Scale In |
|-----------|--------|-----------|----------|
| API services | Request latency > 100ms p99 | +3 replicas | -1 replica after 5 min idle |
| Transcoding | Queue depth > 1000 | +20% GPU nodes | -10% when queue < 100 |
| Kafka brokers | Partition load skew > 20% | Rebalance partitions | Repartition |
| Redis clusters | Memory > 70% | Add shard | Merge shards |
| OpenSearch | Query latency > 200ms | +2 data nodes | -1 after 10 min stable |

## Capacity Planning

### Current Capacity (per region)
- 5M concurrent viewers
- 50K concurrent streams
- 500K concurrent chat connections
- 10K transactions per second

### Growth Model
Spark targets 3x annual growth. Capacity planning uses predictive scaling based on historical patterns and scheduled events. Infrastructure provisioning is automated with 2-week lead time for hardware procurement.

## Database Scaling

### Read Scaling
Read replicas serve cache-miss queries and analytics workloads. Application-layer read-write splitting routes queries appropriately.

### Write Scaling
Write-intensive services use Kafka as a write buffer, decoupling ingestion from processing. Stream service writes queued events that batch-process to PostgreSQL.

### Storage Scaling
Nexus engine scales object storage by adding storage nodes. Erasure coding configuration adjusts for cluster size. Auto-tiering moves cold data to cheaper storage.
