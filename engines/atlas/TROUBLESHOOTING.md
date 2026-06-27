# Troubleshooting

## Service not discoverable
1. Confirm the service TTL heartbeat is arriving - check agent logs for heartbeat.received.
2. Verify metadata tags match the query filter.
3. Check if the service was deregistered: look for service.deregistered event.
4. Force agent cache refresh: POST /v1/discover/:serviceName/refresh.

## High lookup latency
1. Check registrar CPU - add registrar nodes if > 70%.
2. Verify etcd cluster health: etcdctl endpoint health.
3. Reduce agent.cache_ttl if cache miss rate exceeds 5%.
4. Ensure agent can reach registrar on port 8600 without network latency.
