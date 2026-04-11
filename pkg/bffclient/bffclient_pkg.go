package bffclient

import "github.com/brunojet/go-infra-backend/pkg/bffclient/contracts"

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
