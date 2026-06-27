# Redis

Redis serves as SPARK's caching layer, real-time messaging backbone, and session store. Its in-memory architecture provides sub-millisecond latency for critical hot-path operations.

## Deployment

Redis is deployed in a cluster topology with six nodes across three availability zones. Each node runs Redis 7.2 with transparent huge pages disabled and the kernel overcommit memory setting enabled. AOF persistence with fsync every second provides durability without significant performance impact. Sentinel monitors master health and triggers automatic failover within seconds.

## Use Cases

**Session Cache.** User sessions are stored with 24-hour TTL, refreshed on each authenticated request. Session invalidation propagates across all nodes within milliseconds.

**Content Cache.** Frequently accessed content metadata and rendered feed responses are cached with TTLs ranging from 60 seconds to 15 minutes depending on content type. Cache warming scripts pre-populate popular content on deployment.

**Rate Limiting.** Sliding window counters per user and per IP address enable the rate limiting system. Each increment is a single INCR command with EXPIRE, keeping the overhead minimal even at peak traffic.

**Pub/Sub.** Real-time notifications, streaming events, and live comment updates flow through Redis pub/sub channels. Channel naming follows the pattern {service}:{event_type}:{entity_id}.

## Memory Management

Maxmemory is set to 80% of instance RAM with an allkeys-lru eviction policy. Memory usage is monitored through the INFO memory command with alerts at 75% and 90% utilization thresholds.
