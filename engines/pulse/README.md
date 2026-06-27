# Pulse Engine

**Purpose:** Real-time live streaming engine built with Rust and WebRTC.
**Tech Stack:** Rust, WebRTC, Tokio, SRT, MPEG-TS, HLS, RTMP ingest.

Pulse provides ultra-low-latency live video streaming with sub-second glass-to-glass delay. It handles ingest from broadcast encoders via RTMP/SRT, transcodes to adaptive bitrate ladder, and distributes via WebRTC to viewers. Designed for interactive livestream scenarios including gaming, auctions, and live shopping where latency below 500ms is mandatory.
