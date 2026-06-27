# WebRTC Real-Time Streaming Architecture

Spark leverages WebRTC for ultra-low-latency real-time communication between broadcasters and viewers. WebRTC provides sub-second latency essential for interactive features including live chat, virtual gifts, and real-time polling.

## Signaling Infrastructure

WebRTC sessions are established through a global signaling mesh. When a broadcaster starts a stream, the signaling service negotiates SDP offers and answers, exchanges ICE candidates, and establishes peer connections optimized for the lowest latency path.

### Signaling Protocol
`
Broadcaster → Signaling Service → Selective Forwarding Unit (SFU)
     ↓                                          ↓
  Viewer ←─────────────────────────────────── Viewer
`

The signaling service coordinates connection establishment. Once the media path is established, signaling is no longer in the critical path.

## Selective Forwarding Unit (SFU)

Spark deploys a custom SFU that forwards media streams without transcoding to minimize latency. The SFU is deployed at edge locations to keep media paths short.

### SFU Features
- **Simulcast Support**: Broadcasters send multiple quality layers; viewers receive the optimal layer for their connection
- **Forward Error Correction**: Redundant packets reduce retransmission latency
- **Congestion Control**: Google Congestion Control (GCC) adapts to network conditions in real time
- **SVC Scalability**: Scalable Video Coding enables efficient multi-viewer distribution

## Media Path Optimization

### Direct Peer-to-Peer
When broadcaster and viewer are on the same network, a direct P2P connection is established for minimal latency.

### Relay via Edge SFU
For remote viewers, media routes through the nearest edge SFU, which selects the optimal upstream source.

### Fallback to CDN
Viewers with unreliable WebRTC connections fall back to HLS delivery via the CDN, accepting higher latency for stability.

## Data Channels

WebRTC data channels carry real-time chat messages, gift animations, and interactive events. Data channels use ordered, reliable mode for chat and unordered, partially reliable mode for high-frequency events like viewer counts.
