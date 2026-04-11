---
applyTo: "internal/infra/bffclient/**,pkg/infra/bffclient/**"
---

# Infra/BFF — Transport Adapter

## Purpose

HTTP transport adapter that implements `BffClient` — the lowest-level mechanism for
communicating with upstream/legacy APIs. Handles retry, circuit-breaking, and OTel
tracing at the transport layer. Port contracts (repository and service interfaces) are
documented separately:

- Port contracts → [ports/bff/repository.instructions.md](../../ports/bff/repository.instructions.md)
- Service contracts → [ports/bff/service.instructions.md](../../ports/bff/service.instructions.md)

## Package Map

```
pkg/infra/bffclient/
    contracts/contracts.go               ← BffClient, BffHealthChecker, BffClientConfig, BffUpstreamError
    bffclient_pkg.go                     ← public facade (NewNetHttpAdapter)
internal/infra/bffclient/
    adapters/http_client.go              ← sole transport implementation (net/http + backoff + gobreaker)
    adapters/http_client_test.go         ← 12 tests using httptest.NewServer
internal/infra/observability/http_transports/
    otelhttp_transport.go                ← OTel RoundTripper (wraps otelhttp.NewTransport)
pkg/infra/observability/httptransports/
    httptransports_pkg.go                ← public facade for OTel transport
```

## BffClient — Transport Layer

**File**: `pkg/bffclient/contracts/contracts.go`

```go
type BffClient interface {
    Post(ctx context.Context, path string, upstream, downstream any) error
    Get(ctx context.Context, path string, queryParams map[string]string, downstream any) error
    List(ctx context.Context, path string, queryParams map[string]string, downstream any) error
    Patch(ctx context.Context, path, id string, upstream, downstream any) error
    Delete(ctx context.Context, path, id string) error
}
```

**Key decisions**:
- `List` returns only `error` — total count is part of the response body. Extracting it
  is the mapper's responsibility (`ExtractTotal`), not the transport's.
- Path segments are plain strings; the adapter derives full URLs via `BaseURL + path`.
- `BffHealthChecker` is separated so callers that only need CRUD don't depend on health.

## BffClient — Implementation

**File**: `internal/infra/bffclient/adapters/http_client.go`

- `net/http` for transport
- `cenkalti/backoff/v4` for exponential-backoff retry (`BffRetryConfig`)
- `sony/gobreaker` for circuit breaking (`BffCircuitBreakerConfig`)
- `backoff.Permanent` is **only** applied in `run()` via `isPermanentCallError()` — never
  inside `execute()`. This ensures `IsNotFound/IsConflict/IsUnprocessable` see the
  unwrapped `*BffUpstreamError` directly.
- 4xx responses (except 429 Too Many Requests) are **permanent** — not retried.
- `gobreaker.ErrOpenState` is permanent — not retried (circuit is open).

### Constructor

```go
// pkg/infra/bffclient/bffclient_pkg.go
func NewNetHttpAdapter(config contracts.BffClientConfig, transport http.RoundTripper) (contracts.BffClient, contracts.BffHealthChecker, error)
```

## OTel Plugin Pattern

Inject at construction time — never inside method bodies:

```go
// OTel only:
bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))

// OTel + token auth (planned):
bffclient.NewNetHttpAdapter(cfg,
    httptransports.NewOtelHttpTransport(
        bfftransports.NewTokenTransport(nil),
    ),
)
```

`nil` inner transport defaults to `http.DefaultTransport`.

## Error Model

**File**: `pkg/infra/bffclient/contracts/contracts.go`

```go
type BffUpstreamError struct {
    StatusCode int
    Body       []byte
}

var (
    IsUpstreamError func(err error) bool
    IsNotFound      func(err error) bool   // 404
    IsConflict      func(err error) bool   // 409
    IsUnprocessable func(err error) bool   // 422
)
```

## Pending

1. **Token auth RoundTripper** — `pkg/infra/bffclient/transports/`
   - Path 1: forward token from request context (caller → BFF → upstream)
   - Path 2: BFF manages own token cache/refresh
2. **Move `BffUpstreamError`** to `internal/infra/bff/errors/`

## Dependencies

- `github.com/sony/gobreaker v1.0.0`
- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0`
- `github.com/felixge/httpsnoop v1.0.4` (transitive)
- `github.com/cenkalti/backoff/v4`

## Test Strategy

- 12 tests in `internal/infra/bffclient/adapters/http_client_test.go` — `httptest.NewServer`, no mocks
- `newAdapter` helper uses `NewNetHttpAdapter(cfg, nil)` (nil = DefaultTransport)
