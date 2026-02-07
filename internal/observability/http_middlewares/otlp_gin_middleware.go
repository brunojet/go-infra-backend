package middlewares

import (
	"github.com/brunojet/go-infra-backend/internal/config"
	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	// OTELServiceNameEnv is the env var key for overriding the OTel service name
	otelServiceNameEnv = "OTEL_SERVICE_NAME"
	// defaultServiceName is the fallback service name used when env var is not set
	defaultServiceName = "github.com/brunojet/go-infra-backend"
)

// OtelGinMiddleware returns a Gin middleware that starts spans for incoming requests.
// The service name can be configured via OTEL_SERVICE_NAME env var. Default keeps
// the previous hardcoded service name.
func OtelGinMiddleware() gin.HandlerFunc {
	svc := config.GetEnv(otelServiceNameEnv, defaultServiceName)
	return otelgin.Middleware(svc)
}
