# Runbook

## Startup
1. Load blocklists from /data/blocklists/.
2. Start tier 1: ./sentinel-tier1 --config config/sentinel.toml.
3. Start tier 2/3 GPU nodes: ./sentinel-tier2 --gpu-device 0.
4. Deploy models: POST /v1/model/update.
5. Verify: curl -X POST /v1/moderate/text -d '{"content":"test"}'.

## Monitoring
- Dashboard: decisions/sec, tier distribution, latency, false positive rate.
- Alerts: false positive rate > 5%, tier 3 budget exceeded.
