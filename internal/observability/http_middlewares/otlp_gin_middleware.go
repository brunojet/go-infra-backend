package middlewares

import (
	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// OtelGinMiddleware returns a Gin middleware that starts spans for incoming requests.
// It uses the service name from providers.GetServiceName().
func OtelGinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("github.com/brunojet/go-infra-backend")
}
