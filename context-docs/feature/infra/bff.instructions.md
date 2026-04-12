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

**File**: `pkg/infra/bffclient/contracts/contracts.go`

### Stream abstraction

`BffClient` is payload-format agnostic. Callers (BffRepositories) control serialisation
by implementing stream interfaces:

```go
// BffRequestStream — what goes out (format-agnostic)
type BffRequestStream interface {
    Reader()      io.Reader // serialised body; nil = no body (GET, DELETE)
}

// BffResponseStream — what comes in (format-agnostic)
type BffResponseStream interface {
    Decode(r io.Reader) error // deserialise response body
}

// BffHttpRequestStream — HTTP specialisation of BffRequestStream
// Carries routing and format metadata the HTTP adapter needs.
type BffHttpRequestStream interface {
    BffRequestStream
    Method()      string            // "GET", "POST", "PATCH", "PUT", "DELETE"
    Path()        string            // resource path, e.g. "/incidents/INC001"
    Params()      map[string]string // query params; nil = none
    ContentType() string            // "application/json", "multipart/form-data; boundary=x"
}

// BffHttpResponseStream — HTTP specialisation of BffResponseStream
// Currently identical to BffResponseStream; extended when HTTP-specific
// response metadata is needed (e.g. raw status code for partial content).
type BffHttpResponseStream interface {
    BffResponseStream
}
```

**Rationale**: `BffRepository` is the stable port — application code never sees
`BffClient` directly. `BffClient` is a transport detail internal to the bffclient
adapter package. The stream abstraction is needed because BFF clients handle more
than JSON: image downloads, file uploads (multipart), and future binary formats.
For new protocols (gRPC, events), a new adapter implements `BffClient[Req, Resp]`
against `BffRepository` — `BffClient` itself does not need to change.

```go
// BffClient — generic over request/response stream types
type BffClient[Req BffRequestStream, Resp BffResponseStream] interface {
    Emit(ctx context.Context, req Req, resp Resp) error
}
```

**HTTP adapter** concrete type:
```go
var _ BffClient[BffHttpRequestStream, BffHttpResponseStream] = (*netHttpAdapter)(nil)
```

**Provided implementations** (in `internal/infra/bffclient/streams/`):
- `JsonStream[T]` — implements both `BffHttpRequestStream` and `BffHttpResponseStream`
  using `encoding/json`
- `FileDownloadStream` — implements `BffHttpResponseStream`, writes body to `io.Writer`

**Key decisions**:
- `BffHttpRequestStream.Method()` replaces the old per-verb methods (`Post`, `Get`, etc.)
  — the adapter routes based on the method value, not the Go method name.
- Path, params, and content-type live in the request stream — not as adapter method params.
- Headers (auth, correlation IDs) stay in context via `HeadersProxy` — not in the stream.
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
