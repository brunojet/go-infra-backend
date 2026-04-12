package contracts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

// BffMiddlewareConfig configures bidirectional header propagation for a BFF client.
// It is used by both the BFF adapter (to filter headers at transport time) and
// by BffHeadersMiddleware (to filter headers at the Gin layer).
type BffMiddlewareConfig struct {
	RequestHeaders  []string // incoming header keys to forward from ctx → upstream (e.g. "Authorization")
	ResponseHeaders []string // upstream response header keys to capture into ctx (e.g. "X-Request-Id")
}

// BffClientConfig holds all configuration for a BFF client instance.
// It is the bffclient counterpart of database/contracts.DatabaseConfig.
type BffClientConfig struct {
	BaseURL        string
	Timeout        time.Duration
	Headers        map[string]string   // static headers applied to every request (e.g. Content-Type, API keys)
	HeadersProxy   BffMiddlewareConfig // dynamic per-request header propagation (request ctx → upstream, upstream → response ctx)
	Retry          BffRetryConfig
	CircuitBreaker BffCircuitBreakerConfig
}

// BffRetryConfig configures exponential-backoff retry behaviour.
// Set MaxAttempts to 0 to disable retries entirely.
type BffRetryConfig struct {
	MaxAttempts  int           // total attempts including the first call (0 = no retry)
	InitialDelay time.Duration // delay before the first retry
	MaxDelay     time.Duration // upper bound for delay regardless of multiplier
	Multiplier   float64       // backoff growth factor, e.g. 2.0 for exponential
}

// BffCircuitBreakerConfig configures the circuit breaker that wraps upstream calls.
type BffCircuitBreakerConfig struct {
	Enabled          bool
	MaxFailures      int           // consecutive failures before opening the circuit
	ResetTimeout     time.Duration // how long to wait before transitioning to half-open
	HalfOpenRequests int           // probe requests allowed in half-open state
}

// ---------------------------------------------------------------------------
// Circuit breaker state (observable by health/readiness probes)
// ---------------------------------------------------------------------------

// BffCircuitState represents the current state of the circuit breaker.
type BffCircuitState string

const (
	BffCircuitClosed   BffCircuitState = "closed"    // normal operation
	BffCircuitOpen     BffCircuitState = "open"      // failing fast, upstream unreachable
	BffCircuitHalfOpen BffCircuitState = "half-open" // probing upstream recovery
)

// ---------------------------------------------------------------------------
// Middleware (extension point for OpenTelemetry, logging, metrics)
// ---------------------------------------------------------------------------

// BffRequestInfo carries metadata about an outgoing upstream request.
// Passed to middleware hooks so they can create trace spans, emit metrics, etc.
type BffRequestInfo struct {
	Method string
	Path   string
}

// BffResponseInfo carries metadata about a completed upstream call.
// Err is non-nil on transport errors or non-2xx responses mapped to errors.
type BffResponseInfo struct {
	StatusCode int
	Err        error
}

// BffClientMiddleware is the extension point for cross-cutting concerns.
//
// Typical implementations:
//   - OpenTelemetry: OnRequest starts a span and injects trace headers into ctx;
//     OnResponse records status and ends the span.
//   - Structured logging: OnRequest logs the outgoing call; OnResponse logs
//     latency and outcome.
//   - Metrics: OnResponse increments counters and records latency histograms.
//
// OnRequest is called before the request is sent. The returned context is
// propagated to the upstream call and to OnResponse — use it to carry span
// and correlation-ID values downstream.
// OnResponse is always called, even on error.
type BffClientMiddleware interface {
	OnRequest(ctx context.Context, req BffRequestInfo) context.Context
	OnResponse(ctx context.Context, req BffRequestInfo, resp BffResponseInfo)
}

// ---------------------------------------------------------------------------
// Health checker (separated from BffClient so mocks stay minimal)
// ---------------------------------------------------------------------------

// BffHealthChecker exposes the internal state of the BFF client adapter.
// The concrete adapter implements both BffClient and BffHealthChecker.
// Callers that only need CRUD operations depend only on BffClient.
type BffHealthChecker interface {
	CircuitState() BffCircuitState
	IsAvailable() bool
}

// BffRequestStream is the format-agnostic interface for serialising a request
// body to send to an upstream. Implementations control the format: JSON,
// multipart, binary, etc.
type BffRequestStream interface {
	Reader() io.Reader // serialised request body; nil = no body (GET, DELETE)
}

// BffResponseStream is the format-agnostic interface for deserialising a
// response body received from an upstream.
type BffResponseStream interface {
	Decode(r io.Reader) error // deserialise the response body
}

type BffHttpHeaders interface {
	SetHeader(key, value string) // extra per-call headers injected by the mapper
	SetHeaders(http.Header)      // extra per-call headers injected by the mapper
	Headers() http.Header        // read by the adapter to merge into the outgoing request
}

// BffHttpRequestStream extends BffRequestStream with HTTP-specific metadata.
// ContentType is defined here — not in BffRequestStream — because it is an
// HTTP concept; gRPC and event adapters do not use Content-Type headers.
type BffHttpRequestStream interface {
	BffRequestStream
	BffHttpHeaders
	Method() string   // HTTP verb: "GET", "POST", "PATCH", "PUT", "DELETE"
	Path() string     // resource path, e.g. "/incidents/INC001"
	RawQuery() string // pre-encoded query parameters for direct injection into the URL; mutually exclusive with Params()
}

// BffHttpResponseStream extends BffResponseStream for HTTP adapters.
// Currently mirrors BffResponseStream; extended here so future HTTP-specific
// response metadata (e.g. raw status code for 206 Partial Content) can be
// added without touching BffResponseStream.
type BffHttpResponseStream interface {
	BffResponseStream
	BffHttpHeaders
	SetStatusCode(code int)
}

// ---------------------------------------------------------------------------
// Client port
// ---------------------------------------------------------------------------

// BffClient is the protocol-agnostic transport port.
//
// Req carries all outgoing information (routing + serialised body).
// Resp receives and deserialises the response body.
// The concrete type of Req determines which adapter handles the call:
//   - BffHttpRequestStream → net/http adapter
//   - BffGrpcRequestStream → gRPC adapter (future)
//
// BffRepository adapters are responsible for constructing Req/Resp and
// translating domain types before delegating to BffClient.
// Application code never imports BffClient directly.
//
// ctx propagates cancellation, deadlines, and OTel trace context.
// Request/response headers are propagated via the context (HeadersProxy) —
// not via the stream interfaces.
type BffClient[Req BffRequestStream, Resp BffResponseStream] interface {
	Emit(ctx context.Context, req Req, resp Resp) error
}

// ---------------------------------------------------------------------------
// Upstream errors
// ---------------------------------------------------------------------------

// BffUpstreamError represents a non-2xx response from the upstream API.
// BffRepository adapters return this so bffServiceImpl can distinguish
// upstream business errors (404, 409, 422) from transport failures.
//
// Body holds the raw response bytes so that ExtractUpstreamError in the mapper
// can decode them according to the upstream's error contract (e.g. RFC 9457
// Problem Details JSON) without the transport layer pre-processing the content.
type BffUpstreamError struct {
	StatusCode int
	Body       []byte
}

func (e *BffUpstreamError) Error() string {
	return fmt.Sprintf("upstream error %d", e.StatusCode)
}

// IsUpstreamError reports whether err is a BffUpstreamError.
func IsUpstreamError(err error) bool {
	_, ok := err.(*BffUpstreamError)
	return ok
}

// IsNotFound reports whether err represents an upstream 404.
func IsNotFound(err error) bool {
	e, ok := err.(*BffUpstreamError)
	return ok && e.StatusCode == 404
}

// IsConflict reports whether err represents an upstream 409.
func IsConflict(err error) bool {
	e, ok := err.(*BffUpstreamError)
	return ok && e.StatusCode == 409
}

// IsUnprocessable reports whether err represents an upstream 422.
func IsUnprocessable(err error) bool {
	e, ok := err.(*BffUpstreamError)
	return ok && e.StatusCode == 422
}
