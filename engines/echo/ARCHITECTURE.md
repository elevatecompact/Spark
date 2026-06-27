# Architecture

Echo uses a connection gateway plus channel broker pattern. The gateway manages WebSocket connections using Tokio tasks - one per client. Each gateway node maintains an in-memory connection table indexed by session ID. Messages via gRPC publish to Redis Pub/Sub channels fanning out to all gateway nodes. Each gateway delivers to relevant connected clients. For durable delivery, messages persist to Kafka with consumer groups tracking acknowledgements. Channel sharding across Redis cluster slots enables horizontal scaling.
