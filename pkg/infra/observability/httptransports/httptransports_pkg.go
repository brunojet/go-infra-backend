package httptransports

import (
	"net/http"

	internaltransports "github.com/brunojet/go-infra-backend/internal/infra/observability/http_transports"
)

// NewOtelHttpTransport wraps base with OTel instrumentation for outgoing HTTP calls.
//
// Mirrors the gormplugins.NewOtelGormPlugin() / httpmiddlewares.OtelGinMiddleware()
// pattern — callers inject it at construction time:
//
//	bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))
func NewOtelHttpTransport(base http.RoundTripper) http.RoundTripper {
	return internaltransports.NewOtelHttpTransport(base)
}
