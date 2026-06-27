# search-service — Deployment Guide
K8s: k8s/search-service/ — api (3x 1GB), index-worker (2x 1GB consumers from Kafka), curator (CronJob for index maintenance).
Elasticsearch: 3-node cluster (hot nodes for indexing, warm nodes for search). 500GB SSD per node. Snapshot to S3 daily.
Deploy: kubectl apply -f k8s/search-service/, verify ES cluster health yellow/green. Run reindex if mappings changed.
Full reindex: POST /v1/index/reindex — reads from source DB and rebuilds all indices. Takes ~2h for full dataset.
