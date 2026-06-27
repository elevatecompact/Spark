# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| encoder.software | false | Force software encoding |
| encoder.hardware | nvdec | Hardware decoder type |
| encoder.target_usage | 4 | NVENC target usage (1-7) |
| encoder.max_concurrent_sessions | 8 | Max encoder sessions per GPU |
| segment.duration | 6 | Segment duration in seconds |
| segment.format | fmp4 | Segment container format |
| packaging.hls.enabled | true | Enable HLS manifest generation |
| packaging.dash.enabled | true | Enable DASH manifest generation |
| queue.max_jobs | 1000 | Max queued transcoding jobs |
