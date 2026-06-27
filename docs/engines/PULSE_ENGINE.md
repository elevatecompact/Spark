# Pulse Engine — Live Streaming

## Purpose

Pulse is Titan's live streaming engine. It powers real-time video broadcast with sub-second latency across global audiences, handling ingestion from RTMP and WebRTC sources and delivery via WebRTC and HLS.

## Architecture

Pulse is built in Rust for maximum performance and memory safety. The engine uses a distributed media relay architecture with selective forwarding units (SFUs) for WebRTC and a transcoding pipeline for HLS fallback.

## Tech Stack

- **Language**: Rust
- **Transport**: WebRTC (SDP, ICE, DTLS-SRTP), RTMP ingestion
- **Media Pipeline**: GStreamer for transcoding, Opus (audio), H.264/H.265 (video)
- **Coordination**: Redis Streams for signaling, etcd for SFU membership
- **Deployment**: Kubernetes StatefulSet with GPU-accelerated nodes

## Key Features

- **Sub-second latency**: WebRTC end-to-end latency under 500ms
- **Adaptive bitrate**: Per-viewer ABR with configurable ladder (1080p, 720p, 480p, 360p)
- **Simulcast**: Receiver-driven simulcast for bandwidth optimization
- **Scalable SFUs**: Horizontal scaling of selective forwarding units without global state
- **Recording**: On-demand archive recording to Nexus (S3-compatible storage)
- **Health checks**: Real-time stream health monitoring with automatic recovery

## Performance Targets

| Metric | Target |
|--------|--------|
| End-to-end latency | < 500ms (WebRTC) |
| Concurrent streams per SFU | 10,000 |
| Time to first frame | < 2 seconds |
| Recording accuracy | 99.99% |
| Uptime | 99.99% |