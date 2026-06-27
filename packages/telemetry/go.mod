module github.com/elevatecompact/spark/packages/telemetry

go 1.22

require (
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
	go.opentelemetry.io/otel/trace v1.24.0
	github.com/prometheus/client_golang v1.18.0
)
