package bffclient

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/infra/bffclient/types"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	"github.com/gin-gonic/gin"
)

func WithRequestHeaders(c *gin.Context, requestHeaders []string) context.Context {
	requestHeadersMap := make(map[string]string, len(requestHeaders))
	for _, key := range requestHeaders {
		if value := c.Request.Header.Get(key); len(value) > 0 {
			requestHeadersMap[key] = value
		}
	}

	return context.WithValue(c.Request.Context(), types.BffRequestHeadersKey{}, requestHeadersMap)
}

func InitResponseCapture(ctx context.Context, responseHeaders []string) context.Context {
	return context.WithValue(ctx, types.BffResponseHeadersKey{}, types.NewResponseHeadersBag(responseHeaders))
}

func ResponseHeadersFromCtx(c *gin.Context) {
	if bag, ok := c.Request.Context().Value(types.BffResponseHeadersKey{}).(*types.ResponseHeadersBag); ok {
		responseHeadersMap := bag.All()
		for k, v := range responseHeadersMap {
			if len(v) > 0 {
				c.Header(k, v)
			}
		}
	}
}

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
		ctx = InitResponseCapture(ctx, cfg.ResponseHeaders)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if len(c.Errors) == 0 {
			ResponseHeadersFromCtx(c)
		}
	}
}
