# viewer-service — Deployment Guide
K8s manifests in k8s/viewer-service/: 4 replicas, 512MB RAM / 500m CPU. HPA scales on CPU > 70% or Kafka queue depth > 1000. Heavy write pattern on watch_history — ensure adequate WAL capacity. Run partition maintenance via cron: ./viewer partitions manage.
Deploy: kubectl apply -f k8s/viewer-service/ then kubectl rollout status deploy/viewer-service.
History cached in Redis for last 50 entries per viewer (TTL 5min). Bookmarks cached in Redis (TTL 1min).
