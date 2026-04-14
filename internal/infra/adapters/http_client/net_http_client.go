package adapters

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

// netHttpClient is a lightweight adapter around net/http that implements the
// local `HttpClientAdapter` contract. It respects the configuration provided
// by `httpClientConfig` (base URL, headers, timeouts and transport chain).
type netHttpClient struct {
	client  *http.Client
	baseURL string
	headers http.Header
}

// NewNetHttpClient constructs a net/http based client using the provided
// functional options (see net_http_client_config.go). It selects the last
// non-nil RoundTripper from the config (typical usage is to pass an
// instrumentation transport such as an OpenTelemetry transport) and applies a
// connect timeout when the underlying transport is an *http.Transport.
func NewNetHttpClient(opts ...HttpClientOption) (*netHttpClient, error) {
	cfg := newHttpClientConfig(opts...)

	// choose the last non-nil transport (callers typically append instrumentation)
	var transport http.RoundTripper
	for i := len(cfg.rountTrippers) - 1; i >= 0; i-- {
		if cfg.rountTrippers[i] != nil {
			transport = cfg.rountTrippers[i]
			break
		}
	}
	if transport == nil {
		transport = http.DefaultTransport
	}

	// If the chosen transport is an *http.Transport, clone it and apply the
	// connect timeout. This avoids mutating shared/default transports.
	if t, ok := transport.(*http.Transport); ok {
		cloned := t.Clone()
		if cfg.connectTimeoutMs > 0 {
			cloned.DialContext = (&net.Dialer{Timeout: time.Duration(cfg.connectTimeoutMs) * time.Millisecond}).DialContext
		}
		transport = cloned
	}

	// Apply configured middleware builders in registration order so that the
	// first registered middleware becomes the outermost wrapper. We apply in
	// reverse to wrap correctly starting from the innermost transport.
	if len(cfg.middlewares) > 0 {
		for i := len(cfg.middlewares) - 1; i >= 0; i-- {
			if cfg.middlewares[i] == nil {
				continue
			}
			transport = cfg.middlewares[i](transport)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.responseTimeoutMs) * time.Millisecond,
	}

	adapter := &netHttpClient{
		client:  client,
		baseURL: cfg.baseURL,
		headers: cfg.headers,
	}

	return adapter, nil
}

// Do sends the provided request using the context and returns the response.
func (c *netHttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// build request (attach context, merge headers, resolve URL)
	r, err := c.buildRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// perform the request (with optional circuit breaker)
	return c.client.Do(r)
}

// buildRequest prepares the outgoing *http.Request: attaches the context,
// merges configured headers (server config overrides request headers), and
// resolves relative URLs against the client's baseURL.
func (c *netHttpClient) buildRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	if req == nil {
		return nil, errors.New("http_client: request is nil")
	}

	r := req.WithContext(ctx)

	// merge configured headers: server-side config wins (overrides request).
	if c.headers != nil {
		var merged http.Header
		if r.Header == nil {
			merged = make(http.Header)
		} else {
			merged = r.Header.Clone()
		}
		for k, vals := range c.headers {
			merged[k] = append([]string{}, vals...)
		}
		r.Header = merged
	}

	// resolve relative URLs against baseURL when provided
	if !r.URL.IsAbs() && c.baseURL != "" {
		if base, err := url.Parse(c.baseURL); err == nil {
			r.URL = base.ResolveReference(r.URL)
		}
	}

	return r, nil
}
