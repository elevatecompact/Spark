# recommendation-service — Database Schema
## PostgreSQL — Feature store & metadata
### user_embeddings: user_id UUID PK, embedding FLOAT[] (128-dim), model_version, updated_at
### content_embeddings: content_id UUID PK, embedding FLOAT[] (128-dim), model_version, updated_at
### user_content_interactions: user_id+content_id PK, interaction_type(click,watch,rate,subscribe,dismiss), weight FLOAT, timestamp
## Redis — Online feature cache (user features TTL 1h, content features TTL 1h)
## S3 — Model artifacts (TensorFlow SavedModel, ONNX files, training data snapshots)
## Feature computation: Real-time features from Kafka streams, batch features from ClickHouse
