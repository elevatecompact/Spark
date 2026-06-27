# competition-service — Deployment Guide
K8s: k8s/competition-service/ — api (3x 512MB), bracket-engine (2x 1GB CPU-intensive for bracket generation), prize-worker (1x 512MB).
Bracket generation is compute-intensive for large tournaments (128+). Cached in Redis after generation.
Deploy: ./competition migrate up, kubectl apply -f k8s/competition-service/.
Cache: bracket tree TTL 1h (or invalidated on match completion), leaderboard TTL 30s.
Health: /health (DB+Redis+Kafka), /ready (migrations applied), /metrics :4119.
