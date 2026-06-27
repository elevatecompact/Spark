# chat-service — Troubleshooting
## Messages not delivering: WS disconnected (heartbeat timeout), Redis pub/sub channel down, user muted. Check WS state, PUBSUB CHANNELS, mute key in Redis.
## High memory: Too many room subscriptions per node, message buffer not flushing. Reduce rooms/node, increase WS_MESSAGE_BUFFER_SIZE.
## History missing: Retention expired, partition not queried. Check message age vs retention config, verify PG partition pruning.
