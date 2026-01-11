package ginmw

import (
	"strings"

	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
	"github.com/gin-gonic/gin"
)

// RequestID ensures every request has a correlation id.
//
// It reads from X-Request-ID and if absent generates one. It also sets the
// response header and injects it into request context so it is accessible via
// contracts.Context.RequestContext().
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := strings.TrimSpace(c.GetHeader(requestid.Header))
		if id == "" {
			id = requestid.New()
		}

		c.Header(requestid.Header, id)
		c.Request = c.Request.WithContext(requestid.With(c.Request.Context(), id))
		c.Next()
	}
}
