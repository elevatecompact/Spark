# API

## Ingest Endpoints
- POST /ingest/rtmp - Accept RTMP push from encoder.
- POST /ingest/srt - Accept SRT stream with passphrase authentication.
- GET /ingest/status/:streamId - Return ingest bitrate, frame rate, jitter.

## Playback Endpoints
- POST /session - Create a viewer session, return WebRTC SDP offer.
- POST /session/:sessionId/ice - Exchange ICE candidates.
- GET /stream/:streamId/status - Return stream health and connected viewers.
- POST /stream/:streamId/stop - Terminate a live stream.

## Health
- GET /health - Liveness probe.
- GET /metrics - Prometheus metrics.
