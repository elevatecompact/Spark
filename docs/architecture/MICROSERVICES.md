# Microservices Architecture

Spark implements a microservices architecture to enable independent deployment, scaling, and development of loosely coupled services. Each microservice owns a bounded context and communicates through well-defined APIs.

## Service Decomposition

The platform consists of the following primary services:

### Identity Service
Handles user registration, authentication, and profile management. Issues JWT tokens for API authentication. Integrates with third-party OAuth providers. Data store: PostgreSQL with read replicas for authentication queries.

### Stream Service
Manages stream lifecycle from creation to archival. Controls ingest endpoints, transcoding job orchestration, and stream state transitions. Uses Redis for real-time stream status and Kafka for state change events.

### Discovery Service
Provides search and recommendation capabilities. Indexes content metadata into OpenSearch and maintains personalized recommendation models using collaborative filtering and content-based approaches.

### Wallet Service
Manages virtual currency balances, gift transactions, and payment processing. Ensures transactional integrity with optimistic concurrency control. Emits financial events to Kafka for downstream reconciliation.

### Moderation Service
Processes content reports, applies automated filtering via ML models, and records moderation actions. Maintains a fully auditable event log of all enforcement decisions.

### Chat Service
Powers real-time chat via WebRTC data channels with a Kafka-backed persistence layer. Maintains chat history for replay and moderation review.

### Notification Service
Delivers push notifications, in-app alerts, and email digests. Uses a priority queue to ensure time-sensitive notifications are delivered promptly.

## Inter-Service Communication

| Pattern | Protocol | Use Case |
|---------|----------|----------|
| Synchronous | gRPC | Query-heavy, low-latency operations |
| Asynchronous | Kafka | Event notification and state propagation |
| Synchronous | REST | Administrative and external APIs |

## Service Mesh

All services run on Kubernetes with an Istio service mesh providing mutual TLS, traffic management, observability, and fine-grained access control. Sidecar proxies handle inter-service communication with automatic retries and circuit breaking.
