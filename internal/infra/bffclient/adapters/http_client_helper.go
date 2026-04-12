package adapters

import (
	"context"
	"net/http"
	"strings"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/internal/infra/bffclient/types"
)

func SetBffRequestHeadersFromCtx(ctx context.Context, headers *http.Header) {
	debugassert.Assert(headers != nil, "SetBffRequestHeadersFromCtx: headers must not be nil")
	if m, ok := ctx.Value(types.BffRequestHeadersKey{}).(map[string]string); ok {
		for k, v := range m {
			headers.Set(k, v)
		}
	}
}

func SetBffResponseHeadersInCtx(ctx context.Context, headers http.Header) {
	debugassert.Assert(headers != nil, "SetBffResponseHeadersInCtx: headers must not be nil")
	if bag, ok := ctx.Value(types.BffResponseHeadersKey{}).(*types.ResponseHeadersBag); ok {
		for k, v := range headers {
			bag.SetKnown(k, strings.Join(v, ", "))
		}
	}
}
