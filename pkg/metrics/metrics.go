package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// HTTPRequestsTotal tracks total number of HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// HTTPRequestsInFlight tracks current number of HTTP requests being processed
	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
		[]string{"method", "endpoint"},
	)

	// DatabaseQueryDuration tracks database query duration
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// RedisOperationDuration tracks Redis operation duration
	RedisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_operation_duration_seconds",
			Help:    "Duration of Redis operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "key_pattern"},
	)

	// KafkaMessageProcessingDuration tracks Kafka message processing duration
	KafkaMessageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_message_processing_duration_seconds",
			Help:    "Duration of Kafka message processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic", "status"},
	)

	// BusinessMetrics tracks business-specific metrics
	CartOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_operations_total",
			Help: "Total number of cart operations",
		},
		[]string{"operation", "status"},
	)

	OrderOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_operations_total",
			Help: "Total number of order operations",
		},
		[]string{"operation", "status"},
	)

	ProductOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)
) 