package bffclient

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
//   - A mutable response-header bag is installed via InitResponseCapture so BFF
//     adapters can accumulate upstream response headers via CaptureResponseHeader.
//
// On the way out (after c.Next()):
//   - All headers accumulated in the bag are written to the outgoing response.
//
// Usage:
//
//	router.Use(bffclient.BffHeadersMiddleware(bffCfg.HeadersProxy))
func BffHeadersMiddleware(cfg contracts.BffMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestHeaders := make(map[string]string, len(cfg.RequestHeaders))

		for _, key := range cfg.RequestHeaders {
			if value := c.Request.Header.Get(key); value != "" {
				requestHeaders[key] = value
			}
		}

		ctx := WithRequestHeaders(c.Request.Context(), requestHeaders)
		ctx = InitResponseCapture(ctx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		for k, v := range ResponseHeadersFromCtx(c.Request.Context()) {
			c.Header(k, v)
		}
	}
}
