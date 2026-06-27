locals {
  name_prefix = "${var.project_name}-${var.environment}"
  namespace   = "monitoring"
}

resource "kubernetes_namespace_v1" "monitoring" {
  metadata {
    name = local.namespace
    labels = {
      name                           = local.namespace
      "app.kubernetes.io/managed-by" = "Terraform"
      environment                    = var.environment
    }
  }
}

resource "helm_release" "kube_prometheus_stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  namespace  = local.namespace
  version    = "61.3.0"

  values = [
    templatefile("${path.module}/templates/prometheus-values.yaml.tpl", {
      environment             = var.environment
      project_name            = var.project_name
      grafana_admin_password  = var.grafana_admin_password
      alertmanager_slack_url  = var.alertmanager_slack_webhook
      alertmanager_pagerduty  = var.alertmanager_pagerduty_key
    })
  ]

  depends_on = [kubernetes_namespace_v1.monitoring]
}

resource "helm_release" "loki" {
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  namespace  = local.namespace
  version    = "6.6.4"

  values = [
    templatefile("${path.module}/templates/loki-values.yaml.tpl", {
      retention_days = var.loki_retention_days
      environment    = var.environment
    })
  ]

  depends_on = [kubernetes_namespace_v1.monitoring]
}

resource "helm_release" "tempo" {
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "tempo"
  namespace  = local.namespace
  version    = "1.10.3"

  values = [
    templatefile("${path.module}/templates/tempo-values.yaml.tpl", {
      retention_days = var.tempo_retention_days
      environment    = var.environment
    })
  ]

  depends_on = [kubernetes_namespace_v1.monitoring]
}

resource "kubernetes_config_map_v1" "grafana_dashboards" {
  metadata {
    name      = "spark-grafana-dashboards"
    namespace = local.namespace
    labels = {
      "grafana_dashboard" = "1"
      environment         = var.environment
    }
  }

  data = {
    "spark-platform-dashboard.json" = jsonencode({
      title       = "Spark Platform Overview"
      description = "High-level overview of the Spark platform"
      uid         = "spark-platform-overview"
      tags        = ["spark", "platform", var.environment]

      time = {
        from = "now-6h"
        to   = "now"
      }

      timezone    = "browser"
      schemaVersion = 36
      version     = 1

      panels = [
        {
          title      = "API Request Rate"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 0, y = 0 }
          targets = [
            {
              expr       = "sum(rate(http_requests_total{namespace=\"spark\"}[5m]))"
              legendFormat = "Requests/sec"
            }
          ]
        },
        {
          title      = "P99 Latency"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 12, y = 0 }
          targets = [
            {
              expr       = "histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{namespace=\"spark\"}[5m])) by (le))"
              legendFormat = "P99 Latency"
            }
          ]
        },
        {
          title      = "Error Rate"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 0, y = 8 }
          targets = [
            {
              expr       = "sum(rate(http_requests_total{namespace=\"spark\", status=~\"5..\"}[5m])) / sum(rate(http_requests_total{namespace=\"spark\"}[5m])) * 100"
              legendFormat = "Error Rate %"
            }
          ]
        },
        {
          title      = "CPU Usage"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 12, y = 8 }
          targets = [
            {
              expr       = "sum(rate(container_cpu_usage_seconds_total{namespace=\"spark\"}[5m]))"
              legendFormat = "CPU Cores"
            }
          ]
        },
        {
          title      = "Memory Usage"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 0, y = 16 }
          targets = [
            {
              expr       = "sum(container_memory_working_set_bytes{namespace=\"spark\"})"
              legendFormat = "Memory Bytes"
            }
          ]
        },
        {
          title      = "Active Connections"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 12, y = 16 }
          targets = [
            {
              expr       = "sum(nginx_connections_active{namespace=\"spark\"})"
              legendFormat = "Active Connections"
            }
          ]
        }
      ]
    })

    "spark-streaming-dashboard.json" = jsonencode({
      title       = "Spark Streaming Metrics"
      description = "Real-time streaming metrics for the Spark platform"
      uid         = "spark-streaming-metrics"
      tags        = ["spark", "streaming", var.environment]

      time = {
        from = "now-1h"
        to   = "now"
      }

      timezone      = "browser"
      schemaVersion = 36
      version       = 1

      panels = [
        {
          title      = "Streaming Ingest Rate"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 0, y = 0 }
          targets = [
            {
              expr       = "sum(rate(spark_streaming_events_total[5m]))"
              legendFormat = "Events/sec"
            }
          ]
        },
        {
          title      = "Streaming Lag"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 12, y = 0 }
          targets = [
            {
              expr       = "sum(spark_streaming_lag_total)"
              legendFormat = "Total Lag"
            }
          ]
        }
      ]
    })

    "spark-database-dashboard.json" = jsonencode({
      title       = "Spark Database Metrics"
      description = "Database performance metrics for the Spark platform"
      uid         = "spark-database-metrics"
      tags        = ["spark", "database", var.environment]

      time = {
        from = "now-6h"
        to   = "now"
      }

      timezone      = "browser"
      schemaVersion = 36
      version       = 1

      panels = [
        {
          title      = "Active Connections"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 0, y = 0 }
          targets = [
            {
              expr       = "sum(pg_stat_database_numbackends)"
              legendFormat = "Active Connections"
            }
          ]
        },
        {
          title      = "Transactions Per Second"
          type       = "graph"
          gridPos    = { h = 8, w = 12, x = 12, y = 0 }
          targets = [
            {
              expr       = "sum(rate(pg_stat_database_xact_commit[5m]))"
              legendFormat = "TPS"
            }
          ]
        }
      ]
    })
  }

  depends_on = [helm_release.kube_prometheus_stack]
}

resource "kubernetes_config_map_v1" "grafana_dashboard_datasources" {
  metadata {
    name      = "spark-grafana-datasources"
    namespace = local.namespace
    labels = {
      "grafana_datasource" = "1"
      environment          = var.environment
    }
  }

  data = {
    "datasources.yaml" = yamlencode({
      apiVersion = 1
      datasources = [
        {
          name      = "Prometheus"
          type      = "prometheus"
          access    = "proxy"
          url       = "http://kube-prometheus-stack-prometheus.${local.namespace}.svc:9090"
          isDefault = true
        },
        {
          name      = "Loki"
          type      = "loki"
          access    = "proxy"
          url       = "http://loki.${local.namespace}.svc:3100"
        },
        {
          name      = "Tempo"
          type      = "tempo"
          access    = "proxy"
          url       = "http://tempo.${local.namespace}.svc:3200"
        },
        {
          name      = "CloudWatch"
          type      = "cloudwatch"
          access    = "proxy"
          jsonData = {
            authType      = "default"
            defaultRegion = var.region
          }
        }
      ]
    })
  }

  depends_on = [helm_release.kube_prometheus_stack]
}

resource "kubernetes_manifest" "pod_monitor_spark" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "spark-applications"
      namespace = local.namespace
      labels = {
        release    = "kube-prometheus-stack"
        environment = var.environment
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/part-of" = "spark"
        }
      }
      namespaceSelector = {
        any = true
      }
      podMetricsEndpoints = [
        {
          port     = "metrics"
          interval = "15s"
          path     = "/metrics"
        }
      ]
    }
  }

  depends_on = [helm_release.kube_prometheus_stack]
}

resource "kubernetes_manifest" "service_monitor_spark" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    metadata = {
      name      = "spark-services"
      namespace = local.namespace
      labels = {
        release    = "kube-prometheus-stack"
        environment = var.environment
      }
    }
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/part-of" = "spark"
        }
      }
      namespaceSelector = {
        any = true
      }
      endpoints = [
        {
          port     = "metrics"
          interval = "15s"
          path     = "/metrics"
        }
      ]
    }
  }

  depends_on = [helm_release.kube_prometheus_stack]
}

resource "kubernetes_manifest" "prometheus_rule_spark" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PrometheusRule"
    metadata = {
      name      = "spark-alerts"
      namespace = local.namespace
      labels = {
        release    = "kube-prometheus-stack"
        environment = var.environment
      }
    }
    spec = {
      groups = [
        {
          name = "spark-platform"
          rules = [
            {
              alert = "SparkAPIHighErrorRate"
              expr  = "sum(rate(http_requests_total{namespace=\"spark\", status=~\"5..\"}[5m])) / sum(rate(http_requests_total{namespace=\"spark\"}[5m])) > 0.05"
              for   = "5m"
              annotations = {
                summary     = "Spark API error rate is above 5%"
                description = "Error rate is {{ $value | humanizePercentage }} for the last 5 minutes"
              }
              labels = {
                severity = "critical"
              }
            },
            {
              alert = "SparkAPILatencyHigh"
              expr  = "histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{namespace=\"spark\"}[5m])) by (le)) > 5"
              for   = "5m"
              annotations = {
                summary     = "Spark API P99 latency is above 5s"
                description = "P99 latency is {{ $value }}s for the last 5 minutes"
              }
              labels = {
                severity = "warning"
              }
            },
            {
              alert = "SparkStreamingLagHigh"
              expr  = "sum(spark_streaming_lag_total) > 10000"
              for   = "2m"
              annotations = {
                summary     = "Spark streaming lag is above 10,000 events"
                description = "Current lag is {{ $value }} events"
              }
              labels = {
                severity = "critical"
              }
            },
            {
              alert = "SparkLowDiskSpace"
              expr  = "node_filesystem_avail_bytes{mountpoint=\"/\", namespace=\"spark\"} / node_filesystem_size_bytes{mountpoint=\"/\", namespace=\"spark\"} < 0.1"
              for   = "5m"
              annotations = {
                summary     = "Spark node disk space is below 10%"
                description = "Available disk space is {{ $value | humanizePercentage }}"
              }
              labels = {
                severity = "critical"
              }
            },
            {
              alert = "SparkPodCrashLooping"
              expr  = "rate(kube_pod_container_status_restarts_total{namespace=\"spark\"}[30m]) > 1"
              for   = "5m"
              annotations = {
                summary     = "Spark pod is crash looping"
                description = "Pod {{ $labels.pod }} has restarted {{ $value }} times in 30 minutes"
              }
              labels = {
                severity = "warning"
              }
            },
            {
              alert = "SparkDatabaseConnectionPoolExhausted"
              expr  = "sum(pg_stat_database_numbackends) > 80"
              for   = "5m"
              annotations = {
                summary     = "Spark database connection pool is nearly exhausted"
                description = "Current active connections: {{ $value }}"
              }
              labels = {
                severity = "critical"
              }
            }
          ]
        }
      ]
    }
  }

  depends_on = [helm_release.kube_prometheus_stack]
}
