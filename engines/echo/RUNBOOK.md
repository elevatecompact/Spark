# Runbook

## Startup
1. Verify Redis: redis-cli ping.
2. Start Echo gateway: ./echo-gateway --config config/echo.toml.
3. Verify: curl http://localhost:8080/v1/health.
4. Enable durable mode with durable.enabled = true.

## Monitoring
- Dashboard: active connections, message throughput, delivery latency.
- Alerts: connection drop > 50%, message latency > 200ms, Redis Pub/Sub errors.
