# Runbook

## Startup
1. Verify config/pulse.toml exists and ports are available.
2. Start Redis: redis-server --port 6379.
3. Start control plane: ./pulse-control --config config/pulse.toml.
4. Start SFU nodes: ./pulse-sfu --control-addr 127.0.0.1:9000.
5. Verify health: curl http://localhost:8080/health.

## Monitoring
- Prometheus metrics at /metrics: active streams, viewer count, bitrate distribution.
- Critical alerts: SFU node offline, ingest bitrate drop > 50%, viewer disconnect spike.
- Logs are structured JSON to stdout; collect with filebeat.
