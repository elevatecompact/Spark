loki:
  commonConfig:
    replication_factor: 1
  storage:
    type: filesystem
  schemaConfig:
    configs:
      - from: "2024-01-01"
        store: tsdb
        object_store: filesystem
        schema: v13
        index:
          prefix: index_
          period: 24h
  ingester:
    chunk_encoding: snappy
  querier:
    max_concurrent: 4
  compactor:
    retention_enabled: true
  limits_config:
    retention_period: ${retention_days}d
    max_query_lookback: ${retention_days}d
    max_line_size: 1MB
    reject_old_samples: true
    reject_old_samples_max_age: 168h

write:
  replicas: ${environment == "production" ? 3 : 1}
  persistence:
    size: ${environment == "production" ? "50Gi" : "10Gi"}

read:
  replicas: ${environment == "production" ? 3 : 1}
  persistence:
    size: ${environment == "production" ? "50Gi" : "10Gi"}

backend:
  replicas: ${environment == "production" ? 3 : 1}
  persistence:
    size: ${environment == "production" ? "20Gi" : "5Gi"}

gateway:
  enabled: true
  replicas: ${environment == "production" ? 2 : 1}

monitoring:
  dashboards:
    enabled: true
  rules:
    enabled: true
  alerts:
    enabled: true

test:
  enabled: false
