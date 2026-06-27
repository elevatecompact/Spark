# Video Streaming Pipeline

Spark's video streaming pipeline ingests live video from broadcasters, processes it through transcoding and packaging, and delivers optimized streams to viewers worldwide. The pipeline is designed for sub-second ingest-to-playback latency.

## Ingest

Broadcasters connect via RTMP or WebRTC. The ingest layer runs on edge nodes closest to the broadcaster to minimize upload latency.

### Ingest Protocols
| Protocol | Use Case | Latency |
|----------|----------|---------|
| RTMP | Professional broadcasters with OBS | 2-5 seconds |
| WebRTC | Mobile broadcasters | <500ms |
| SRT | High-reliability scenarios | 1-3 seconds |

## Transcoding

The transcoding farm runs on Kubernetes with GPU-backed nodes for hardware-accelerated encoding. A job queue manages transcoding tasks with priority based on stream popularity.

### Output Ladder
`
Source: 1080p  60fps  8Mbps
   ├── 1080p  30fps  6Mbps  (high)
   ├── 720p   30fps  3Mbps  (medium)
   ├── 480p   30fps  1.5Mbps  (low)
   └── 360p   30fps  800Kbps  (mobile)
`

Each rendition is segmented into 2-second HLS chunks and 1-second CMAF fragments for low-latency DASH playback.

## Packaging & Delivery

Transcoded segments are written to S3-compatible storage via the Nexus engine and immediately published via CDN cache invalidation.

### Packaging Formats
- **HLS**: Primary format for broad compatibility
- **CMAF**: Low-latency format for modern browsers
- **LL-HLS**: Apple's low-latency extension for iOS devices

## Adaptive Bitrate (ABR)

The player client monitors bandwidth and buffer health, selecting the optimal rendition dynamically. ABR logic prefers the highest sustainable bitrate with a 10-second buffer target.

### ABL Decision Factors
- Measured download speed (moving average)
- Buffer occupancy
- Device capabilities (screen size, codec support)
- Network type (WiFi, cellular, ethernet)

Each rendition is independently cacheable at CDN edges for maximum hit rates.
