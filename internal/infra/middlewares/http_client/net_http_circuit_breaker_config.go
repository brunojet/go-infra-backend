package middlewares

import "time"

// circuitBreakerConfig configures the circuit breaker used by the
// Breaker middleware.
type circuitBreakerConfig struct {
	MaxFailures      int           // consecutive failures before opening
	ResetTimeout     time.Duration // how long to wait before transitioning to half-open
	HalfOpenRequests int           // probe requests allowed in half-open state
}

// BreakerOption configures a CircuitBreakerConfig.
type BreakerOption func(cfg *circuitBreakerConfig)

// newCircuitBreakerConfig builds a CircuitBreakerConfig applying provided
// functional options.
func newCircuitBreakerConfig(opts ...BreakerOption) circuitBreakerConfig {
	cfg := circuitBreakerConfig{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&cfg)
	}
	return cfg
}

// WithCircuitBreakerMaxFailures sets the number of consecutive failures
// required to open the circuit.
func WithCircuitBreakerMaxFailures(n int) BreakerOption {
	return func(c *circuitBreakerConfig) { c.MaxFailures = n }
}

// WithCircuitBreakerResetTimeout sets the reset timeout for the breaker.
func WithCircuitBreakerResetTimeout(d time.Duration) BreakerOption {
	return func(c *circuitBreakerConfig) { c.ResetTimeout = d }
}

// WithCircuitBreakerHalfOpenRequests sets how many probe requests are
// allowed while the breaker is half-open.
func WithCircuitBreakerHalfOpenRequests(n int) BreakerOption {
	return func(c *circuitBreakerConfig) { c.HalfOpenRequests = n }
}
