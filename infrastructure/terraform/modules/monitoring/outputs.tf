output "namespace" {
  description = "Monitoring namespace"
  value       = local.namespace
}

output "prometheus_service" {
  description = "Prometheus service endpoint"
  value       = "kube-prometheus-stack-prometheus.${local.namespace}.svc:9090"
}

output "grafana_service" {
  description = "Grafana service endpoint"
  value       = "kube-prometheus-stack-grafana.${local.namespace}.svc:80"
}

output "loki_service" {
  description = "Loki service endpoint"
  value       = "loki.${local.namespace}.svc:3100"
}

output "tempo_service" {
  description = "Tempo service endpoint"
  value       = "tempo.${local.namespace}.svc:3200"
}

output "alertmanager_service" {
  description = "AlertManager service endpoint"
  value       = "kube-prometheus-stack-alertmanager.${local.namespace}.svc:9093"
}
