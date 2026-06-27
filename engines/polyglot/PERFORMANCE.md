# Performance

## Targets
- Translation latency (batch): < 200ms P99 per sentence (GPU)
- Translation latency (streaming): < 500ms P99 per segment
- Language detection: < 10ms P99
- Throughput per GPU: > 500 sentences/second
- Batch throughput: > 10,000 sentences/second (batch 64)
- Supported languages: 200+ detection, 100+ translation

## Benchmarks
A10G GPU: NLLB-200 at 200 sentences/sec (batch 32). Opus-MT (en-es) at 800 sentences/sec. Language detection at 3ms.
