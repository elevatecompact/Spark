# advertising-service — Deployment Guide
K8s: k8s/advertising-service/ — api (3x 1GB), ad-server (5x 1GB latency-optimized), impression-worker (4x 1GB for event processing), fraud-detection (2x 1GB Spark).
Ad server is latency-critical (p99 < 50ms). Must be co-located with stream service for pre-roll decisions.
Deploy: ./advertising migrate up, kubectl apply -f k8s/advertising-service/. Warm caches on deploy.
Health: /health (DB+Redis+Kafka+ClickHouse), /ready (inventory loaded), /metrics :4120.
