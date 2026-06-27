global:
  evaluationInterval: 30s
  scrapeInterval: 30s

alertmanager:
  enabled: true
  config:
    global:
      resolveTimeout: 5m
    route:
      groupBy: ["alertname", "severity", "cluster"]
      groupWait: 30s
      groupInterval: 5m
      repeatInterval: 4h
      receiver: "default"
      routes:
        - match:
            severity: critical
          receiver: "critical"
          repeatInterval: 1h
        - match:
            severity: warning
          receiver: "warning"
          repeatInterval: 4h
    receivers:
      - name: "default"
        slack_configs:
          - api_url: "${alertmanager_slack_url}"
            channel: "#spark-alerts"
            title: "[${environment}] Spark Alert"
            text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}\n{{ end }}'
            sendResolved: true
        pagerduty_configs:
          - routing_key: "${alertmanager_pagerduty}"
            severity: "{{ .Labels.severity }}"
            description: '{{ .Annotations.summary }}'
            details:
              environment: "${environment}"
              alert: '{{ .Annotations.description }}'
      - name: "critical"
        slack_configs:
          - api_url: "${alertmanager_slack_url}"
            channel: "#spark-critical"
            title: "[CRITICAL] [${environment}] {{ .GroupLabels.alertname }}"
            text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}\n{{ end }}'
            sendResolved: true
        pagerduty_configs:
          - routing_key: "${alertmanager_pagerduty}"
            severity: "critical"
            description: '{{ .Annotations.summary }}'
  ingress:
    enabled: false

grafana:
  enabled: true
  adminPassword: "${grafana_admin_password}"
  defaultDashboardsEnabled: true
  sidecar:
    dashboards:
      enabled: true
      label: grafana_dashboard
      searchNamespace: ALL
    datasources:
      enabled: true
      label: grafana_datasource
      searchNamespace: ALL
  datasources:
    datasources.yaml:
      apiVersion: 1
      datasources:
        - name: Prometheus
          type: prometheus
          access: proxy
          url: http://kube-prometheus-stack-prometheus.monitoring.svc:9090
          isDefault: true
        - name: Loki
          type: loki
          access: proxy
          url: http://loki.monitoring.svc:3100
        - name: Tempo
          type: tempo
          access: proxy
          url: http://tempo.monitoring.svc:3200
  ingress:
    enabled: false

kubeApiServer:
  enabled: true

kubeControllerManager:
  enabled: true

kubeDns:
  enabled: true

kubeEtcd:
  enabled: true

kubeProxy:
  enabled: true

kubeScheduler:
  enabled: true

kubeStateMetrics:
  enabled: true

kubelet:
  enabled: true

nodeExporter:
  enabled: true

prometheus:
  prometheusSpec:
    retention: ${environment == "production" ? "30d" : "7d"}
    retentionSize: ${environment == "production" ? "100GB" : "50GB"}
    podMonitorSelector:
      matchLabels:
        release: kube-prometheus-stack
    serviceMonitorSelector:
      matchLabels:
        release: kube-prometheus-stack
    ruleSelector:
      matchLabels:
        release: kube-prometheus-stack
    resources:
      requests:
        cpu: ${environment == "production" ? "500m" : "200m"}
        memory: ${environment == "production" ? "2Gi" : "1Gi"}
      limits:
        cpu: ${environment == "production" ? "2" : "1"}
        memory: ${environment == "production" ? "8Gi" : "4Gi"}
    storageSpec:
      volumeClaimTemplate:
        spec:
          storageClassName: gp3
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: ${environment == "production" ? "100Gi" : "20Gi"}

thanosRuler:
  enabled: false
