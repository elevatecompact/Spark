# Velocity Engine

**Purpose:** CDN orchestration and edge caching engine for media delivery.
**Tech Stack:** Go, Fastly API, CloudFront API, Akamai API, Redis, Envoy, gRPC, Prometheus.

Velocity manages content distribution across multiple CDN providers - cache warming, origin shield, purging, failover, and real-time traffic steering based on performance metrics. Provides unified API on top of diverse CDN backends for multi-CDN reliability and cost optimisation.
