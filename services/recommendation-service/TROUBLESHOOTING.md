# recommendation-service — Troubleshooting
## Empty feeds or poor recommendations: Embedding computation failed, feature store stale, model returned all-zero scores. Check embedding tables, verify feature freshness, run offline model eval.
## High inference latency: GPU OOM, model too large for ONNX batch size, Redis cache miss storm. Scale inference replicas, optimize ONNX graph, pre-warm cache.
## CTR dropping: Model drift, feature distribution shift, cold start failure for new content. Rollback to previous model version, retrain on recent data, investigate feature distributions.
