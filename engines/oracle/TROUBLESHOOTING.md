# Troubleshooting

## Empty recommendations
1. Verify user exists in Redis: redis-cli GET feature:user:{userId}.
2. Check Milvus collection is populated: milvus-cli count collection items.
3. Confirm model is deployed: GET /v1/model/active.
4. Check ranker diversity threshold.

## High inference latency
1. Reduce model.batch_size for latency-sensitive paths.
2. Verify ONNX uses CUDA: check logs for Execution Provider: CUDA.
3. Scale inference nodes behind load balancer.
4. Use Redis pipeline for multi-key fetches.
