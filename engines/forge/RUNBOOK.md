# Runbook

## Startup
1. Check GPU: nvidia-smi or radeontop.
2. Verify Redis queue: redis-cli ping.
3. Start Forge worker: ./forge-worker --config config/forge.toml.
4. Start API gateway: ./forge-api --queue-addr localhost:6379.
5. Submit test job.

## Monitoring
- Dashboard: active jobs, GPU utilization, queue depth, encoding FPS.
- Alerts: GPU memory > 90%, job failure rate > 2%, queue backlog > 500.
