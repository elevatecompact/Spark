# Performance

## Targets
- P99 inference latency: < 20ms (single), < 100ms (batch 64)
- Feature lookup latency: < 5ms P99 from Redis
- Vector search latency: < 50ms P99 for top-200 candidates
- End-to-end recommendation: < 150ms P99
- Throughput per node: > 2000 QPS

## Benchmarks
ONNX runtime on c6i.8xlarge: embedding model at 8ms per inference. Redis cluster with 6 shards handles 500K feature read/write per second.
