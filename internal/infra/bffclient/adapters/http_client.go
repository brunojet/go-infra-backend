package adapters

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"

	"github.com/brunojet/go-infra-backend/debugassert"
	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
)

// netHttpAdapter implements BffClient[BffHttpRequestStream, BffHttpResponseStream]
// and BffHealthChecker using net/http, cenkalti/backoff for exponential-backoff
// retry, and sony/gobreaker for circuit-breaking.
//
// OTel tracing is injected via http.RoundTripper (e.g. otelhttp.NewTransport),
// following the same pattern used by the GORM adapter (gorm.Plugin) and
// the Gin adapter (gin.HandlerFunc).
type netHttpAdapter struct {
	client  *http.Client
	config  bffcts.BffClientConfig
	breaker *gobreaker.CircuitBreaker // nil when circuit breaker is disabled
}

// NewNetHttpAdapter creates a new adapter implementing
// BffClient[BffHttpRequestStream, BffHttpResponseStream] and BffHealthChecker.
//
// transport is the http.RoundTripper used for every outgoing request.
// Pass otelhttp.NewTransport(http.DefaultTransport) to enable OTel tracing,
// or nil to use http.DefaultTransport as-is.
func NewNetHttpAdapter(config bffcts.BffClientConfig, transport http.RoundTripper) (*netHttpAdapter, error) {
	if config.BaseURL == "" {
		return nil, errBaseURLRequired
	}

	if transport == nil {
		transport = http.DefaultTransport
	}

	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	adapter := &netHttpAdapter{
		client: &http.Client{Timeout: timeout, Transport: transport},
		config: config,
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
// BffClient[BffHttpRequestStream, BffHttpResponseStream]
// ---------------------------------------------------------------------------

// Emit executes the HTTP request described by req and deserialises the
// response into resp. Both req and resp are required:
//   - req nil: not valid (method and path are required)
//   - resp nil: not valid — always provide a stream; use NoBodyResponseStream for body-less operations (e.g. DELETE)
func (a *netHttpAdapter) Emit(ctx context.Context, bffRequest bffcts.BffHttpRequestStream, bffResponse bffcts.BffHttpResponseStream) error {
	debugassert.Assert(bffRequest != nil, "bffclient: bffRequest stream is required")
	debugassert.Assert(bffResponse != nil, "bffclient: bffResponse stream is required")
	url := a.buildURL(bffRequest.Path())
	SetBffClientHeaders(a.config.Headers, bffRequest) // must be last to be not overridden by context headers

	execute := func() error {
		httpReq, err := http.NewRequestWithContext(ctx, bffRequest.Method(), url, bffRequest.Reader())
		if err != nil {
			return err
		}
		httpReq.Header = bffRequest.Headers().Clone()
		httpReq.URL.RawQuery = bffRequest.RawQuery()

		httpResp, err := a.client.Do(httpReq)
		if err != nil {
			return err
		}
		defer httpResp.Body.Close()
		return a.buildResponse(httpResp, bffResponse)
	}

	return a.run(execute)
}

// ---------------------------------------------------------------------------
// retry + circuit breaker
// ---------------------------------------------------------------------------

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
// Only circuit-breaker state errors are permanent; all other transport
// errors are retried up to MaxAttempts.
func isPermanentCallError(err error) bool {
	return errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests)
}

func (a *netHttpAdapter) buildResponse(httpResp *http.Response, bffResp bffcts.BffHttpResponseStream) error {
	debugassert.Assert(bffResp != nil, "bffclient: resp stream is required — use NoBodyResponseStream for body-less operations")
	debugassert.Assert(httpResp != nil, "buildResponse: httpResp must not be nil")
	bffResp.SetStatusCode(httpResp.StatusCode)
	bffResp.MergeHeader(httpResp.Header)
	if httpResp.StatusCode == http.StatusNoContent || httpResp.StatusCode >= http.StatusBadRequest {
		return nil
	}
	return bffResp.Decode(httpResp.Body)
}

// buildURL joins BaseURL and the resource path provided by the request stream.
// Path is normalised (leading/trailing slashes).
func (a *netHttpAdapter) buildURL(path string) string {
	base := strings.TrimRight(a.config.BaseURL, "/")
	p := strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s/%s", base, p)
}
