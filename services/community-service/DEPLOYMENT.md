# community-service — Deployment Guide
K8s: k8s/community-service/ — api (3x 1GB), post-worker (2x 1GB for async processing of post notifications).
Moderate write load on posts/comments. Key indexes: community_id+created_at, author_id.
Deploy: ./community migrate up, then kubectl apply -f k8s/community-service/.
Cache: community metadata TTL 5min, member counts TTL 1min, post lists TTL 30s.
Health: /health (DB+Redis+Kafka), /ready (migrations applied), /metrics :4117.
