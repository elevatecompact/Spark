# search-service — Runbook
## Alerts: SearchLatencyP99 > 500ms, ESClusterStatusYellow, ESWriteRejected > 100/min, IndexWorkerLag > 10000, AutocompleteLatency > 200ms
## Cluster health: kubectl exec es-pod -- curl -s localhost:9200/_cluster/health
## Force merge: POST /v1/admin/force-merge — merges ES segments for better query performance.
## Clear cache: POST /v1/admin/clear-cache
## Reindex corrupt index: POST /v1/index/{type}/reindex
## Hot thread: GET /_nodes/hot_threads on ES node to identify slow queries.
