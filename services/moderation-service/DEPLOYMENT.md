# moderation-service — Deployment Guide
Components: API server (3 replicas), ML inference (GPU, 2 replicas for text, 2 for image), review worker (2 replicas for async queue processing).
K8s: k8s/moderation-service/ — api (3x 1GB), ml-text (2x GPU 8GB), ml-image (2x GPU 16GB).
ML models run on dedicated GPU nodes with nodeSelector. Models: BERT-based toxicity classifier, EfficientNet for NSFW detection.
Deploy: kubectl apply -f k8s/moderation-service/. ML models deployed separately via ML platform CI/CD.
Health: /health (DB+Redis+ML endpoints), /ready (models loaded), /metrics :4116.
