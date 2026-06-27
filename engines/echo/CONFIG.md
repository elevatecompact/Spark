# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| gateway.port | 8080 | WebSocket listener port |
| gateway.max_connections | 100000 | Max concurrent connections per node |
| gateway.heartbeat_interval | 30 | Ping interval in seconds |
| channel.max_subscribers | 10000 | Max subscribers per channel |
| channel.message_size_limit | 256KB | Max message payload size |
| redis.pubsub.endpoints | localhost:6379 | Redis Pub/Endpoints |
| kafka.brokers | localhost:9092 | Kafka brokers for durable delivery |
| durable.enabled | false | Enable Kafka-backed durable delivery |
