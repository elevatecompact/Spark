# Oracle Engine — Recommendation Engine

## Purpose

Oracle is Titan's machine learning recommendation engine. It powers personalized content discovery across the platform, delivering relevant videos, streams, and clips to every user based on their behavior, preferences, and contextual signals.

## Architecture

Oracle follows a two-stage retrieval plus ranking architecture. Candidate generation uses collaborative filtering and content-based embeddings, followed by a lightweight ranking model that scores candidates in real time.

## Tech Stack

- **Language**: Python (ML training), Go (inference serving)
- **ML Framework**: TensorFlow, PyTorch for model training
- **Feature Store**: Redis-backed real-time feature store with S3 cold storage
- **Vector Database**: Milvus for approximate nearest neighbor (ANN) search
- **Inference**: TensorFlow Serving with GPU acceleration
- **Orchestration**: Airflow for batch training pipelines

## Key Features

- **Personalized ranking**: Per-user model with real-time feature updates (watch history, dwell time, likes, shares)
- **Cold start**: Content-based embeddings for new users and new content
- **Contextual awareness**: Time-of-day, device type, network conditions influence recommendations
- **A/B testing framework**: Multi-armed bandit for model comparison without full rollout
- **Real-time feedback loop**: User interactions update embeddings within seconds
- **Explainability**: Feature attribution scores returned with every recommendation for transparency

## Performance Targets

| Metric | Target |
|--------|--------|
| P99 inference latency | < 50ms |
| Throughput per node | 5,000 req/s |
| Recommendation relevance (NDCG@10) | > 0.75 |
| Training pipeline frequency | Daily (incremental), Weekly (full retrain) |
| Feature freshness | < 5 seconds from user action to feature update |