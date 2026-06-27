# Event Sourcing Strategy

Spark uses event sourcing for domains where auditability, temporal queries, and state reconstruction are critical. Instead of storing current state, the system persists a sequence of state-changing events and derives current state by replaying them.

## Applicable Domains

Event sourcing is applied selectively based on business requirements:

### Wallet and Financial Transactions
Every financial event is immutably recorded. Balance is projected from the event stream. This provides a complete audit trail and enables point-in-time balance queries. Compensation events undo prior transactions without mutating history.

### Stream Lifecycle
Stream state changes (scheduled, started, paused, ended, archived) are captured as an ordered event sequence. This enables analytics on stream duration, interruption frequency, and viewer engagement patterns over time.

### Moderation Actions
Moderation decisions are permanently recorded. The event stream supports compliance audits and enables reversal workflows where a moderator overturns a prior decision.

### User Activity History
Login events, profile changes, and preference updates are logged for security analysis and user support.

## Event Store Implementation

Events are stored in Kafka compacted topics with Avro serialization. Key characteristics:

- **Immutable**: Events are append-only and never modified
- **Versioned**: Each event carries an aggregate version for concurrency control
- **Partitioned**: By aggregate ID for ordered replay
- **Retained**: Forever for regulatory compliance

## Snapshot Strategy

To prevent unbounded replay times, periodic snapshots capture aggregate state at version intervals:

`
SnapshotInterval = 100 events
`

Snapshots are stored in S3-compatible storage with the Nexus engine for fast retrieval. On load, the system reads the latest snapshot and replays only subsequent events.

## State Reconstruction

`json
{
  "Stream:abc123": {
    "snapshot": { "state": "live", "startedAt": "...", "viewers": 1500 },
    "version": 42,
    "pendingEvents": [
      {"type": "ViewerJoined", "version": 43},
      {"type": "ViewerLeft", "version": 44}
    ]
  }
}
`

Projections consume the event stream to build read models for queries, analytics, and search indexes. They maintain their own offset tracking for exactly-once semantics.
