# commerce-service — Deployment Guide
K8s: k8s/commerce-service/ — api (3x 1GB), fulfillment-worker (2x 1GB for digital delivery), merchant-payout (CronJob monthly).
Fulfillment: digital goods delivered via S3 presigned URLs (72h expiry). Inventory managed with optimistic locking for limited items.
Deploy: ./commerce migrate up, ./commerce seed-categories, kubectl apply -f k8s/commerce-service/.
Cache: product detail TTL 1min, storefront TTL 30s, cart persisted in DB (cached in Redis).
Health: /health (DB+Redis+Kafka+S3), /ready (migrations applied), /metrics :4121.
