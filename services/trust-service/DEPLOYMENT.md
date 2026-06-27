# trust-service — Deployment Guide
K8s: k8s/trust-service/ — api (3x 1GB), reputation-worker (2x 1GB for daily recalculation), fraud-detector (2x 2GB for ML-based fraud).
Reputation recalculation runs daily at 3AM. Processes all users with signals updated since last run (~5min for 1M users).
Fraud detection uses rule-based + ML model (XGBoost) deployed on fraud-detector pods.
Deploy: ./trust migrate up, kubectl apply -f k8s/trust-service/. Seed default risk rules on first deploy.
Health: /health (DB+Redis+Kafka), /ready (reputation model loaded), /metrics :4125.
