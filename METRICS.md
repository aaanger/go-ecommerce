# Prometheus Metrics for Ecommerce Application

This document describes the Prometheus metrics implementation for the ecommerce application.

## Available Metrics

### HTTP Metrics
- `http_request_duration_seconds` - Duration of HTTP requests in seconds
- `http_requests_total` - Total number of HTTP requests
- `http_requests_in_flight` - Current number of HTTP requests being processed

### Database Metrics
- `database_query_duration_seconds` - Duration of database queries in seconds

### Redis Metrics
- `redis_operation_duration_seconds` - Duration of Redis operations in seconds

### Kafka Metrics
- `kafka_message_processing_duration_seconds` - Duration of Kafka message processing in seconds

### Business Metrics
- `cart_operations_total` - Total number of cart operations
- `order_operations_total` - Total number of order operations
- `product_operations_total` - Total number of product operations

## Usage

### Starting the Application

1. Start your ecommerce application:
```bash
go run cmd/main.go
```

2. Start Prometheus and Grafana:
```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

### Accessing Metrics

- **Application Metrics**: http://localhost:8080/metrics
- **Health Check**: http://localhost:8080/health
- **Prometheus UI**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

### Example Queries

#### HTTP Request Duration (95th percentile)
```
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, endpoint))
```

#### Request Rate by Endpoint
```
sum(rate(http_requests_total[5m])) by (endpoint)
```

#### Cart Operations Success Rate
```
sum(rate(cart_operations_total{status="success"}[5m])) / sum(rate(cart_operations_total[5m]))
```

#### Error Rate
```
sum(rate(http_requests_total{status_code=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

## Adding Custom Metrics

### In Handlers
```go
import "github.com/aaanger/ecommerce/pkg/middleware"

type MyHandler struct {
    metrics *middleware.BusinessMetricsMiddleware
}

func NewMyHandler() *MyHandler {
    return &MyHandler{
        metrics: middleware.NewBusinessMetricsMiddleware(),
    }
}

func (h *MyHandler) MyMethod(c *gin.Context) {
    // Your logic here
    if err != nil {
        h.metrics.RecordCartOperation("my_operation", "error")
        return
    }
    h.metrics.RecordCartOperation("my_operation", "success")
}
```

### Database Operations
```go
import "github.com/aaanger/ecommerce/pkg/middleware"

func (r *Repository) GetData() {
    defer middleware.DatabaseMetricsMiddleware("select", "users")()
    // Your database query here
}
```

### Redis Operations
```go
import "github.com/aaanger/ecommerce/pkg/middleware"

func (r *RedisClient) GetData() {
    defer middleware.RedisMetricsMiddleware("get", "user:*")()
    // Your Redis operation here
}
```

## Grafana Dashboards

After starting Grafana, you can create dashboards to visualize:

1. **HTTP Metrics Dashboard**
   - Request duration percentiles
   - Request rate by endpoint
   - Error rate
   - In-flight requests

2. **Business Metrics Dashboard**
   - Cart operations success rate
   - Order processing metrics
   - Product operation metrics

3. **Infrastructure Dashboard**
   - Database query performance
   - Redis operation latency
   - Kafka message processing time

## Alerting

You can set up alerts in Prometheus for:

- High error rates (>5%)
- Slow response times (>2s for 95th percentile)
- High database query duration
- Cart operation failures

Example alert rule:
```yaml
groups:
  - name: ecommerce_alerts
    rules:
      - alert: HighErrorRate
        expr: sum(rate(http_requests_total{status_code=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}"
``` 