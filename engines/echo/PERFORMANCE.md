# Performance

## Targets
- Message delivery latency: < 10ms P99 (same node), < 50ms P99 (cross-node)
- Connection throughput: > 50,000 conns/second per node
- Message throughput: > 500,000 msgs/second per node
- Memory per connection: < 20KB
- Durable delivery: < 500ms P99 end-to-end

## Benchmarks
c6i.4xlarge: 100K concurrent connections at 32% CPU. Broadcast to 10K subscribers in 8ms P99. Redis Pub/Sub fan-out to 50 nodes adds 35ms P99.
