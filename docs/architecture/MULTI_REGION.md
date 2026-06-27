# Multi-Region Deployment Strategy

Spark operates across multiple cloud regions to provide low-latency access, regional compliance, and disaster resilience. The multi-region strategy ensures that users connect to the nearest healthy region with automatic failover.

## Region Topology

The platform is deployed in six active regions across three continents:

| Region | Provider | Purpose |
|--------|----------|---------|
| us-east-1 | AWS | Primary region, full services |
| us-west-2 | AWS | Secondary region, full services |
| eu-west-1 | AWS | European users, GDPR compliance |
| eu-central-1 | GCP | European redundancy |
| ap-southeast-1 | AWS | Asia-Pacific users |
| ap-northeast-1 | GCP | East Asia users |

## Deployment Model

### Active-Active
All regions serve production traffic simultaneously. Each region runs a full stack of services with its own data stores. Traffic is routed based on DNS-based geo-routing with health checks.

### Data Sovereignty
User data is stored and processed within regional boundaries. Cross-region data transfer is limited to metadata and analytics. Each region maintains its own Kafka cluster and database instances.

## Regional Services

### Stateful Services
- PostgreSQL: Each region has independent sharded clusters
- Redis: Regional clusters with local data only
- OpenSearch: Cross-cluster replication for search indices

### Stateless Services
- API gateway and microservices: Deployed identically in every region
- Transcoding workers: Regional GPU fleets
- Edge SFUs: Regional WebRTC processing nodes

## Traffic Routing

`
User → DNS (Latency-based) → Regional Gateway → Services
                    ↓
               Cloudflare (Global Load Balancer)
`

Global traffic management uses Cloudflare's load balancer with:
- Latency-based steering to the closest healthy region
- Geo-IP steering for compliance-required traffic
- Availability steering with circuit breaker thresholds

## Cross-Region Communication

Inter-region traffic is minimized. When necessary, it uses:
- Kafka MirrorMaker 2 for cross-region event replication (bidirectional)
- Private cloud interconnects for data synchronization
- Async replication with conflict resolution for user metadata

Health checks between regions detect cascading failures and trigger automated failover procedures.
