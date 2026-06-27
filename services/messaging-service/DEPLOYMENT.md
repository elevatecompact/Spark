# messaging-service — Deployment Guide
Architecture: Stateless REST API + WS gateway for real-time. PostgreSQL for messages, Redis for presence/typing.
K8s: k8s/messaging-service/ — api-deployment (3x 1GB), ws-deployment (4x 2GB connection-heavy).
High write on messages — ensure conversation_id index, partition by month.
Deploy: ./messaging migrate up then kubectl apply -f k8s/messaging-service/.
Fallback: If Redis down, poll-based delivery. Attachments served from S3 via CDN.
