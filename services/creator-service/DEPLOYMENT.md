# creator-service — Deployment Guide

## Prerequisites
- PostgreSQL 15+, Redis 7+, S3 bucket for verification docs
- Kubernetes 1.28+ with ingress for file uploads

## Manifests: k8s/creator-service/
3 replicas, 1GB RAM / 1 CPU. HPA: CPU > 65%, min 2 max 8.

## Deployment
`ash
./creator migrate up
kubectl apply -f k8s/creator-service/
kubectl rollout status deploy/creator-service
`
Verification docs use presigned S3 URLs with 24h expiration. Progressive delivery via Argo: 10% → 50% → 100%.

## Health
- /health — DB + Redis + S3
- /ready — Metrics view refreshed
- /metrics — Port 4102
