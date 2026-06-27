# Feature Store for ML

The Spark Feature Store (built on Feast) serves as the centralized repository for ML features, ensuring consistency between training and serving, reducing duplication, and enabling feature sharing across teams.

## Architecture

### Online Store
- **Backend**: Redis Cluster with read-replicas for low-latency feature retrieval
- **Use Case**: Real-time inference features (user watch history embeddings, session context, content features)
- **API**: gRPC with batch key lookup and point-in-time queries
- **Freshness**: Sub-second feature updates via Kafka stream ingestion

### Offline Store
- **Backend**: S3-compatible object store with Parquet format, partitioned by date and entity
- **Use Case**: Training dataset generation, batch inference, historical analysis
- **API**: Spark DataFrame and SQL interfaces for feature retrieval
- **Freshness**: Daily batch updates with hourly micro-batches for critical features

## Feature Categories

| Category | Examples | Update Frequency |
|---|---|---|
| User Features | Embedding vectors, engagement aggregates, topic affinities | Real-time + Daily |
| Content Features | Multi-modal embeddings, metadata tags, quality scores | On-ingestion |
| Context Features | Trending scores, seasonal boosts, viral coefficients | Real-time |
| Derived Features | Rate-of-change metrics, interaction ratios, recency decay | Hourly |

## Feature Registration

Features are defined in a declarative YAML spec with name, owner, data type, description, source, freshness SLA, and monitoring rules. Registration triggers automated validation including type checking, null ratio analysis, and distribution comparison with production.

## Point-in-Time Correctness

The feature store guarantees point-in-time correct join for training datasets — features are joined using the exact timestamp of the label event, preventing data leakage and ensuring training/serving consistency.
