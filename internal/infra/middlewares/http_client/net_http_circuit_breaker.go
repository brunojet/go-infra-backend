package middlewares

import (
	"errors"
	"net/http"

	"github.com/sony/gobreaker"
)

// breakerRoundTripper wraps a downstream RoundTripper and executes requests
// under a circuit breaker. When the breaker is nil it delegates directly to
// the next RoundTripper.
type breakerRoundTripper struct {
	next http.RoundTripper
	cb   *gobreaker.CircuitBreaker
}

// NewBreakerMiddleware returns a middleware builder: a function that accepts
// the next RoundTripper and returns a RoundTripper that applies circuit
// breaking according to cfg. This keeps middleware composition the client's
// responsibility (the client composes builders into a final transport).
func NewBreakerMiddleware(opts ...BreakerOption) func(next http.RoundTripper) http.RoundTripper {
	cfg := newCircuitBreakerConfig(opts...)
	return func(next http.RoundTripper) http.RoundTripper {
		settings := gobreaker.Settings{
			Name:        "breaker",
			MaxRequests: uint32(cfg.HalfOpenRequests),
			Timeout:     cfg.ResetTimeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= uint32(cfg.MaxFailures)
			},
		}
		return &breakerRoundTripper{next: next, cb: gobreaker.NewCircuitBreaker(settings)}
	}
}

func (b *breakerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if b.next == nil {
		return nil, errors.New("breaker middleware: next RoundTripper is nil; ensure client provides base transport")
	}
	if b.cb == nil {
		return b.next.RoundTrip(req)
	}
	var resp *http.Response
	_, err := b.cb.Execute(func() (any, error) {
		r, err := b.next.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		resp = r
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
