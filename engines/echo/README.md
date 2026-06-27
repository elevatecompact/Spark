# Echo Engine

**Purpose:** Real-time messaging and notification engine for live events and alerts.
**Tech Stack:** Rust, WebSocket, Tokio, Redis Pub/Sub, gRPC, Kafka.

Echo provides bidirectional real-time communication between Titan services and connected clients. It handles WebSocket connections, fan-out messaging to channels, presence tracking, and reliable delivery with exactly-once semantics. The backbone for live chat, notifications, and server-sent events across the platform.
