package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/propagation"

	bffcts "github.com/brunojet/go-infra-backend/pkg/bffclient/contracts"
)

const (
	headerTotalCount  = "X-Total-Count"
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
	defaultTimeout    = 30 * time.Second
)

// netHttpAdapter implements BffClient and BffHealthChecker using net/http,
// cenkalti/backoff for exponential-backoff retry, and sony/gobreaker for
// circuit-breaking.
type netHttpAdapter struct {
	client      *http.Client
	config      bffcts.BffClientConfig
	middlewares []bffcts.BffClientMiddleware
	breaker     *gobreaker.CircuitBreaker // nil when circuit breaker is disabled
	propagator  propagation.TextMapPropagator
}

// NewNetHttpAdapter creates a new adapter that implements BffClient and BffHealthChecker.
// The caller can supply zero or more middlewares (e.g. OtelMiddleware) that will be
// called in order for every request.
func NewNetHttpAdapter(config bffcts.BffClientConfig, middlewares ...bffcts.BffClientMiddleware) (*netHttpAdapter, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("bffclient: BaseURL is required")
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	adapter := &netHttpAdapter{
		client:      &http.Client{Timeout: timeout},
		config:      config,
		middlewares: middlewares,
		propagator:  propagation.TraceContext{},
	}

	if config.CircuitBreaker.Enabled {
		cb := config.CircuitBreaker
		settings := gobreaker.Settings{
			Name:        config.BaseURL,
			MaxRequests: uint32(cb.HalfOpenRequests),
			Timeout:     cb.ResetTimeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= uint32(cb.MaxFailures)
			},
		}
		adapter.breaker = gobreaker.NewCircuitBreaker(settings)
	}

	return adapter, nil
}

// ---------------------------------------------------------------------------
// BffHealthChecker
// ---------------------------------------------------------------------------

func (a *netHttpAdapter) CircuitState() bffcts.BffCircuitState {
	if a.breaker == nil {
		return bffcts.BffCircuitClosed
	}
	switch a.breaker.State() {
	case gobreaker.StateOpen:
		return bffcts.BffCircuitOpen
	case gobreaker.StateHalfOpen:
		return bffcts.BffCircuitHalfOpen
	default:
		return bffcts.BffCircuitClosed
	}
}

func (a *netHttpAdapter) IsAvailable() bool {
	return a.CircuitState() != bffcts.BffCircuitOpen
}

// ---------------------------------------------------------------------------
// BffClient
// ---------------------------------------------------------------------------

func (a *netHttpAdapter) Post(ctx context.Context, path string, upstream, downstream any) error {
	return a.do(ctx, http.MethodPost, path, "", nil, upstream, downstream, nil)
}

func (a *netHttpAdapter) Get(ctx context.Context, path string, queryParams map[string]string, downstream any) error {
	return a.do(ctx, http.MethodGet, path, "", queryParams, nil, downstream, nil)
}

func (a *netHttpAdapter) List(ctx context.Context, path string, queryParams map[string]string, downstream any) (int64, error) {
	var total int64
	err := a.do(ctx, http.MethodGet, path, "", queryParams, nil, downstream, &total)
	return total, err
}

func (a *netHttpAdapter) Patch(ctx context.Context, path, id string, upstream, downstream any) error {
	return a.do(ctx, http.MethodPatch, path, id, nil, upstream, downstream, nil)
}

func (a *netHttpAdapter) Delete(ctx context.Context, path, id string) error {
	return a.do(ctx, http.MethodDelete, path, id, nil, nil, nil, nil)
}

// ---------------------------------------------------------------------------
// core executor
// ---------------------------------------------------------------------------

func (a *netHttpAdapter) do(
	ctx context.Context,
	method, path, id string,
	queryParams map[string]string,
	upstream, downstream any,
	totalOut *int64,
) error {
	reqInfo := bffcts.BffRequestInfo{Method: method, Path: path}

	// Notify middleware chain — OnRequest may add trace spans to the context.
	enriched := ctx
	for _, m := range a.middlewares {
		enriched = m.OnRequest(enriched, reqInfo)
	}

	// execute is the atomic unit of work, re-invoked on each retry attempt.
	// It rebuilds the request body reader each time so retries see a fresh stream.
	// It returns raw typed errors — the retry wrapper (run) decides what is permanent.
	var statusCode int
	execute := func() error {
		body, err := marshalBody(upstream)
		if err != nil {
			return err
		}

		httpReq, err := http.NewRequestWithContext(enriched, method, a.buildURL(path, id), body)
		if err != nil {
			return err
		}

		// Static headers from config (e.g. Authorization, Accept).
		for k, v := range a.config.Headers {
			httpReq.Header.Set(k, v)
		}
		if upstream != nil {
			httpReq.Header.Set(headerContentType, contentTypeJSON)
		}

		// Inject W3C traceparent/tracestate from the enriched context into HTTP headers.
		a.propagator.Inject(enriched, propagation.HeaderCarrier(httpReq.Header))

		// Apply caller-supplied query parameters.
		if len(queryParams) > 0 {
			q := httpReq.URL.Query()
			for k, v := range queryParams {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		resp, err := a.client.Do(httpReq)
		if err != nil {
			return err // transport error — retriable
		}
		defer resp.Body.Close()

		statusCode = resp.StatusCode

		if resp.StatusCode >= 400 {
			msg, _ := io.ReadAll(resp.Body)
			return &bffcts.BffUpstreamError{
				StatusCode: resp.StatusCode,
				Message:    strings.TrimSpace(string(msg)),
			}
		}

		// Extract pagination total from the de-facto standard header.
		if totalOut != nil {
			if v, parseErr := strconv.ParseInt(resp.Header.Get(headerTotalCount), 10, 64); parseErr == nil {
				*totalOut = v
			}
		}

		if downstream != nil && resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(resp.Body).Decode(downstream); err != nil {
				return fmt.Errorf("bffclient: decode response: %w", err)
			}
		}

		return nil
	}

	// Wrap execute with circuit breaker and/or retry.
	execErr := a.run(execute)

	// Notify middleware chain after the call completes (or fails).
	respInfo := bffcts.BffResponseInfo{StatusCode: statusCode, Err: execErr}
	for _, m := range a.middlewares {
		m.OnResponse(enriched, reqInfo, respInfo)
	}

	return execErr
}

// run wraps execute with the configured retry and circuit-breaker policies.
// Retry is the outer shell so each probe goes through the circuit breaker.
func (a *netHttpAdapter) run(execute func() error) error {
	cfg := a.config.Retry
	if cfg.MaxAttempts <= 1 {
		return a.runOnce(execute)
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = cfg.InitialDelay
	bo.MaxInterval = cfg.MaxDelay
	bo.Multiplier = cfg.Multiplier
	bo.MaxElapsedTime = 0 // controlled exclusively by MaxAttempts

	limited := backoff.WithMaxRetries(bo, uint64(cfg.MaxAttempts-1))

	return backoff.Retry(func() error {
		err := a.runOnce(execute)
		if err == nil {
			return nil
		}
		if isPermanentCallError(err) {
			return backoff.Permanent(err)
		}
		return err
	}, limited)
}

// runOnce calls execute through the circuit breaker (if enabled).
func (a *netHttpAdapter) runOnce(execute func() error) error {
	if a.breaker == nil {
		return execute()
	}
	_, err := a.breaker.Execute(func() (any, error) {
		return nil, execute()
	})
	return err
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// isPermanentCallError reports whether err should NOT be retried.
// 4xx client errors (except 429) and circuit-open states are permanent.
func isPermanentCallError(err error) bool {
	if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
		return true
	}
	upErr, ok := err.(*bffcts.BffUpstreamError)
	return ok && upErr.StatusCode >= 400 && upErr.StatusCode < 500 && upErr.StatusCode != http.StatusTooManyRequests
}

// buildURL joins BaseURL, resource path, and optional ID.
// e.g. ("https://sn.example.com", "incidents", "INC001") → "https://sn.example.com/incidents/INC001"
func (a *netHttpAdapter) buildURL(path, id string) string {
	base := strings.TrimRight(a.config.BaseURL, "/")
	p := strings.TrimLeft(path, "/")
	if id != "" {
		return fmt.Sprintf("%s/%s/%s", base, p, id)
	}
	return fmt.Sprintf("%s/%s", base, p)
}

// marshalBody serialises v to JSON and returns a fresh reader.
// Returns nil when v is nil (no-body requests such as GET, DELETE).
func marshalBody(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("bffclient: marshal request body: %w", err)
	}
	return bytes.NewReader(data), nil
}
