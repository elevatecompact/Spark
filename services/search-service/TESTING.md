# search-service — Testing Guide
## Unit: Query parsing and validation, filter building, relevance scoring logic, synonym expansion, autocomplete prefix matching.
## Integration: Index→search→verify result, update index→re-search, delete from index, filter combinations, pagination correctness.
## Relevance: Precision@10, Recall@20, NDCG@10 measured against labeled test set. Query latency p50/p95/p99 monitoring.
## Load: k6 scripts in tests/load/ — search-throughput.js, autocomplete-burst.js. Target: 1000 QPS with p99 < 300ms.
