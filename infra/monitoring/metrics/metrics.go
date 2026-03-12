package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector 指標收集器
type MetricsCollector struct {
	// HTTP 請求指標

	// HttpRequestTotal 記錄 HTTP 請求總數，標籤包括 method、path、status、api_version
	HttpRequestTotal    *prometheus.CounterVec
	// HttpRequestDuration 記錄 HTTP 請求持續時間，標籤包括 method、path、status
	HttpRequestDuration *prometheus.HistogramVec
	// HttpRequestSize 記錄 HTTP 請求大小，標籤包括 method、path、status
	HttpRequestSize     *prometheus.SummaryVec
	// HttpResponseSize 記錄 HTTP 回應大小，標籤包括 method、path、status
	HttpResponseSize    *prometheus.SummaryVec

	// 業務指標
	// BusinessEvents 記錄業務事件，標籤包括 event_type、domain、status
	BusinessEvents         *prometheus.CounterVec
	// DatabaseQueryDuration 記錄資料庫查詢持續時間，標籤包括 operation、table、success
	DatabaseQueryDuration  *prometheus.HistogramVec
	// CacheOperationDuration 記錄快取操作持續時間，標籤包括 operation、cache_type、success
	CacheOperationDuration *prometheus.HistogramVec

	// 系統指標
	// GoroutineCount 記錄當前 Goroutine 數量
	GoroutineCount prometheus.Gauge
	// MemoryUsage 記錄記憶體使用量
	MemoryUsage    prometheus.Gauge
	// CPUUsage 記錄 CPU 使用率
	CPUUsage       prometheus.Gauge
}

// NewMetricsCollector 建立新的 MetricsCollector 實例
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		HttpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "api_version"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		BusinessEvents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_events_total",
				Help: "Total number of business events",
			},
			[]string{"event_type", "domain", "status"},
		),
		DatabaseQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"operation", "table", "success"},
		),
		GoroutineCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "go_goroutines",
				Help: "Number of goroutines",
			},
		),
		MemoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "go_memory_usage_bytes",
				Help: "Memory usage in bytes",
			},
		),
	}
}

// RecordHttpRequest 記錄 HTTP 請求指標
func (m *MetricsCollector) RecordHttpRequest(method, path, status, apiVersion string, duration float64) {
	m.HttpRequestTotal.WithLabelValues(method, path, status, apiVersion).Inc()
	m.HttpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
}

// RecordBusinessEvent 記錄業務事件
func (m *MetricsCollector) RecordBusinessEvent(eventType, domain, status string) {
	m.BusinessEvents.WithLabelValues(eventType, domain, status).Inc()
}

// RecordDatabaseQuery 記錄資料庫查詢
func (m *MetricsCollector) RecordDatabaseQuery(operation, table string, success bool, duration float64) {
	m.DatabaseQueryDuration.WithLabelValues(operation, table, boolToString(success)).Observe(duration)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
