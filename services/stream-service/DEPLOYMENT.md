# stream-service — Deployment Guide
Components: API server (stateless), Ingest Nodes (Nginx-RTMP edge), Transcoder Fleet (FFmpeg GPU), Origin Nodes.
K8s: Standard deployment for API, DaemonSets for ingest/transcoder on GPU nodes (NVIDIA T4+).
Deploy: kubectl apply -f k8s/stream-service/ + kubectl apply -f k8s/stream-transcoder/.
Health: /health (DB+Redis+S3+transcoder), /ready (ingest configured), /metrics :4103.
