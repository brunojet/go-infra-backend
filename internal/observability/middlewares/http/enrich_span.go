package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// EnrichAppSpan attaches application-specific attributes (like user/tenant)
// to the current span if one exists in the request context.
func EnrichAppSpan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		span := trace.SpanFromContext(ctx)
		if span != nil && span.IsRecording() {
			userID := c.GetHeader("X-User-ID")
			tenantID := c.GetHeader("X-Tenant-ID")
			if userID != "" {
				span.SetAttributes(attribute.String("app.user_id", userID))
			}
			if tenantID != "" {
				span.SetAttributes(attribute.String("app.tenant_id", tenantID))
			}
		}
		c.Next()
	}
}
