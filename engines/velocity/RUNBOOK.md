# Runbook

## Startup
1. Verify Vault accessible and provider credentials loaded.
2. Start orchestrator: ./velocity-orchestrator --config config/velocity.toml.
3. Verify providers: GET /v1/providers shows all healthy.
4. Start Envoy sidecars with steering config.
5. Test purge: POST /v1/cache/purge.

## Monitoring
- Dashboard: cache hit rate by CDN, purge backlog, warming throughput.
- Alerts: purge backlog > 10K, provider error > 5%, hit rate drop > 10%.
