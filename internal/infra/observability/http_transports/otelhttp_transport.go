package transports

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewOtelHttpTransport wraps base with OTel instrumentation.
// It creates a child span for each outgoing request, injects the W3C
// traceparent/tracestate headers automatically, and records the response
// status and any error — using the global TracerProvider registered by the
// observability bootstrap.
//
// Usage — analogous to NewOtelGormPlugin() for GORM:
//
//	bffclient.NewNetHttpAdapter(cfg, transports.NewOtelHttpTransport(http.DefaultTransport))
//
// Pass nil to use http.DefaultTransport as the base.
func NewOtelHttpTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return otelhttp.NewTransport(base)
}
