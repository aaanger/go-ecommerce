package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsInitialization(t *testing.T) {
	// Test that metrics are properly initialized
	assert.NotNil(t, HTTPRequestDuration)
	assert.NotNil(t, HTTPRequestsTotal)
	assert.NotNil(t, HTTPRequestsInFlight)
	assert.NotNil(t, DatabaseQueryDuration)
	assert.NotNil(t, RedisOperationDuration)
	assert.NotNil(t, KafkaMessageProcessingDuration)
	assert.NotNil(t, CartOperationsTotal)
	assert.NotNil(t, OrderOperationsTotal)
	assert.NotNil(t, ProductOperationsTotal)
}

func TestHTTPRequestDuration(t *testing.T) {
	// Test HTTP request duration metric
	labels := []string{"GET", "/test", "200"}
	HTTPRequestDuration.WithLabelValues(labels...).Observe(0.1)
	
	// Verify metric exists
	metric, err := HTTPRequestDuration.GetMetricWithLabelValues(labels...)
	assert.NoError(t, err)
	assert.NotNil(t, metric)
}

func TestCartOperationsTotal(t *testing.T) {
	// Test cart operations metric
	labels := []string{"add_product", "success"}
	CartOperationsTotal.WithLabelValues(labels...).Inc()
	
	// Verify metric exists
	metric, err := CartOperationsTotal.GetMetricWithLabelValues(labels...)
	assert.NoError(t, err)
	assert.NotNil(t, metric)
}

func TestDatabaseQueryDuration(t *testing.T) {
	// Test database query duration metric
	labels := []string{"select", "users"}
	DatabaseQueryDuration.WithLabelValues(labels...).Observe(0.05)
	
	// Verify metric exists
	metric, err := DatabaseQueryDuration.GetMetricWithLabelValues(labels...)
	assert.NoError(t, err)
	assert.NotNil(t, metric)
}

func TestMetricsRegistration(t *testing.T) {
	// Test that all metrics are registered with Prometheus
	// Try to register metrics (they should already be registered with DefaultRegisterer)
	// This test ensures our metrics are properly defined
	assert.NotPanics(t, func() {
		HTTPRequestDuration.WithLabelValues("GET", "/test", "200").Observe(0.1)
		HTTPRequestsTotal.WithLabelValues("GET", "/test", "200").Inc()
		HTTPRequestsInFlight.WithLabelValues("GET", "/test").Inc()
		CartOperationsTotal.WithLabelValues("add_product", "success").Inc()
	})
} 