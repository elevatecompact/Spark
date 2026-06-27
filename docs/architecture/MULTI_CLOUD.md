# Multi-Cloud Architecture

Spark adopts a multi-cloud strategy using AWS and Google Cloud Platform to eliminate single-provider dependency, optimize costs, and leverage each provider's unique capabilities. The architecture avoids provider lock-in through abstraction layers and portable infrastructure.

## Provider Distribution

### AWS (Primary)
AWS hosts the majority of production workloads including compute, managed databases, and the primary object storage tier. Services used: EKS for Kubernetes, RDS PostgreSQL, ElastiCache Redis, and MSK Kafka.

### GCP (Secondary)
GCP provides complementary infrastructure leveraging strengths in data analytics, ML inference, and network backbone. Services used: GKE for Kubernetes, BigQuery for analytics, Cloud GPUs for transcoding, and Cloud CDN as secondary CDN.

## Abstraction Layers

### Compute Abstraction
All services run in Kubernetes, providing a consistent deployment target across both clouds. A unified Istio service mesh spans providers via inter-cluster gateways. Workloads use provider-agnostic container images.

### Storage Abstraction
The Nexus engine abstracts underlying object storage, supporting both S3 and GCS backends. Data is replicated across providers for durability. A virtual storage layer handles provider-specific authentication and endpoint configuration.

### Messaging Abstraction
Kafka clusters run on each provider with MirrorMaker 2 bridging topics. Applications use a thin abstraction library that handles producer/consumer configuration differences.

## Data Residency

### Primary Storage
Each provider stores a full copy of critical data with configurable replication. User data is geo-restricted to its home region regardless of provider.

### Cross-Provider Replication
`
AWS S3 ← → Nexus Replication ← → GCP GCS
MSK Kafka ← → MirrorMaker 2 ← → GCP Pub/Sub (bridge)
`

## Failover Between Providers

Provider-level failover is reserved for catastrophic scenarios. The runbook covers:
1. DNS cutover to secondary provider
2. Data replication verification
3. Service mesh traffic migration
4. Gradual traffic ramp with canary validation

## Cost Optimization

Workloads are placed on the most cost-effective provider for their profile. Spot/preemptible instances handle batch transcoding. Reserved instances cover baseline capacity. A FinOps dashboard tracks provider spend with chargeback per service.
