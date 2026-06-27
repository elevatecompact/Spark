# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| backend.type | etcd | Backend store: etcd or consul |
| backend.endpoints | localhost:2379 | Comma-separated backend endpoints |
| agent.port | 8600 | Agent HTTP API port |
| agent.health_check_interval | 10 | Health check interval in seconds |
| agent.cache_ttl | 30 | Local cache TTL in seconds |
| registrar.replication_factor | 3 | Number of registrar nodes |
| routing.strategy | weighted_random | Routing strategy |
| routing.blue_green.enabled | false | Enable blue-green routing |
