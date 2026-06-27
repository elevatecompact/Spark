# Performance

## Targets
- Real-time factor: < 0.2 (5x real-time on NVENC for 1080p)
- Transcoding latency (live): < 2s end-to-end
- GPU memory per 4K stream: < 1.5GB
- Max concurrent sessions per GPU: 8 (NVENC)
- Job queue throughput: > 500 concurrent VOD jobs per node

## Benchmarks
On A10G GPU: 8 simultaneous 1080p-to-720p transcodes at 120fps each. Per-title encoding saves 35% bitrate at same SSIM.
