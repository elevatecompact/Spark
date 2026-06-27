# Performance

## Targets
- Ingestion throughput: > 1M events/second per node
- Query latency (aggregated): < 50ms P99 for 30-day window
- Query latency (raw): < 500ms P99 for 7-day window
- Dashboard render: < 2s P99
- Storage compression: > 5:1 (LZ4), > 8:1 (ZSTD)

## Benchmarks
c6i.4xlarge: 2.5M events/sec with LZ4. Aggregation over 1B rows in 120ms. Materialised view reduces query time 40x.
