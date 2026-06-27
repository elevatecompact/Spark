# Runbook

## Startup
1. Start etcd cluster: etcd --config-file etcd.yaml on 3+ nodes.
2. Start registrar: ./atlas-registrar --config config/atlas.toml.
3. Start agent on each node: ./atlas-agent --registrar-addr 10.0.0.1:8600.
4. Verify: curl http://localhost:8600/v1/health.

## Monitoring
- Dashboard: registration rate, cache hit ratio, propagation latency.
- Alerts: registrar leader election, agent disconnection, propagation latency > 1s.
- Logs: structured JSON with fields serviceId, event, latency.
