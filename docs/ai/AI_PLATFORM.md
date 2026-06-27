# AI Platform Architecture

Spark AI is built on a layered architecture separating infrastructure, model serving, orchestration, and application logic.

## Layer 1 — Infrastructure

Kubernetes cluster (EKS/GKE) with GPU node pools (A100/H100 for training, T4/L4 for inference). All models are containerized using NVIDIA CUDA base images and orchestrated via KubeFlow. Storage is handled by S3-compatible object stores for datasets and model artifacts, with EFS/NFS for shared feature stores.

## Layer 2 — Model Serving

The Inference Gateway (Python/FastAPI + Triton Inference Server) routes requests to the appropriate model. Serving strategies include:
- **Real-time**: Triton with dynamic batching for latency-sensitive workloads (moderation, voice)
- **Batch**: Spark jobs via KubeFlow for offline inference (clip generation, thumbnail scoring)
- **Streaming**: Kafka streams for continuous inference on live data (recommendation updates, translation)

## Layer 3 — Orchestration

KubeFlow Pipelines orchestrate multi-step ML workflows: data ingestion → feature engineering → training → evaluation → deployment. MLflow tracks experiments, parameters, and metrics. Feast serves online and offline feature values.

## Layer 4 — Application

Each AI capability exposes a gRPC API consumed by the Spark backend services. The API Gateway (Envoy) handles authentication, rate limiting, and canary routing. Event provenance is tracked via OpenTelemetry for full observability.

## Cross-Cutting Concerns

- **Observability**: Prometheus + Grafana dashboards for model latency, throughput, and drift
- **Security**: Vault for secrets, mTLS between services, model input/output scanning
- **Governance**: Every inference carries a traceable lineage ID for audit and replay
