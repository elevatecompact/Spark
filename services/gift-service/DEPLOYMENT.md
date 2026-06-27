# gift-service — Deployment Guide
Requires PG15+, Redis 7+, Kafka. K8s: 3 replicas 512MB/500m, HPA min 2 max 6.
Deploy: ./gift migrate up, ./gift seed-catalog, kubectl apply -f k8s/gift-service/.
Cache: gift catalog TTL 5min, leaderboard TTL 60s, gift card validity checked in Redis first.
Health: /health (DB+Redis+Kafka), /ready (catalog loaded), /metrics :4109.
