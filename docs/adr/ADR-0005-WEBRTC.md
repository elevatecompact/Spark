# ADR-0005: WebRTC for Real-Time Streaming

## Status

Accepted

## Context

Spark's core value proposition includes sub-second latency livestreaming for interactive experiences such as live shopping, real-time Q&A, and collaborative creator-viewer sessions. The platform must support millions of concurrent viewers with latency under 500 milliseconds while maintaining video quality across variable network conditions. HLS and DASH provide excellent scalability and compatibility but introduce 6-30 seconds of latency unsuitable for interactivity. WebRTC, SRT, and low-latency CMAF were evaluated. SRT offered reliable transport over unpredictable networks but required custom player integration. Low-latency CMAF reduced latency to 2-4 seconds but could not meet the sub-second requirement. WebRTC provided the lowest latency with native browser support and adaptive bitrate streaming.

## Decision

Use WebRTC for real-time streaming with Selective Forwarding Units (SFUs) for multi-party scalability. The streaming pipeline captures video via browser MediaRecorder, encodes with hardware-accelerated H.264, and transmits over WebRTC to a regional media server cluster. SFUs forward streams without transcoding to minimize latency and compute cost. A fallback path delivers HLS via a concurrent encoding pipeline for viewers on unsupported networks or devices. Adaptive bitrate encoding uses a ladder of presets tuned for 144p to 4K resolutions. TURN servers provide connectivity for NAT-traversed peers.

## Consequences

### Positive
- Sub-300ms latency enables genuine real-time interactivity
- Native browser support without plugins or proprietary players
- Adaptive bitrate handles variable network conditions gracefully
- SFU architecture scales horizontally for large viewer counts

### Negative
- CPU overhead for encoding on client devices, especially mobile
- TURN relay costs scale with the number of NAT-traversed connections
- Debugging WebRTC issues requires specialized tooling and expertise
- Fallback to HLS increases pipeline complexity and operational cost
