// Package metrics 提供應用程式的指標收集功能，使用 Prometheus 作為指標收集和暴露的工具。
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector 指標收集器
type Collector struct {
	// HTTP 請求指標

	// HTTPRequestTotal 記錄 HTTP 請求總數，標籤包括 method、path、status、api_version
	HTTPRequestTotal *prometheus.CounterVec
	// HTTPRequestDuration 記錄 HTTP 請求持續時間，標籤包括 method、path、status
	HTTPRequestDuration *prometheus.HistogramVec

	// 業務指標
	// BusinessEvents 記錄業務事件，標籤包括 event_type、domain、status
	BusinessEvents *prometheus.CounterVec

	// 系統指標
	// GoroutineCount 記錄當前 Goroutine 數量
	GoroutineCount prometheus.Gauge
	// MemoryUsage 記錄記憶體使用量
	MemoryUsage prometheus.Gauge
	// CPUUsage 記錄 CPU 使用率
	CPUUsage prometheus.Gauge
}

// NewMetricsCollector 建立新的 MetricsCollector 實例
func NewMetricsCollector(namespace string) *Collector {
	collector := &Collector{
		HTTPRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "api_version"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		BusinessEvents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "business_events_total",
				Help:      "Total number of business events",
			},
			[]string{"event_type", "domain", "status"},
		),
		GoroutineCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "go_goroutines",
				Help:      "Number of goroutines",
			},
		),
		MemoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "go_memory_usage_bytes",
				Help:      "Memory usage in bytes",
			},
		),
		CPUUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "go_cpu_usage_percent",
				Help:      "CPU usage in percent",
			},
		),
	}

	prometheus.MustRegister(collector.HTTPRequestTotal)
	prometheus.MustRegister(collector.HTTPRequestDuration)
	prometheus.MustRegister(collector.BusinessEvents)
	prometheus.MustRegister(collector.GoroutineCount)
	prometheus.MustRegister(collector.MemoryUsage)
	prometheus.MustRegister(collector.CPUUsage)
	return collector
}

// RecordHTTPRequest 記錄 HTTP 請求指標
func (m *Collector) RecordHTTPRequest(method, path, status, apiVersion string, duration float64) {
	m.HTTPRequestTotal.WithLabelValues(method, path, status, apiVersion).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
}

// RecordBusinessEvent 記錄業務事件
func (m *Collector) RecordBusinessEvent(eventType, domain, status string) {
	m.BusinessEvents.WithLabelValues(eventType, domain, status).Inc()
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
