# Scaling

Echo scales by adding gateway nodes. Each node is stateless regarding connections. Redis Pub/Sub handles cross-node delivery. For global deployments, deploy gateway nodes per region with regional Redis cluster; cross-region messages use Kafka mirroring. Channel affinity ensures high-traffic channels are not over-distributed. Use client IP hash or session ID hash for consistent routing.
