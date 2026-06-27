# Troubleshooting

## Client cannot connect
1. Check JWT token validity and expiry.
2. Verify WebSocket URL protocol (ws vs wss).
3. Check gateway connection count: GET /v1/health.
4. Look for session.connected event.

## Messages not delivered
1. Verify client is subscribed to correct channel.
2. Check message size under channel.message_size_limit.
3. If durable mode, verify Kafka consumer group is not lagging.
4. Test via REST: POST /v1/channel/test/publish.
