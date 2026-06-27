# Ranking Engine — Content Ranking

## Purpose

The Ranking engine is responsible for ordering content across all Titan surfaces — feeds, search results, recommendations, and browse pages. It combines relevance signals, business rules, and real-time engagement data to produce optimal content ordering.

## Architecture

The Ranking engine operates as a lightweight scoring service invoked by Atlas (search/browse) and Oracle (recommendations). It applies a configurable cascade of scoring functions, from simple business rule filters through ML-scored relevance models.

## Tech Stack

- **Language**: Go
- **Scoring Cache**: Redis (precomputed scores for popular content)
- **Feature Store**: Redis + S3 for cold features
- **ML Inference**: ONNX Runtime for cross-platform model execution
- **Configuration**: Dynamic scoring rules via etcd with hot reload

## Key Features

- **Hybrid scoring**: Combines content relevance, user preference signals, freshness, popularity, and business rules
- **Configurable cascades**: Multi-stage ranking pipeline (filter, coarse score, fine score, rerank)
- **Real-time feature injection**: Engagement signals (clicks, likes, shares, watch time) incorporated within seconds
- **Diversity constraints**: Ensure topical, creator, and format diversity in ranked results
- **Fairness constraints**: Exposure guarantees for underrepresented creators and content categories
- **A/B testable configurations**: Scoring parameters are dynamically configurable per experiment bucket
- **Explanation API**: Returns contributing factors for each scored item for transparency

## Performance Targets

| Metric | Target |
|--------|--------|
| P99 scoring latency | < 20ms |
| Throughput per node | 10,000 items scored/second |
| Feature freshness | < 3 seconds |
| Cascade stages | 5 (filter, freshness, relevance, popularity, diversity rerank) |
| Configuration update latency | < 10 seconds (hot reload) |