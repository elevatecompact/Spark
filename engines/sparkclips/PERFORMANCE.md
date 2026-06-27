# Performance

## Targets
- Highlight detection (real-time): < 30s delay from live event
- Highlight detection (VOD): < 10% of video duration
- Clip generation: < 60s for 30-second 1080p clip
- Throughput: > 50 VOD hours/hour per GPU

## Benchmarks
A10G GPU: Whisper medium processes 1-hour audio in 8 minutes. CLIP at 5fps on 1080p. Audio spectrogram at 2x real-time on CPU.
