# Domain Model

The Spark domain model captures the core business concepts and their relationships. Following domain-driven design principles, the model is organized into aggregates that enforce consistency boundaries.

## Core Entities

### User
The root identity for all platform participants. Users have roles: Viewer, Broadcaster, Moderator, and Admin. Each User carries a profile, authentication credentials, and preference settings.

### Broadcaster
A User role that creates and manages live streams. Broadcasters own channels, configure stream settings, and manage monetization options.

### Stream
A live video broadcast event. Each Stream has metadata (title, category, tags), an ingest endpoint, transcoding configuration, and analytics tracking. Streams transition through states: scheduled, live, paused, ended, and archived.

### Viewer
A User role that consumes live content. Viewers can chat, send virtual gifts, subscribe to channels, and interact with polls or Q&A sessions.

### Channel
A persistent content hub owned by a Broadcaster. Channels organize streams historically and serve as the primary subscription target for viewers.

### VirtualGood
Digital items viewers can purchase and send during streams. Includes currency units, gifts, badges, and subscription tiers.

### ModerationAction
Records of moderation decisions including user mutes, content removal, and stream termination. Immutable and fully auditable.

## Aggregates

| Aggregate | Root Entity | Key Invariants |
|-----------|-------------|----------------|
| Identity | User | Email uniqueness, authentication rules |
| Streaming | Stream | Single active stream per broadcaster |
| Wallet | Account | Non-negative balance, transaction atomicity |
| Moderation | ModerationAction | Immutable audit trail |

## Value Objects

- **GeoLocation**: Latitude/longitude pair used for edge routing
- **BitrateProfile**: Resolution, frame rate, and bitrate settings
- **StreamMetrics**: Concurrent viewers, duration, quality scores
- **PaymentRecord**: Transaction ID, amount, currency, timestamp

## Domain Events

Significant domain events include StreamStarted, StreamEnded, GiftSent, UserSubscribed, ModerationActionTaken, and PaymentProcessed. These events flow through Kafka and drive downstream projections and analytics.
