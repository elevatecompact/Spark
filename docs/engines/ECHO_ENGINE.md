# Echo Engine — Real-Time Messaging

## Purpose

Echo is Titan's real-time messaging engine. It powers live chat during streams, direct messaging between users, and platform-wide notification delivery. Echo guarantees ordered, durable message delivery at scale.

## Architecture

Echo uses a fan-out publish/subscribe model. Messages are ingested via WebSocket connections, persisted to a message log (Kafka), and fanned out to subscribers via per-channel WebSocket connections managed by connection brokers.

## Tech Stack

- **Language**: Go
- **Transport**: WebSocket (with SockJS fallback for browsers)
- **Message Log**: Apache Kafka (retention-based, configurable per channel type)
- **Connection Management**: Redis for presence tracking and channel membership
- **Deployment**: Kubernetes with HPA based on WebSocket connection count

## Key Features

- **Ordered delivery**: Messages within a channel are delivered in strict order with deduplication IDs
- **Presence**: Real-time online/offline and typing indicators
- **Channel types**: Public (stream chat), private (DMs), system (notifications)
- **History**: Configurable message retention per channel (stream: 7 days, DMs: permanent)
- **Moderation integration**: Echo streams all messages to Sentinel for real-time moderation
- **Scalability**: Horizontal scaling via consistent hashing of channel ID to broker node
- **Delivery guarantees**: At-least-once delivery with application-level acknowledgments

## Performance Targets

| Metric | Target |
|--------|--------|
| End-to-end delivery latency | < 100ms (p99) |
| Connections per node | 100,000 concurrent WebSocket connections |
| Message throughput | 1M messages/second |
| Channel fan-out | 100,000 subscribers per channel |
| Uptime | 99.999% |