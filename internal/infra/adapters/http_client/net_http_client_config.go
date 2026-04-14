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
	roundTrippers     []http.RoundTripper
}

// defaultHttpClientConfig returns a sensible default configuration.
func defaultHttpClientConfig() httpClientConfig {
	return httpClientConfig{
		baseURL:           "",
		connectTimeoutMs:  DefaultConnectTimeoutMs,
		responseTimeoutMs: DefaultResponseTimeoutMs,
		headers:           make(http.Header),
		roundTrippers:     make([]http.RoundTripper, 0),
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
		c.headers.Set(key, value)
	}
}

// WithRoundTripper adds a custom http.RoundTripper to the client's transport chain.
func WithRoundTripper(rt http.RoundTripper) HttpClientOption {
	return func(c *httpClientConfig) {
		c.roundTrippers = append(c.roundTrippers, rt)
	}
}
