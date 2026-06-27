# Architecture

Pulse follows a peer-to-peer hybrid architecture with selective forwarding units (SFUs). The ingest pipeline receives RTMP or SRT streams and feeds them into a Rust-based media router. The router demuxes, repacketizes into RTP, and fans out to viewer-connected SFU nodes. Each SFU maintains a WebRTC peer connection per viewer and handles congestion control, NACK retransmission, and adaptive bitrate selection. Control plane uses a separate gRPC channel for signalling, room management, and stream health monitoring.
