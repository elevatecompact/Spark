# Performance

## Targets
- Payment intent creation: < 50ms P99
- Payment confirmation: < 200ms P99
- Subscription query: < 10ms P99 (cached)
- Invoice generation: < 500ms P99
- Webhook processing: < 100ms P99
- Throughput: > 1000 intent creations/second

## Benchmarks
c6i.2xlarge: Stripe API at 120ms end-to-end. Invoice generation at 200ms. Redis-cached lookups at 2ms P99.
