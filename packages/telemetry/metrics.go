package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requestDuration   *prometheus.HistogramVec
	requestsTotal     *prometheus.CounterVec
	activeConnections prometheus.Gauge
	errorsTotal       *prometheus.CounterVec
}

func InitMetrics(serviceName string) *Metrics {
	m := &Metrics{
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      "http_request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds",
				Namespace: serviceName,
				Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path", "status"},
		),
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
				Namespace: serviceName,
			},
			[]string{"method", "path", "status"},
		),
		activeConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:      "active_connections",
				Help:      "Number of active connections",
				Namespace: serviceName,
			},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:      "errors_total",
				Help:      "Total number of errors by type",
				Namespace: serviceName,
			},
			[]string{"type"},
		),
	}

	return m
}

func MetricsMiddleware(m *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			m.activeConnections.Inc()
			defer m.activeConnections.Dec()

			srw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(srw, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(srw.statusCode)

			m.requestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
			m.requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		})
	}
}

func RecordRequest(m *Metrics, duration time.Duration, method, path, status string) {
	if m == nil {
		return
	}
	m.requestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
	m.requestsTotal.WithLabelValues(method, path, status).Inc()
}

func RecordActiveConnections(m *Metrics, delta float64) {
	if m == nil {
		return
	}
	m.activeConnections.Add(delta)
}

func RecordError(m *Metrics, errorType string) {
	if m == nil {
		return
	}
	m.errorsTotal.WithLabelValues(errorType).Inc()
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *statusResponseWriter) WriteHeader(code int) {
	srw.statusCode = code
	srw.ResponseWriter.WriteHeader(code)
}
