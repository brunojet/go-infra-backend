package bffclient

import (
	internaladapters "github.com/brunojet/go-infra-backend/internal/bffclient/adapters"
	"github.com/brunojet/go-infra-backend/pkg/bffclient/contracts"
)

// ---- Contracts ----

type (
	BffClientConfig         = contracts.BffClientConfig
	BffRetryConfig          = contracts.BffRetryConfig
	BffCircuitBreakerConfig = contracts.BffCircuitBreakerConfig
	BffCircuitState         = contracts.BffCircuitState
	BffRequestInfo          = contracts.BffRequestInfo
	BffResponseInfo         = contracts.BffResponseInfo
	BffClientMiddleware     = contracts.BffClientMiddleware
	BffHealthChecker        = contracts.BffHealthChecker
	BffClient               = contracts.BffClient
	BffUpstreamError        = contracts.BffUpstreamError
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
func NewNetHttpAdapter(config contracts.BffClientConfig, middlewares ...contracts.BffClientMiddleware) (contracts.BffClient, contracts.BffHealthChecker, error) {
	a, err := internaladapters.NewNetHttpAdapter(config, middlewares...)
	return a, a, err
}

// NewOtelMiddleware returns a BffClientMiddleware that traces each upstream
// call using the global OTel TracerProvider registered by the observability bootstrap.
func NewOtelMiddleware() contracts.BffClientMiddleware {
	return internaladapters.NewOtelMiddleware()
}
