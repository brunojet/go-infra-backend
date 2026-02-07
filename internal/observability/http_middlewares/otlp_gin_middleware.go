package middlewares

import (
	"github.com/brunojet/go-infra-backend/internal/config"
	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// OtelGinMiddleware returns a Gin middleware that starts spans for incoming requests.
// The service name can be configured via OTEL_SERVICE_NAME env var. Default keeps
// the previous hardcoded service name.
func OtelGinMiddleware() gin.HandlerFunc {
	svc := config.Get("OTEL_SERVICE_NAME", "github.com/brunojet/go-infra-backend")
	return otelgin.Middleware(svc)
}
