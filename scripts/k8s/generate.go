//go:build ignore
// +build ignore

// Generator for per-service Kubernetes manifests. The Go template avoids
// PowerShell quoting headaches when emitting many similar YAML files.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Service struct {
	Name         string
	Port         string
	Database     string
	Replicas     int
	CPU          string
	Memory       string
	CPULimit     string
	MemoryLimit  string
	NeedsClickHouse bool
}

var services = []Service{
	{Name: "advertising-service", Port: "4001", Database: "spark_advertising", Replicas: 2, CPU: "150m", Memory: "192Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "analytics-service", Port: "4002", Database: "spark_analytics", Replicas: 2, CPU: "300m", Memory: "384Mi", CPULimit: "1000m", MemoryLimit: "1Gi", NeedsClickHouse: true},
	{Name: "commerce-service", Port: "4003", Database: "spark_commerce", Replicas: 2, CPU: "150m", Memory: "192Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "community-service", Port: "4004", Database: "spark_community", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "competition-service", Port: "4007", Database: "spark_competition", Replicas: 2, CPU: "150m", Memory: "192Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "creator-service", Port: "4008", Database: "spark_creator", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "discovery-service", Port: "4010", Database: "spark_discovery", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "event-service", Port: "4011", Database: "spark_event", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "gift-service", Port: "4015", Database: "spark_gift", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "identity-service", Port: "4016", Database: "spark_identity", Replicas: 3, CPU: "250m", Memory: "256Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "licensing-service", Port: "4017", Database: "spark_licensing", Replicas: 2, CPU: "150m", Memory: "192Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "media-service", Port: "4018", Database: "spark_media", Replicas: 3, CPU: "400m", Memory: "512Mi", CPULimit: "1500m", MemoryLimit: "2Gi"},
	{Name: "messaging-service", Port: "4019", Database: "spark_messaging", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "moderation-service", Port: "4020", Database: "spark_moderation", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "notification-service", Port: "4021", Database: "spark_notification", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "payment-service", Port: "4022", Database: "spark_payment", Replicas: 3, CPU: "250m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "recommendation-service", Port: "4023", Database: "spark_recommendation", Replicas: 2, CPU: "400m", Memory: "512Mi", CPULimit: "1500m", MemoryLimit: "2Gi"},
	{Name: "search-service", Port: "4024", Database: "spark_search", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "stream-service", Port: "4025", Database: "spark_stream", Replicas: 3, CPU: "500m", Memory: "512Mi", CPULimit: "1500m", MemoryLimit: "2Gi"},
	{Name: "subscription-service", Port: "4026", Database: "spark_subscription", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "translation-service", Port: "4027", Database: "spark_translation", Replicas: 2, CPU: "150m", Memory: "192Mi", CPULimit: "500m", MemoryLimit: "512Mi"},
	{Name: "trust-service", Port: "4028", Database: "spark_trust", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "viewer-service", Port: "4029", Database: "spark_viewer", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
	{Name: "wallet-service", Port: "4030", Database: "spark_wallet", Replicas: 2, CPU: "200m", Memory: "256Mi", CPULimit: "750m", MemoryLimit: "768Mi"},
}

const tmpl = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: spark
  labels:
    app: {{ .Name }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{ .Name }}
  template:
    metadata:
      labels:
        app: {{ .Name }}
    spec:
      containers:
        - name: {{ .Name }}
          image: ghcr.io/elevatecompact/{{ .Name }}:latest
          ports:
            - containerPort: {{ .Port }}
          env:
            - name: DATABASE_URL
              value: "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/{{ .Database }}?sslmode=disable"
            - name: REDIS_URL
              value: "redis://$(REDIS_HOST):$(REDIS_PORT)/0"
            - name: KAFKA_BROKERS
              valueFrom:
                configMapKeyRef:
                  name: spark-config
                  key: KAFKA_BROKERS
            - name: LOG_LEVEL
              valueFrom:
                configMapKeyRef:
                  name: spark-config
                  key: LOG_LEVEL
            {{- if .NeedsClickHouse }}
            - name: ANALYTICS_CLICKHOUSE_URL
              value: "clickhouse://$(CLICKHOUSE_USER):$(CLICKHOUSE_PASSWORD)@clickhouse.spark.svc.cluster.local:9000/spark"
            {{- end }}
          envFrom:
            - secretRef:
                name: spark-secrets
          livenessProbe:
            httpGet:
              path: /health
              port: {{ .Port }}
            initialDelaySeconds: 10
            periodSeconds: 15
          readinessProbe:
            httpGet:
              path: /ready
              port: {{ .Port }}
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            requests:
              cpu: "{{ .CPU }}"
              memory: "{{ .Memory }}"
            limits:
              cpu: "{{ .CPULimit }}"
              memory: "{{ .MemoryLimit }}"
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}
  namespace: spark
spec:
  ports:
    - name: http
      port: {{ .Port }}
      targetPort: {{ .Port }}
  selector:
    app: {{ .Name }}
`

func main() {
	out := flag.String("out", "infrastructure/kubernetes/services", "output directory")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	t, err := template.New("svc").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	for _, svc := range services {
		f, err := os.Create(filepath.Join(*out, svc.Name+".yaml"))
		if err != nil {
			log.Fatal(err)
		}
		if err := t.Execute(f, svc); err != nil {
			f.Close()
			log.Fatal(err)
		}
		f.Close()
	}
	log.Printf("generated %d manifests in %s", len(services), *out)
}
