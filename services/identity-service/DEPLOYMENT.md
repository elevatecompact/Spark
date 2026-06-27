# identity-service — Deployment Guide

## Prerequisites
- PostgreSQL 15+ with pgcrypto extension
- Redis 7+ cluster with persistence enabled
- HashiCorp Vault agent for JWT secret injection
- Kubernetes 1.28+ with cert-manager for TLS

## Kubernetes Manifests
Located in k8s/identity-service/:
- deployment.yaml — 3 replicas, 512Mi memory / 500m CPU limits
- service.yaml — ClusterIP on port 4001
- hpa.yaml — CPU target 70%, min 3 max 10 (auth is latency-sensitive)
- pdb.yaml — Min available 2 (critical path service)
- ault-annotation.yaml — Vault sidecar for secret injection

## Migration
`ash
./identity migrate up
./identity migrate status
`
Migrations create tables, indexes, and seed default admin role.

## Rollout Strategy
Canary via Istio: 10% traffic for 5 minutes, roll forward if error rate < 0.1%. Zero-downtime guaranteed via PDB. PreStop hook drains connections for 15s.

## Health
- /health — Validates DB, Redis, Vault connectivity
- /ready — All migrations applied, caches populated
- /metrics — Prometheus on port 4101
