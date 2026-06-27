# API Gateway Architecture

Spark uses Envoy Proxy as the API gateway to manage ingress traffic, perform authentication, enforce rate limits, and route requests to backend services. Envoy is deployed as a standalone fleet alongside the Kubernetes service mesh.

## Gateway Topology

`
Client → Cloudflare CDN → Envoy Gateway → Backend Services
                   ↓
              Edge Cache (Velocity)
`

Two Envoy deployment tiers exist: edge gateways at each Cloudflare point of presence handle TLS termination and regional routing, while internal gateways within each Kubernetes cluster manage service-to-service routing.

## Key Capabilities

### TLS Termination
All external traffic terminates TLS at the edge. Envoy uses automated certificate management via Let's Encrypt with SPIFFE-compatible identity for mTLS between gateways and services.

### Authentication
Every request passes through an authentication filter that validates JWT tokens issued by the Identity Service. Tokens carry claims for user ID, roles, and session metadata. Invalid or expired tokens are rejected with 401 responses.

### Rate Limiting
Global and per-user rate limits are enforced at the gateway. Configuration is stored in Redis and supports sliding window counters. Different rate limits apply per endpoint group:

| Endpoint Group | Global (req/s) | Per User (req/s) |
|----------------|---------------|-------------------|
| Stream API | 10000 | 100 |
| Chat API | 50000 | 500 |
| Auth API | 2000 | 20 |
| Admin API | 500 | 50 |

### Routing
Requests are routed based on path prefixes, HTTP methods, and header values. gRPC-web traffic is transcoded automatically. Canary routing uses weight-based traffic splitting for gradual rollouts.

### Observability
Envoy exports detailed metrics to Prometheus and distributed trace spans to OpenTelemetry. Each request carries a unique trace ID propagated through the service mesh.
