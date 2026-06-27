# ADR-0010: Zero-Trust Security Model

## Status

Accepted

## Context

Spark's architecture spans multiple cloud providers, regions, Kubernetes clusters, and a complex mesh of internal services, APIs, and external integrations. The traditional perimeter-based security model, which trusts anything inside the corporate network, is insufficient for a distributed system where workloads run on shared infrastructure and access patterns cross trust boundaries. Threat actors targeting creator platforms increasingly exploit lateral movement after compromising a single service. Regulatory requirements including GDPR, CCPA, and SOC 2 demand strict access controls, audit logging, and data isolation. The evaluation compared zero-trust, perimeter-based with VPN, and beyondcorp-style models. Perimeter-based approaches were insufficient for multi-cloud. BeyondCorp provided inspiration but required Google-specific infrastructure.

## Decision

Adopt a zero-trust security model where no entity is trusted by default, regardless of network location. Every request is authenticated, authorized, and encrypted. Implementation includes: mutual TLS (mTLS) via Istio for all service-to-service communication; OAuth 2.0 with OIDC for user authentication; RBAC with attribute-based access control (ABAC) extensions for fine-grained authorization; SPIFFE/SPIRE for workload identity; network policies enforcing least-privilege at the pod level; and mandatory audit logging for all data access. Secrets are managed through HashiCorp Vault with dynamic database credentials and automatic rotation. All access decisions are logged and monitored through a centralized security information and event management pipeline.

## Consequences

### Positive
- Eliminates implicit trust based on network location, critical for multi-cloud deployments
- Workload identity through SPIFFE/SPIRE prevents credential theft and replay attacks
- Fine-grained authorization enables least-privilege access down to the API resource level
- Comprehensive audit trail satisfies regulatory compliance requirements

### Negative
- mTLS overhead adds latency and CPU consumption for every service call
- Istio sidecar footprint increases resource consumption per pod
- Certificate rotation and SPIRE agent management add operational complexity
- Teams must adopt security-first mindset; expanded blast radius of misconfigured policies
