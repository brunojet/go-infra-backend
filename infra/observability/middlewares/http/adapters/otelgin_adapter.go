package adapters

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// OtelGinMiddleware returns the otelgin middleware configured with the provided service name.
// It relies on the global TracerProvider/Propagators to be configured in bootstrap.
func OtelGinMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
