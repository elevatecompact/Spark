# Multi-Layer Caching Strategy

Spark implements a multi-layer caching strategy to minimize latency, reduce origin load, and deliver content efficiently at global scale. Each cache layer targets different data types and operates with independent configurations.

## Cache Layers

### Layer 1: Browser Cache
The client-side cache stores static assets and API responses with appropriate Cache-Control headers. Video segments are cached in the browser's Media Source Extension buffer. Service workers cache API responses for offline-capable experiences.

`
Cache Duration: Assets (1 year), API (30s), Segments (stream duration)
Invalidation: Versioned URLs, service worker update
Hit Ratio: 40% for API, 90% for static assets
`

### Layer 2: CDN Edge Cache
Cloudflare edge caches deliver content from 330+ PoPs worldwide. Video segments have the highest priority with pre-warming on stream start. API responses use Cache-Tag headers for granular purging.

`
Cache Duration: Segments (24h), Manifests (5s), API (0-300s)
Invalidation: Tag-based purge on content change
Hit Ratio: 85% for video segments, 60% for API
`

### Layer 3: Velocity Engine Cache
The Velocity engine provides intelligent video segment caching with predictive preloading. It maintains a hot cache of segments likely to be requested next based on viewership patterns.

`
Cache Duration: Life of stream session
Invalidation: Session end or segment deprecation
Hit Ratio: 95% for active stream segments
`

### Layer 4: Application Cache (Redis)
Redis caches session data, user profiles, leaderboard rankings, and rate limiter counters. Redis clusters are deployed in each region with read replicas for scaling.

`
Cache Duration: Configurable per key type
Eviction Policy: LRU with per-key TTL
Hit Ratio: 90% for session data, 95% for leaderboards
`

## Cache Invalidation

A central cache invalidation service processes invalidation events from Kafka and purges the appropriate cache layers. The service uses Content-Addressable Storage keys for precise invalidation without false positives.

## Consistency Strategy

- **Strong consistency**: Wallet balances, stream state
- **Read-your-writes**: User profile updates
- **Eventual consistency**: Discovery, recommendations, trends

Each cache layer respects the consistency requirements of the data it serves.
