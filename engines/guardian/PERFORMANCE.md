# Performance

## Targets
- Token validation (cache hit): < 1ms P99
- Token validation (cache miss): < 5ms P99
- Login flow: < 100ms P99 (password)
- RBAC evaluation: < 2ms P99
- Throughput per node: > 50,000 token validations/second

## Benchmarks
c6i.2xlarge: 80K token validations/sec at 25% CPU. Ed25519 signing at 4us per token. Redis cache hit rate > 99.7%.
