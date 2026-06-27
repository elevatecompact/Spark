# Kubernetes Orchestration

## Architecture

Titan runs on Amazon EKS with multiple clusters segmented by environment (dev, staging, prod-us, prod-eu). Each cluster runs Kubernetes 1.28+ with Calico for network policies and OIDC-based authentication via AWS IAM.

## Cluster Topology

- **Control Plane**: AWS-managed, multi-AZ, with private API server endpoint.
- **Node Groups**: Graviton-based (arm64) for general workloads, p4d.24xlarge for GPU-accelerated ML training jobs.
- **Namespace Strategy**: Each engine gets a dedicated namespace with resource quotas and network policies enforcing least-privilege communication.

## Workloads

All engines run as Kubernetes Deployments (stateless) or StatefulSets (stateful — e.g., Vault, Nexus). Sidecar patterns include Envoy for service mesh (mTLS, traffic splitting), Fluent Bit for log shipping, and Node Exporter for host-level metrics.

## Scheduling & Autoscaling

- **Cluster Autoscaler**: Scales node groups based on pending pods.
- **HPA / VPA**: Horizontal pod autoscaling on CPU/memory/custom metrics; vertical pod autoscaling for batch jobs.
- **Descheduler**: Evicts pods to balance resource utilization.

## Security

Pod Security Standards (restricted profile), network policies deny by default, secrets stored in AWS Secrets Manager synced via External Secrets Operator, and OPA Gatekeeper enforces policy (no privileged containers, required resource limits, allowed registries only).

## Day 2 Operations

Velero for backup and restore of persistent volumes and cluster state, Kyverno for mutation and validation webhooks, Prometheus Operator for monitoring stack lifecycle, and certified upgrades through end-to-end conformance testing in staging before production promotion.