package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// MetricsHandler handles the metrics endpoint
type MetricsHandler struct{}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// GetMetrics returns the Prometheus metrics endpoint
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	// Use promhttp.HandlerFor to serve metrics
	promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

// RegisterMetricsRoutes registers the metrics routes
func RegisterMetricsRoutes(router *gin.Engine) {
	handler := NewMetricsHandler()
	
	// Metrics endpoint
	router.GET("/metrics", handler.GetMetrics)
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
} 