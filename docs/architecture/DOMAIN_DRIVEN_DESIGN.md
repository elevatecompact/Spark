# Domain-Driven Design

Spark adopts Domain-Driven Design (DDD) as the primary methodology for modeling complex business domains. DDD enables the team to maintain a shared understanding of the problem space and align software boundaries with business capabilities.

## Bounded Contexts

The platform is decomposed into the following bounded contexts, each with its own ubiquitous language and internal model:

### Identity and Access Context
Manages user registration, authentication, authorization, and profile management. Uses OAuth 2.0 and OpenID Connect for federation. The core entity is the User aggregate with role-based access control.

### Streaming Context
Owns the live broadcasting lifecycle including stream ingestion, transcoding, and delivery. The Stream aggregate enforces the invariant that a broadcaster may only have one active stream at a time. WebRTC signaling and media negotiation live within this context.

### Discovery Context
Handles content recommendation, search, and categorization. Operates primarily on read models projected from domain events. Uses OpenSearch for full-text search and collaborative filtering for recommendations.

### Monetization Context
Manages virtual currency, gift purchasing, subscriptions, and revenue sharing. The Account and Transaction aggregates ensure financial consistency. All financial events are immutably logged.

### Moderation Context
Governs content review, user reporting, and enforcement actions. ModerationAction entities are immutable and provide a complete audit trail.

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Spark | A live streaming event on the platform |
| Stage | A channel or content hub |
| Flame | Virtual currency unit |
| Spotlight | Featured content promotion |
| Echo | Real-time chat message in a stream |

## Strategic Design

Context mapping uses the following relationship patterns:

- **Partnership**: Streaming and Monetization contexts collaborate closely on subscription flows
- **Conformist**: Discovery context consumes events from Streaming context downstream
- **Open Host Service**: Identity context publishes user lifecycle events via Kafka
- **Anti-Corruption Layer**: Legacy analytics integration uses an ACL to translate between models

Tactical patterns include Entity, Value Object, Aggregate Root, Domain Event, Repository, and Domain Service where appropriate.
