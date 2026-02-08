package httpmiddlewares

import (
	internalmiddlewares "github.com/brunojet/go-infra-backend/internal/observability/http_middlewares"
	"github.com/gin-gonic/gin"
)

// OtelGinMiddleware returns a Gin middleware that starts spans for incoming requests.
// The service name can be configured via OTEL_SERVICE_NAME env var.
func OtelGinMiddleware() gin.HandlerFunc {
	return internalmiddlewares.OtelGinMiddleware()
}
