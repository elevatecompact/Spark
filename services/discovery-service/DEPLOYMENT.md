# discovery-service — Deployment Guide
K8s: k8s/discovery-service/ — api (3x 1GB), trending-worker (2x 1GB for velocity computation), editorial-api (1x 512MB admin).
Feed construction combines: algorithmic recs from recommendation-service, trending from local trending engine, editorial picks from DB. Multi-source merge with deduplication and diversity rules.
Deploy: kubectl apply -f k8s/discovery-service/. Cache warming on deploy: pre-compute top 10 most popular feed configurations.
Health: /health (DB+Redis+Kafka+recommendation-service), /ready (feed cache warming), /metrics :4124.
