# Rollback Strategy

## Philosophy

Rollbacks are a first-class design concern in Titan, not an afterthought. Every deployment includes an automated rollback plan that is validated before the deployment begins.

## Automated Rollbacks

The CD pipeline triggers an instantaneous rollback when error rate exceeds baseline by more than 2% for 60 seconds post-deployment, p99 latency exceeds baseline by more than 30%, health check endpoints return non-200 on more than 10% of probes, or synthetic transaction tests fail. For blue-green deployments, the load balancer flips back to the previous environment. For canary deployments, the canary is scaled to zero.

## Manual Rollbacks

A rollback can be initiated at any time by an SRE or on-call engineer via a GitHub Actions workflow (`rollback.yaml` with a target version parameter), the ArgoCD UI (one-click sync to a previous revision), or the CLI (`argo rollback`).

## Database Rollbacks

Titan follows the expand-contract pattern for schema migrations: expand with new columns/tables (app remains backward-compatible), migrate data in the background, contract by removing old columns in a subsequent release. This ensures that an application rollback never requires a database rollback. In rare cases of destructive migration, point-in-time recovery from RDS snapshot is executed.

## Post-Rollback

After a rollback, an incident postmortem is automatically drafted, the failed version is quarantined in the image registry, a fix branch is created from the last known good commit, and the CI pipeline runs with the fix before re-attempting deployment.