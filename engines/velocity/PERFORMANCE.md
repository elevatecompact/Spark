# Performance

## Targets
- Cache purge propagation: < 1s P99 (single), < 5s P99 (tag)
- Cache warming throughput: > 10,000 URLs/min per provider
- Steering decision latency: < 10ms P99
- Provider failover time: < 30s (auto), < 5s (manual)
- Origin shield hit rate: > 90%

## Benchmarks
Fastly purge: 300ms mean single URL, 2.5s for 500-URL batch. Warming at 12,000 URLs/min per thread. Edge probes measure latency within 5ms accuracy.
