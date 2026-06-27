# ADR-0006: Event-Driven Architecture

## Status

Accepted

## Context

Spark's feature set spans content ingestion, transcoding, AI analysis, moderation, notification, recommendation, and payout processing. Early implementations used synchronous HTTP calls between services, resulting in tight coupling, cascading failures when downstream services degraded, and difficulty adding new consumers without modifying producers. As the service count grew, request tracing became complex, and timeout tuning became a source of recurring incidents. The team evaluated pure request-response, event-driven with message broker, and choreographed vs. orchestrated sagas. Purely synchronous architectures were rejected due to coupling. Orchestrated sagas added a central coordinator that became a bottleneck. A choreographed event-driven approach using Kafka provided the best decoupling.

## Decision

Adopt a choreographed event-driven architecture where services communicate exclusively through events published to Apache Kafka topics, with no synchronous inter-service dependencies except for idempotent read paths. Each service owns its event schemas and publishes state-change events. Downstream services subscribe to relevant topics without the publisher knowing their identities. Event schemas are versioned in Avro with Schema Registry enforcing compatibility. Compensating transactions are implemented as event handlers for failure scenarios. Correlation IDs are propagated through event headers for end-to-end tracing.

## Consequences

### Positive
- Loose coupling allows independent deployment and scaling of services
- New consumers can be added without modifying producers
- Event log provides an audit trail of all state changes
- Natural fit for CQRS and event-sourcing patterns

### Negative
- Eventual consistency requires careful handling of read-your-writes guarantees
- Debugging multi-hop event flows demands robust tracing infrastructure
- Schema evolution must be managed with backward/forward compatibility
- Exactly-once delivery guarantees are difficult to achieve end-to-end
