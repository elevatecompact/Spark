# Troubleshooting

## Purge not taking effect
1. Check purge job status: GET /v1/cache/purge/status/:purgeId.
2. Verify provider propagation time (CloudFront 5-10s).
3. Check if cached behind origin shield with its own TTL.
4. Ensure token has purge:write scope.

## Failover not triggering
1. Check failover_threshold vs current error rate.
2. Verify edge probes running: GET /v1/steering/probes.
3. Check backup provider marked healthy: GET /v1/providers.
4. Verify Envoy has latest routing rules.
