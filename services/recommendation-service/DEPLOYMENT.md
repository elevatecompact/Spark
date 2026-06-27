# recommendation-service — Deployment Guide
Components: API server (stateless), Inference server (ONNX Runtime, GPU), Feature computation job (Spark, 4 executors), Model training pipeline (TensorFlow, scheduled weekly).
K8s: k8s/recommendation-service/ — api (3x 1GB), inference (2x GPU 16GB), feature-job (CronJob).
Model deployment: Blue-green with A/B testing. New model receives 10% traffic for 24h before full cutover.
Deploy: kubectl apply -f k8s/recommendation-service/, verify inference latency, monitor CTR.
