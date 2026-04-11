package contracts

import (
	"context"
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

// ---------------------------------------------------------------------------
// Client port
// ---------------------------------------------------------------------------

// BffClient is the port for the upstream transport layer.
//
// It operates on raw path segments and serialisable payloads (any).
// BffRepository adapters are responsible for translating domain types
// into these primitives before delegating to BffClient.
//
// All methods propagate ctx for cancellation, deadline propagation,
// and OpenTelemetry trace-context injection by middleware.
type BffClient interface {
	// Post serialises upstream and sends it as POST to path.
	// The upstream response is deserialised into downstream.
	Post(ctx context.Context, path string, upstream, downstream any) error

	// Get sends GET to path with the given query parameters and
	// deserialises the response into downstream.
	Get(ctx context.Context, path string, queryParams map[string]string, downstream any) error

	// List sends GET to path with pagination query parameters and deserialises
	// the response body into downstream. Total item count, when available, is
	// part of the response body — extracting it is the responsibility of the
	// BffRepository adapter, not the transport client.
	List(ctx context.Context, path string, queryParams map[string]string, downstream any) error

	// Patch serialises upstream and sends it as PATCH to path/id.
	// The updated upstream resource is deserialised into downstream.
	Patch(ctx context.Context, path, id string, upstream, downstream any) error

	// Delete sends DELETE to path/id.
	Delete(ctx context.Context, path, id string) error
}

// ---------------------------------------------------------------------------
// Upstream errors
// ---------------------------------------------------------------------------

// BffUpstreamError represents a non-2xx response from the upstream API.
// BffRepository adapters return this so bffServiceImpl can distinguish
// upstream business errors (404, 409, 422) from transport failures.
type BffUpstreamError struct {
	StatusCode int
	Message    string
}

func (e *BffUpstreamError) Error() string {
	return e.Message
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
