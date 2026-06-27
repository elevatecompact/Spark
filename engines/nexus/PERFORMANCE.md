# Performance

## Targets
- Upload throughput: > 1Gbps per node to S3
- Download latency: < 5ms P99 presigned URL generation
- Metadata lookup (cached): < 3ms P99
- Image transformation: < 100ms P99 (resize)
- Video thumbnail: < 3s P99 for 1080p
- Upload durability: 99.999999999%

## Benchmarks
Direct-to-S3 at 800Mbps on c6i.2xlarge. Redis metadata at 1.2ms mean. ImageMagick resize at 45ms for 4K to 320px.
