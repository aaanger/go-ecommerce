package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

	"github.com/aaanger/ecommerce/pkg/metrics"
)

// PrometheusMiddleware returns a gin middleware that tracks HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Track in-flight requests
		inflight := metrics.HTTPRequestsInFlight.WithLabelValues(method, path)
		inflight.Inc()
		defer inflight.Dec()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Record request duration
		metrics.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		// Record request count
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	}
}

// DatabaseMetricsMiddleware returns a middleware for tracking database operations
func DatabaseMetricsMiddleware(operation, table string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		metrics.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration)
	}
}

// RedisMetricsMiddleware returns a middleware for tracking Redis operations
func RedisMetricsMiddleware(operation, keyPattern string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		metrics.RedisOperationDuration.WithLabelValues(operation, keyPattern).Observe(duration)
	}
}

// KafkaMetricsMiddleware returns a middleware for tracking Kafka operations
func KafkaMetricsMiddleware(topic, status string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		metrics.KafkaMessageProcessingDuration.WithLabelValues(topic, status).Observe(duration)
	}
}

// BusinessMetricsMiddleware provides helper functions for business metrics
type BusinessMetricsMiddleware struct{}

// RecordCartOperation records cart operation metrics
func (b *BusinessMetricsMiddleware) RecordCartOperation(operation, status string) {
	metrics.CartOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordOrderOperation records order operation metrics
func (b *BusinessMetricsMiddleware) RecordOrderOperation(operation, status string) {
	metrics.OrderOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordProductOperation records product operation metrics
func (b *BusinessMetricsMiddleware) RecordProductOperation(operation, status string) {
	metrics.ProductOperationsTotal.WithLabelValues(operation, status).Inc()
}

// NewBusinessMetricsMiddleware creates a new business metrics middleware
func NewBusinessMetricsMiddleware() *BusinessMetricsMiddleware {
	return &BusinessMetricsMiddleware{}
}
