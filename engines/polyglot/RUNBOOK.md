# Runbook

## Startup
1. Verify model files in /models/.
2. Start server: ./polyglot serve --config config/polyglot.toml.
3. Verify: POST /v1/translate/detect -d '{"text":"Hello"}'.
4. Test: POST /v1/translate -d '{"text":"Hello","targetLang":"es"}'.
5. Start streaming workers if needed.

## Monitoring
- Dashboard: QPS, latency by pair, cache hit rate, GPU utilization.
- Alerts: latency > 1s P99, quality < 0.6, GPU memory > 90%.
