package adapters

import (
	"net/http"
)

const (
	DefaultConnectTimeoutMs  = 15 * 1000
	DefaultResponseTimeoutMs = 30 * 1000
)

type httpClientConfig struct {
	baseURL           string
	headers           http.Header
	connectTimeoutMs  int
	responseTimeoutMs int
	rountTrippers     []http.RoundTripper
	// middlewares are RoundTripper decorators applied by the client when
	// assembling the final transport chain. Each middleware receives the
	// next RoundTripper and must return a RoundTripper wrapping it.
	middlewares []func(http.RoundTripper) http.RoundTripper
}

// defaultHttpClientConfig returns a sensible default configuration.
func defaultHttpClientConfig() httpClientConfig {
	return httpClientConfig{
		baseURL:           "",
		headers:           make(http.Header),
		connectTimeoutMs:  DefaultConnectTimeoutMs,
		responseTimeoutMs: DefaultResponseTimeoutMs,
		rountTrippers:     []http.RoundTripper{http.DefaultTransport},
		middlewares:       nil,
	}
}

// HttpClientOption configures a HttpClientConfig when constructing a client.
type HttpClientOption func(cfg *httpClientConfig)

// newHttpClientConfig builds a HttpClientConfig applying provided options.
func newHttpClientConfig(opts ...HttpClientOption) httpClientConfig {
	cfg := defaultHttpClientConfig()
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&cfg)
	}
	return cfg
}

// NOTE: Circuit breakers are provided as standard http.RoundTripper
// middleware. Use `WithRoundTripper(...)` and supply the BreakerRoundTripper
// (or any other RoundTripper) to enable circuit-breaking behaviour.

// NOTE: We intentionally avoid a single-blob `WithConfig` helper because the
// functional options are small and composable (WithBaseURL/WithHeaders/etc.).

// WithBaseURL sets the base URL used by the client.
func WithBaseURL(url string) HttpClientOption {
	return func(c *httpClientConfig) { c.baseURL = url }
}

// WithTimeout sets the client timeout in milliseconds.
func WithTimeout(connectionTimeoutMs, responseTimeoutMs int) HttpClientOption {
	return func(c *httpClientConfig) {
		c.connectTimeoutMs = connectionTimeoutMs
		c.responseTimeoutMs = responseTimeoutMs
	}
}

// WithHeader sets a single header key/value.
func WithHeader(key, value string) HttpClientOption {
	return func(c *httpClientConfig) {
		if c.headers == nil {
			c.headers = make(http.Header)
		}
		c.headers.Set(key, value)
	}
}

// WithRoundTripper adds a custom http.RoundTripper to the client's transport chain.
func WithRoundTripper(rt http.RoundTripper) HttpClientOption {
	return func(c *httpClientConfig) {
		c.rountTrippers = append(c.rountTrippers, rt)
	}
}

// WithRoundTripperMiddleware appends a RoundTripper middleware builder that
// will be applied by the client when assembling the final transport chain.
// Middleware builders have signature `func(next http.RoundTripper) http.RoundTripper`.
func WithRoundTripperMiddleware(mw func(http.RoundTripper) http.RoundTripper) HttpClientOption {
	return func(c *httpClientConfig) {
		c.middlewares = append(c.middlewares, mw)
	}
}
