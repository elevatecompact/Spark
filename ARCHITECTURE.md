# Architecture Overview

Spark is built on a **microservices architecture** with event-driven communication, edge-optimized delivery, and cloud-agnostic infrastructure.

## High-Level Diagram

```
+-------------+   +--------------+   +--------------+
¦  Web (Nx.js) ¦   ¦ Mobile (Fltr)¦   ¦  3rd-Party   ¦
¦  + React/TS  ¦   ¦  + Dart      ¦   ¦  via SDK/API ¦
+--------------+   +--------------+   +--------------+
       ¦                  ¦                  ¦
       +------------------+------------------+
                          ¦
                   +------?------+
                   ¦   Envoy     ¦
                   ¦  API Gateway¦
                   +-------------+
                          ¦
          +---------------+---------------+
          ¦               ¦               ¦
   +------?------+ +------?------+ +------?------+
   ¦  Auth Svc   ¦ ¦ Content Svc ¦ ¦  Stream Svc  ¦
   ¦  (Go)       ¦ ¦  (Rust)     ¦ ¦  (Go)        ¦
   +-------------+ +-------------+ +-------------+
          ¦               ¦               ¦
   +------?---------------?---------------?------+
   ¦               Kafka (Event Bus)               ¦
   +---------------------------------------------+
          ¦               ¦               ¦
   +------?------+ +------?------+ +------?------+
   ¦ PostgreSQL  ¦ ¦  ClickHouse ¦ ¦  OpenSearch ¦
   ¦  (OLTP)     ¦ ¦  (Analytics)¦ ¦  (Search)   ¦
   +-------------+ +-------------+ +-------------+
```

## Key Design Decisions

- **API Gateway (Envoy)**: Routes all external traffic, handles TLS termination, rate limiting, and auth token validation.
- **Event Bus (Kafka)**: Decouples all services. Events are schema-registered with Avro.
- **CQRS + Event Sourcing**: Write operations go to PostgreSQL; reads are served from materialized views or ClickHouse.
- **Media Pipeline**: Upload ? transcode (FFmpeg/Rust) ? CDN (edge) ? playback. Fully asynchronous via Kafka.
- **Content Delivery**: Static assets served from CDN. Live video via WebRTC ingestion, relayed through selective forwarding units, and delivered via HLS at the edge.

See `docs/architecture/` for detailed ADRs and service specifications.
