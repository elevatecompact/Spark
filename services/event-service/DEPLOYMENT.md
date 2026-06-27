# event-service — Deployment Guide
K8s: k8s/event-service/ — api (3x 512MB), reminder-worker (2x 512MB for sending event reminders).
Event reminders: CronJob that checks events starting in 24h/1h and triggers notification. Ticket sales: inventory check with optimistic locking.
Deploy: ./event migrate up, kubectl apply -f k8s/event-service/.
Cache: event detail TTL 5min, ticket availability TTL 30s (near-sold-out events uncached).
Health: /health (DB+Redis+Kafka), /ready (migrations applied), /metrics :4118.
