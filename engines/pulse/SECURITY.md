# Security

- DTLS-SRTP encrypts all WebRTC media streams.
- Stream tokens - HMAC-signed JWT tokens required for viewer session creation. Tokens embed streamId, viewerId, and expiry.
- Ingest authentication - RTMP/SRT streams require a pre-shared stream key validated on push.
- Origin validation - SFU nodes verify the control plane's gRPC TLS certificate before accepting assignments.
- Rate limiting - Per-IP connection limits and token bucket for ICE candidate submissions.
- STUN/TURN - TURN relay credentials are time-limited and per-session.
