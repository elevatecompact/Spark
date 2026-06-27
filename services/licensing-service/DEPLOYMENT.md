# licensing-service — Deployment Guide
K8s: k8s/licensing-service/ — api (3x 512MB), royalty-worker (2x 1GB for monthly calculation), compliance-worker (1x 512MB for alerts).
Royalty calculation is compute-intensive for large catalogs — run as monthly batch job with progress tracking.
Deploy: ./licensing migrate up, kubectl apply -f k8s/licensing-service/.
Cache: license validity TTL 5min (short for near-expiry accuracy), content rights TTL 1min.
Health: /health (DB+Redis+Kafka), /ready (migrations applied), /metrics :4123.
