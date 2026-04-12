package bffclient

import (
	"net/http"

	internalbff "github.com/brunojet/go-infra-backend/internal/infra/bffclient"
	internaladapters "github.com/brunojet/go-infra-backend/internal/infra/bffclient/adapters"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	"github.com/gin-gonic/gin"
)

// ---- Contracts ----

type (
	BffClientConfig         = contracts.BffClientConfig
	BffMiddlewareConfig     = contracts.BffMiddlewareConfig
	BffRetryConfig          = contracts.BffRetryConfig
	BffCircuitBreakerConfig = contracts.BffCircuitBreakerConfig
	BffCircuitState         = contracts.BffCircuitState
	BffHealthChecker        = contracts.BffHealthChecker
	BffUpstreamError        = contracts.BffUpstreamError
	BffRequestStream        = contracts.BffRequestStream
	BffResponseStream       = contracts.BffResponseStream
	BffHttpRequestStream    = contracts.BffHttpRequestStream
	BffHttpResponseStream   = contracts.BffHttpResponseStream

	// BffHttpClient is the concrete BffClient instantiation for HTTP transports.
	// Use this type when wiring the HTTP adapter — BffClient[Req,Resp] itself
	// is rarely referenced directly outside adapter construction.
	BffHttpClient = contracts.BffClient[contracts.BffHttpRequestStream, contracts.BffHttpResponseStream]
)

// ---- Circuit state constants ----

const (
	BffCircuitClosed   = contracts.BffCircuitClosed
	BffCircuitOpen     = contracts.BffCircuitOpen
	BffCircuitHalfOpen = contracts.BffCircuitHalfOpen
)

// ---- Upstream error helpers ----

var (
	IsUpstreamError = contracts.IsUpstreamError
	IsNotFound      = contracts.IsNotFound
	IsConflict      = contracts.IsConflict
	IsUnprocessable = contracts.IsUnprocessable
)

// ---- Constructors ----

// NewNetHttpAdapter creates an HTTP adapter implementing BffClient and BffHealthChecker.
// Both return values point to the same underlying instance; callers may store
// them separately depending on what interface they need to expose.
//
// transport is the http.RoundTripper for every outgoing request.
// Pass httptransports.NewOtelHttpTransport(nil) to enable OTel tracing —
// analogous to passing gormplugins.NewOtelGormPlugin() to NewDatabaseManager.
// Pass nil to use http.DefaultTransport as-is.
func NewNetHttpAdapter(config contracts.BffClientConfig, transport http.RoundTripper) (BffHttpClient, contracts.BffHealthChecker, error) {
	a, err := internaladapters.NewNetHttpAdapter(config, transport)
	return a, a, err
}

// BffHeadersMiddleware returns a Gin middleware for bidirectional header propagation.
// Pass cfg.HeadersProxy from the BffClientConfig used to create the adapter.
func BffHeadersMiddleware(cfg contracts.BffMiddlewareConfig) gin.HandlerFunc {
	return internalbff.BffHeadersMiddleware(cfg)
}

// ---- Context helpers ----

var (
	WithRequestHeaders     = internalbff.WithRequestHeaders
	InitResponseCapture    = internalbff.InitResponseCapture
	ResponseHeadersFromCtx = internalbff.ResponseHeadersFromCtx
)
