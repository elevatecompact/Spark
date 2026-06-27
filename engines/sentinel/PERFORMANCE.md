# Performance

## Targets
- Tier 1 latency: < 1ms P99
- Tier 2 latency: < 50ms P99 (text), < 150ms P99 (image)
- Tier 3 latency: < 500ms P99 (text), < 2s P99 (image)
- Video moderation: < 30s for 10-minute video
- Throughput: > 500 req/s (tier 1+2), > 50 req/s (tier 3)

## Benchmarks
On g5.xlarge: DistilBERT at 12ms, ResNet-18 at 65ms. Blocklist matching at 0.3ms for 10K entries. CLIP analysis at 1.2s per image.
