# Performance

## Targets
- Glass-to-glass latency: < 400ms P95
- Ingest to SFU distribution: < 50ms P99
- Viewer connection time: < 1s (ICE + DTLS handshake)
- Maximum viewers per node: 10,000 (with hardware TURN)
- Packet loss recovery: < 200ms for NACK retransmission
- CPU per 1000 streams: < 4 cores (modern Xeon)

## Benchmarks
Measured on c5.4xlarge instances: 500 concurrent 1080p30 streams consumed 35% CPU. WebRTC simulcast at 3 bitrate tiers used 12 Mbps egress per viewer.
