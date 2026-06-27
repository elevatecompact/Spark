# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| model.path | /models/current.onnx | Path to deployed ONNX model |
| model.batch_size | 64 | Inference batch size |
| redis.host | localhost | Redis host for feature store |
| redis.ttl | 86400 | User feature TTL in seconds |
| milvus.host | localhost | Milvus host for embedding store |
| milvus.top_k | 200 | Candidates retrieved from vector search |
| ranker.diversity_alpha | 0.3 | Diversity vs relevance tradeoff |
| ranker.exploration_epsilon | 0.05 | Epsilon-greedy exploration rate |
| kafka.brokers | localhost:9092 | Kafka broker list |
