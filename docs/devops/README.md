# DevOps — Titan Platform

## Overview

The DevOps discipline within Titan governs the entire infrastructure lifecycle — from provisioning bare-metal and cloud resources through deployment orchestration, observability, and incident response. Every operation is treated as code, every change is version-controlled, and every execution is auditable.

## Principles

- **Reproducibility**: All infrastructure is defined declaratively via Terraform and Kubernetes manifests. No manual mutations.
- **Observability by Default**: Every service emits metrics, logs, and traces. No deployment is considered healthy without verified telemetry.
- **Gradual Delivery**: Changes propagate through blue-green and canary strategies. Rollbacks are automated and tested.
- **Immutable Artifacts**: Every build produces a versioned, signed container image. Binaries are never patched in place.

## Core Stack

| Layer | Tooling |
|-------|---------|
| Infrastructure | Terraform, AWS CDK |
| Orchestration | Kubernetes (EKS) |
| Packaging | Docker, Helm |
| CI/CD | GitHub Actions, ArgoCD |
| Observability | Prometheus, Grafana, Loki, Tempo |
| Alerting | Alertmanager, PagerDuty |

## Repository Structure

```
infra/              Terraform modules and live configs
charts/             Helm charts for all services
k8s/                Raw Kubernetes manifests (overlays)
ci/                 Reusable CI workflows
docker/             Dockerfiles per service
runbooks/           Incident response procedures
```

Every engineer is expected to understand the delivery pipeline end-to-end and to contribute improvements to the DevOps tooling as part of normal development workflow.