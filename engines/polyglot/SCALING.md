# Scaling

Polyglot scales by language pair sharding. High-traffic pairs get dedicated GPU nodes with pre-loaded models. Long-tail pairs share nodes with model swapping. Redis-backed router dispatches requests to correct node. Batch translation uses Kafka for job distribution. Streaming requires dedicated nodes pinning models in GPU memory for sub-500ms latency. Shared Redis cache for maximum hit rate across all nodes.
