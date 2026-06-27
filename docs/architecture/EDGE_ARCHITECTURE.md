# Edge Computing Architecture

Spark's edge computing layer brings computation closer to end users to reduce latency, offload central infrastructure, and enable real-time processing at global scale. The edge architecture spans three tiers: Cloudflare Workers, regional edge nodes, and on-premise broadcaster edges.

## Edge Tiers

### Tier 1: Cloudflare Workers (Code at Edge)
Serverless functions execute at Cloudflare's 330+ points of presence. Workers handle request authentication, header manipulation, A/B testing, and lightweight API aggregation. They operate with sub-millisecond startup and a 10ms CPU execution limit.

### Tier 2: Regional Edge Nodes (Compute at Edge)
Bare-metal Kubernetes nodes deployed in strategic metro locations host latency-sensitive services including:
- WebRTC Selective Forwarding Units
- Stream transcoding initiation
- Real-time moderation filtering
- Local analytics aggregation

Each regional node runs a lightweight control plane that synchronizes with the central orchestrator.

### Tier 3: Broadcaster Edge (Ingest at Edge)
On-premise or ISP-located appliances provide low-latency ingest for professional broadcasters. These devices handle first-mile encoding, packet loss recovery, and local monitoring before forwarding to the regional tier.

## Edge Processing Use Cases

### Real-Time Moderation
Video frames are sampled at the edge and analyzed by lightweight ML models for content policy violations. Suspicious streams trigger alerts without sending full video to central servers.

### Localized Transcoding
Popular streams are transcoded at regional edges to reduce backbone bandwidth. Edge nodes cache transcoding configurations for fast spin-up.

### Predictive Prefetching
Edge nodes predict viewer behavior and preconnect to upstream sources, reducing join latency by 40% on average.

## Edge-to-Cloud Sync

Edge nodes maintain eventual consistency with central services via:
- **Local writes with async replication** for non-critical state
- **Read-through to origin** for authoritative data
- **Conflict-free replicated data types (CRDTs)** for collaborative moderation state

The edge architecture processes over 60% of all API requests without reaching origin servers.
