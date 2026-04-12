package middlewares

import (
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	"github.com/gin-gonic/gin"
)

// BffHeadersMiddleware returns a Gin middleware that enables bidirectional
// header propagation between the incoming HTTP request and upstream BFF calls.
//
// cfg is the BffMiddlewareConfig from the BffClientConfig used to create the adapter
// for the single upstream this BFF serves.
//
// On the way in (before c.Next()):
//   - Only headers listed in cfg.RequestHeaders are stored in the context via
//     WithRequestHeaders, keeping the context map minimal.
//   - A mutable response-header bag is installed via InitResponseHeaders so BFF
//     adapters can accumulate upstream response headers automatically after each response.
//
// On the way out (after c.Next()):
//   - All headers accumulated in the bag are written to the outgoing response.
//
// Usage:
//
//	router.Use(bffclient.BffHeadersMiddleware(bffCfg.HeadersProxy))
func BffHeadersMiddleware(cfg contracts.BffMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := WithRequestHeaders(c, cfg.RequestHeaders)
		ctx = InitResponseHeaders(ctx, cfg.ResponseHeaders)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if len(c.Errors) == 0 {
			ResponseHeadersFromCtx(c)
		}
	}
}
