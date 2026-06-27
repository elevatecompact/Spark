tempo:
  metricsGenerator:
    enabled: true
    processor:
      serviceGraphs:
        dimensions:
          - namespace
      spanMetrics:
        dimensions:
          - namespace
          - service
          - span_kind
          - status_code

  storage:
    trace:
      backend: local
      wal:
        path: /var/tempo/wal
      local:
        path: /var/tempo/blocks

  retention:
    max_block_duration: 2h
    min_block_duration: 30m

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: "0.0.0.0:4317"
        http:
          endpoint: "0.0.0.0:4318"

ingester:
  trace_idle_period: 10s
  max_block_duration: 30m
  max_block_bytes: 524288000

compactor:
  compaction:
    block_retention: ${retention_days * 24}h

querier:
  max_concurrent_queries: 4

queryFrontend:
  max_outstanding_per_tenant: 100

metricsGenerator:
  enabled: true

overrides:
  defaults:
    metricsGenerator:
      processors:
        - service-graphs
        - span-metrics

storage:
  trace:
    backend: s3
    s3:
      bucket: spark-tempo-traces
      endpoint: s3.amazonaws.com

ingester:
  lifecycler:
    ring:
      replication_factor: ${environment == "production" ? 3 : 1}

memberlist:
  abort_if_cluster_join_fails: false
