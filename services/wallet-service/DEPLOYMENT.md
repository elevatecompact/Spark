# wallet-service — Deployment Guide
Critical: PostgreSQL sync replication, Redis Cluster for idempotency, HSM/Vault for processor secrets, SOC2 compliance.
K8s: 3 replicas, 2GB RAM/1 CPU. PDB min 2 available. WAL: logical for PITR. Backups: 30min snapshots + continuous WAL.
Deploy: ./wallet migrate up --dry-run then --confirm. Canary 1 replica 15min.
Health: /health (DB+Redis+Stripe), /ready (replica sync), /metrics :4108.
