# Runbook

## Startup
1. Verify Redis cluster: redis-cli --cluster check 127.0.0.1:6379.
2. Start Milvus: milvus run standalone.
3. Start Oracle inference: ./oracle-server --config config/oracle.toml.
4. Deploy model: curl -X POST /v1/model/deploy -d '{"modelId":"v3.2"}'.
5. Verify: curl -X POST /v1/recommend -d '{"userId":"test","count":5}'.

## Monitoring
- Dashboard: QPS, P99 latency, cache hit rate, model score distribution.
- Alerts: latency > 200ms P99, error rate > 1%, model staleness.
