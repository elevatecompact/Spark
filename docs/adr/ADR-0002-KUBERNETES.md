# ADR-0002: Kubernetes for Orchestration

## Status

Accepted

## Context

Spark runs dozens of microservices ranging from stateless HTTP APIs to stateful streaming processors and GPU-accelerated AI inference workloads. The platform must support multi-region deployment, automated scaling based on real-time metrics, zero-downtime rolling updates, and self-healing in the face of node or service failures. The evaluation compared managed container orchestration platforms including Amazon ECS, HashiCorp Nomad, Docker Swarm, and Kubernetes. ECS offered simplicity but limited portability across cloud providers. Nomad provided scheduling flexibility but a smaller ecosystem. Docker Swarm was easier to operate but lacked advanced orchestration features. Kubernetes provided the strongest ecosystem, multi-cloud portability, and community investment.

## Decision

Use Kubernetes (k8s) as the container orchestration platform across all environments. Deployments use Helm charts templated per environment with values overrides. Cluster provisioning uses Terraform with EKS (AWS), AKS (Azure), and GKE (GCP) modules for multi-cloud portability. Service meshing is handled by Istio for traffic management, observability, and mTLS. Node autoscaling uses Karpenter for intelligent bin-packing and spot instance utilization. GPU nodes are tainted and tolerated for AI inference workloads only.

## Consequences

### Positive
- Multi-cloud portability with consistent deployment semantics
- Rich ecosystem for service mesh, observability, and scaling
- Broad community support and industry-standard tooling
- Declarative infrastructure with self-healing capabilities

### Negative
- Significant operational complexity; requires dedicated SRE team
- Steep learning curve for developers unfamiliar with k8s concepts
- Resource overhead from control plane components and sidecar proxies
