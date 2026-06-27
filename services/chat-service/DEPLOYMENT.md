# chat-service — Deployment Guide
Architecture: WebSocket nodes (stateful, sticky sessions via session affinity), REST API nodes (stateless), Redis pub/sub for inter-node fan-out.
K8s: k8s/chat-service/ — websocket-deployment.yaml (5 replicas, 1GB), api-deployment.yaml (3 replicas, 512MB).
WebSocket nodes scale on active connections. Connection draining: 30s graceful shutdown.
Deploy: kubectl apply -f k8s/chat-service/ then kubectl rollout restart deploy/chat-websocket.
