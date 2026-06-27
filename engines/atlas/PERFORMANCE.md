# Performance

## Targets
- Lookup latency: < 2ms P99 (cache hit), < 20ms P99 (cache miss)
- Registration propagation: < 500ms P99 across all agents
- Watch event delivery: < 100ms P99
- Agent memory: < 50MB per node (10,000 services)
- Registrar throughput: > 100,000 registrations/second
- Health check throughput: > 10,000 checks/second per registrar

## Benchmarks
etcd-backed cluster with 3 registrar nodes handles 50,000 service registrations with 200ms P99 propagation. Agent cache hit rate exceeds 99.5% under steady-state workload.
