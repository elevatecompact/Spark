# CD Pipeline

## Overview

The Titan CD pipeline takes a validated, signed image and delivers it to production through a graduated series of environments. Every promotion is gated by automated verification and, for production, human approval.

## Environments

| Environment | Purpose | Promotion Gate |
|-------------|---------|---------------|
| Build | Image is built, scanned, signed | CI green |
| Dev | Automated deployment on merge to `main` | None |
| Staging | Integration testing, canary prep | Smoke tests pass |
| Production | Live user traffic | Manual approval + canary |

## Pipeline Flow

### Dev
ArgoCD watches the GitOps repo for image tag updates. On merge to `main`, CI pushes new tag and opens a PR in the GitOps repo. Auto-merge after 2 minutes. ArgoCD syncs dev cluster.

### Staging
Release branch created from `main` with release notes. PR merges after engineering manager approval. ArgoCD syncs staging cluster. Full integration test suite runs post-deploy. Load generators simulate realistic traffic patterns.

### Production
Release engineer initiates deploy via GitHub Actions workflow_dispatch. System validates that staging passed, image is signed, and SBOM is current. A 1% canary is deployed for 15 minutes with metric comparison. If canary passes: gradual rollout (25%, 50%, 100%), each step lasting 10 minutes. If canary fails: automatic rollback and incident created.

## Rollback Automation

If error rate spikes > 1% or p99 latency increases > 20% compared to baseline at any stage, the pipeline immediately rolls back and notifies on-call via PagerDuty.