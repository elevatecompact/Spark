# Changelog

## 0.3.0 (2026-06-15)
- Added SRT ingest protocol support.
- Improved GCC congestion control with neural bandwidth estimator.
- Reduced P95 latency from 600ms to 400ms.
- Added per-stream quality metrics endpoint.

## 0.2.0 (2026-04-10)
- Introduced SFU node registry with Redis.
- Added Prometheus metrics for all media pipelines.
- WebRTC simulcast support with 3-tier bitrate ladder.

## 0.1.0 (2026-02-01)
- Initial release with RTMP ingest and WebRTC playback.
- Single-node SFU with up to 500 concurrent viewers.
