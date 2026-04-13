package middlewares

import (
	"context"
	"strings"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	"github.com/gin-gonic/gin"
)

type BffRequestHeadersKey struct{}

type BffResponseHeadersKey struct{}

// Middleware helpers section: these functions are used by the BffHeadersMiddleware and the netHttpAdapter to manage the flow of headers between incoming requests, context, and outgoing requests/responses.

func WithRequestHeaders(c *gin.Context, requestHeaders []string) context.Context {
	requestHeadersMap := make(map[string]string, len(requestHeaders))
	for _, key := range requestHeaders {
		if value := c.Request.Header.Get(key); len(value) > 0 {
			requestHeadersMap[key] = value
		}
	}

	return context.WithValue(c.Request.Context(), BffRequestHeadersKey{}, requestHeadersMap)
}

func InitResponseHeaders(ctx context.Context, responseHeaders []string) context.Context {
	return context.WithValue(ctx, BffResponseHeadersKey{}, NewResponseHeadersBag(responseHeaders))
}

func ResponseHeadersFromCtx(c *gin.Context) {
	if bag, ok := c.Request.Context().Value(BffResponseHeadersKey{}).(*ResponseHeadersBag); ok {
		responseHeadersMap := bag.All()
		for k, v := range responseHeadersMap {
			if len(v) > 0 {
				c.Header(k, v)
			}
		}
	}
}

// Service helpers section: these functions are used by the BffService to set headers on the outgoing request stream and capture headers from the upstream response stream into the context bag for propagation back to the client.

func SetBffRequestHeadersFromCtx(ctx context.Context, requestStream contracts.BffHttpRequestStream) {
	debugassert.Assert(requestStream != nil, "SetBffRequestHeadersFromCtx: requestStream must not be nil")
	debugassert.Assert(ctx != nil, "SetBffRequestHeadersFromCtx: ctx must not be nil")
	if m, ok := ctx.Value(BffRequestHeadersKey{}).(map[string]string); ok {
		for k, v := range m {
			requestStream.SetHeader(k, v)
		}
	}
}

func SetBffResponseHeadersInCtx(ctx context.Context, responseStream contracts.BffHttpResponseStream) {
	debugassert.Assert(ctx != nil, "SetBffResponseHeadersInCtx: ctx must not be nil")
	debugassert.Assert(responseStream != nil, "SetBffResponseHeadersInCtx: responseStream must not be nil")
	if bag, ok := ctx.Value(BffResponseHeadersKey{}).(*ResponseHeadersBag); ok {
		for k, v := range responseStream.Headers() {
			bag.SetKnown(k, strings.Join(v, ", "))
		}
	}
}
