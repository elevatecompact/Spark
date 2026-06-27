# Event-Driven Architecture

Spark relies on an event-driven architecture to enable loose coupling, scale independently, and provide near-real-time reactivity across the platform. Apache Kafka serves as the central event backbone.

## Event Topology

The event infrastructure is organized into distinct domains with dedicated Kafka clusters for isolation:

### Core Event Cluster
Handles high-volume domain events including stream lifecycle, user activity, and financial transactions. Topics are partitioned by aggregate ID for ordered processing within a partition. Retention is configured for 7 days with compacted topics for key-value state.

### Analytics Event Cluster
Captures raw behavioral events for offline processing. Includes viewer heartbeats, stream quality metrics, and interaction events. Higher throughput with larger retention (30 days) supports batch analytics and ML training pipelines.

### Audit Event Cluster
Stores immutable records of compliance-relevant events including moderation actions, financial transactions, and access control changes. Infinite retention with tiered storage to S3 via Kafka Connect.

## Event Schema

All events are serialized using Apache Avro with schema registry enforcement. Schemas evolve through backward-compatible additions. Each event includes mandatory headers:

`json
{
  "eventId": "uuid",
  "eventType": "StreamStarted",
  "source": "stream-service",
  "timestamp": "2026-06-27T12:00:00Z",
  "correlationId": "uuid",
  "data": {}
}
`

## Key Event Flows

### Stream Lifecycle
StreamService publishes StreamScheduled, StreamStarted, StreamEnded. WalletService subscribes to bill broadcasters based on duration. AnalyticsService updates dashboards. NotificationService alerts subscribers.

### Monetization
User sends a gift → WalletService emits GiftSent → StreamService updates gift leaderboard → ChatService broadcasts gift animation → AnalyticsService records revenue event.

## Processing Guarantees

- **At-least-once** delivery for all domain events
- **Exactly-once** semantics for financial event processing using Kafka transactions
- **Ordered processing** within stream aggregates via partition key routing

Dead letter queues capture failed events for manual investigation and replay.
