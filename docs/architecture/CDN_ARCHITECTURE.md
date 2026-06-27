# CDN Architecture

Spark delivers live and on-demand video content through a multi-CDN architecture with Cloudflare as the primary provider, augmented by the proprietary Velocity edge delivery engine for specialized streaming workloads.

## CDN Topology

`
Origin (S3/Nexus) → Cloudflare CDN → Velocity Engine → End User
       ↓                                      ↓
  Regional Cache                         Edge Caches
`

Cloudflare's global network provides the base delivery layer with 330+ points of presence. The Velocity engine sits on top, optimizing video-specific delivery with intelligent segment caching and predictive preloading.

## Cloudflare Integration

### Caching Strategy
- **HLS/CMAF segments**: Cache with TTL of 24 hours with immediate purge on stream end
- **Manifest files**: Short TTL of 5 seconds for live streams
- **API responses**: Cache-control headers set per endpoint with edge TTLs from 0 to 300 seconds
- **Static assets**: Immutable content with year-long TTL and content hash in URL

### Argo Smart Routing
Traffic between Cloudflare PoPs routes over optimized paths, reducing origin latency by 30%. Argo tunnels provide secure, persistent connections to origin servers.

### Workers for CDN Logic
Cloudflare Workers at the edge rewrite manifest URLs, insert ad markers, apply geo-blocking rules, and implement token-based authentication for premium content.

## Velocity Engine

The Velocity engine is Spark's proprietary delivery optimization layer:

### Predictive Caching
Velocity analyzes viewership patterns to predict which segments viewers will request next. Predicted segments are preloaded into edge cache before the player requests them, reducing segment fetch time by 60%.

### Multi-Protocol Delivery
The engine negotiates the optimal delivery protocol per client: HLS for Apple devices, CMAF for modern browsers, and WebRTC for low-latency applications. Protocol selection is transparent to the player.

### Real-Time Metrics
Velocity exports per-stream delivery metrics including time-to-first-frame, rebuffer ratio, average bitrate, and error rates. These feed into the streaming quality dashboard and automated alerting.

## Origin Shield

An origin shield layer prevents cache misses from overwhelming origin storage. Shield nodes aggregate requests from downstream caches and maintain a hot cache of popular content.
