---
applyTo: "internal/bffclient/**,pkg/bffclient/**,pkg/ports/bff/**,demoapp/repositories/**,demoapp/services/**"
---

# BFF Layer — Architectural Context

## Purpose

The BFF (Backend For Frontend) layer bridges the application's domain model against
asymmetric upstream/legacy APIs (e.g. ServiceNow). It follows the same
Ports & Adapters structure used by the GORM database layer, replacing database types
with HTTP types.

## Package Map

```
pkg/ports/bff/
    repositories/contracts/contracts.go      ← BffRepository, BffNestedRepository ports
    repositories/bff_repositories_pkg.go     ← public facades (NewBffRepository, NewBffNestedRepository)
    services/contracts/contracts.go          ← BffServiceMapper, BffService, BffNestedService ports
    services/bff_services_pkg.go             ← public facades (NewBffServiceImpl, NewBffNestedServiceImpl)
pkg/bffclient/
    contracts/contracts.go               ← BffClient, BffHealthChecker, BffClientConfig, BffUpstreamError
    bffclient_pkg.go                     ← public facade (NewNetHttpAdapter)
internal/bffclient/
    adapters/http_client.go              ← sole transport implementation (net/http + backoff + gobreaker)
    adapters/http_client_test.go         ← 12 tests using httptest.NewServer
internal/ports/bff/
    repositories/bff_repository_impl.go      ← BffRepository + BffNestedRepository concrete impl
    repositories/bff_repository_impl_test.go ← 11 tests using httptest.NewServer
    services/bff_service_impl.go             ← BffService + BffNestedService concrete impl
    services/bff_service_impl_test.go        ← 22 tests using stubs
internal/observability/http_transports/
    otelhttp_transport.go                ← OTel RoundTripper (wraps otelhttp.NewTransport)
pkg/observability/httptransports/
    httptransports_pkg.go                ← public facade for OTel transport
```

## Type Parameter Conventions

Upstream APIs are asymmetric — request and response shapes differ per verb.
Three explicit type parameters are always used:

| Param | Role | Example |
|---|---|---|
| `CE` | upstream create DTO (POST body) | `ServiceNowCreateEntity` |
| `RE` | upstream read DTO (GET/List response) | `ServiceNowReadEntity` |
| `UE` | upstream update DTO (PATCH body) | `ServiceNowUpdateEntity` |
| `C` | domain create DTO | `CreateApplication` |
| `R` | domain read/response DTO | `Application` |
| `U` | domain update DTO | `UpdateApplication` |

All `*E` types must implement `BffEntity` (provides `ResourceName() string`).

## Contract Hierarchy

```
BffClient                          ← transport only: Post/Get/List/Patch/Delete
    ↓ used by
BffRepository[CE, RE, UE]          ← resource-aware: Create/GetByID/List/Update/Delete
BffNestedRepository[CE, RE, UE]    ← adds CreateNested/ListNested (parent path)
    ↓ used by
BffServiceMapper[C,R,U,CE,RE,UE]   ← translates domain ↔ upstream DTOs
BffService[C,R,U]                  ← same interface as database-backed Service
BffNestedService[C,R,U]            ← same as NestedService
```

Handlers are **identical** regardless of whether the backing adapter is a database or an
upstream HTTP API — they always depend on `BffService` / `BffNestedService`.

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

**File**: `internal/bffclient/adapters/http_client.go`

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
func NewNetHttpAdapter(config BffClientConfig, transport http.RoundTripper) (*netHttpAdapter, error)
```

Public facade re-exports as:

```go
// pkg/bffclient/bffclient_pkg.go
func NewNetHttpAdapter(config contracts.BffClientConfig, transport http.RoundTripper) (contracts.BffClient, contracts.BffHealthChecker, error)
```

## OTel / Observability Plugin Pattern

Same pattern used across all adapters — inject at construction time, never inside method bodies:

| Layer | Plugin mechanism |
|---|---|
| GORM | `db.Use(gormplugins.NewOtelGormPlugin())` |
| Gin | `router.Use(httpmiddlewares.OtelGinMiddleware())` |
| BFF/net/http | `bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))` |

`NewOtelHttpTransport` wraps `otelhttp.NewTransport` from
`go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`.
It uses the global `TracerProvider` registered by the observability bootstrap —
no extra configuration needed.

### Transport composition chain

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

## BffServiceMapper Contract

**File**: `pkg/ports/bff/services/contracts/contracts.go`

```go
type BffServiceMapper[C, R, U, CE, RE, UE any] interface {
    ToUpstreamPost(dto C, upstream *CE) error
    ToUpstreamPatch(dto U, upstream *UE) error
    ToDomainDTO(upstream *RE, dto *R) error
    GetUpstreamID(id string) (string, error)
    ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
    // ExtractUpstreamTotal reads the total item count from a slice of upstream
    // read DTOs. Only the mapper knows the RE shape and where the count lives
    // (e.g. a metadata field in the first element), so extraction is delegated here
    // rather than handled by the repository or transport layers.
    ExtractUpstreamTotal(upstream []RE) int64
}
```

**Decision** (implemented): `ExtractUpstreamTotal(upstream []RE) int64` was added to `BffServiceMapper`.
Rationale: `bffRepositoryImpl` does not know the DTO structure; only the mapper knows `RE`
and can extract the total from the response body.
`BffRepository.List` and `BffNestedRepository.ListNested` now return `error` only — the
total is surfaced by the service layer after calling `ExtractUpstreamTotal`.

## BffRepository Contract

**File**: `pkg/ports/bff/repositories/contracts/contracts.go`

```go
type BffRepository[CE, RE, UE BffEntity] interface {
    Create(ctx context.Context, upstream CE, downstream *RE) error
    GetByID(ctx context.Context, id string, downstream *RE) error
    // List populates downstream with the upstream page results. The total item count
    // is extracted by the mapper via ExtractUpstreamTotal, not returned here.
    List(ctx context.Context, params BffListParams, downstream *[]RE) error
    Update(ctx context.Context, id string, upstream UE, downstream *RE) error
    Delete(ctx context.Context, id string) error
}

type BffNestedRepository[CE, RE, UE BffEntity] interface {
    BffRepository[CE, RE, UE]
    CreateNested(ctx context.Context, parentID string, upstream CE, downstream *RE) error
    // ListNested populates downstream with the upstream page results scoped to parentID.
    // The total item count is extracted by the mapper via ExtractUpstreamTotal.
    ListNested(ctx context.Context, parentID string, params BffListParams, downstream *[]RE) error
}
```

## Error Model

**File**: `pkg/bffclient/contracts/contracts.go` (planned move: `internal/bff/errors/`)

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

## Pending Implementation Work

In priority order:

1. **Token auth RoundTripper plugin** — package TBD (`pkg/bffclient/transports/`?)
   Two paths under consideration:
   - Path 1: Forward token from request context (POS → BFF → upstream)
   - Path 2: BFF manages own token cache/refresh internally
2. **Move `BffUpstreamError`** to `internal/bff/errors/`

## Go Version & Dependencies

- Go: `1.25.0` (upgraded from 1.24.0 — required by `otelhttp v0.68.0`)
- `github.com/sony/gobreaker v1.0.0`
- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0`
- `github.com/felixge/httpsnoop v1.0.4` (transitive, otelhttp)
- `github.com/cenkalti/backoff/v4` (was transitive, promoted to direct)

## Test Strategy

- All BFF transport and repository tests use `httptest.NewServer` — real HTTP round-trip, no mocks
- Service tests use hand-written stubs (no gomock) — same pattern as DB service tests
- `newAdapter` helper uses `NewNetHttpAdapter(cfg, nil)` (nil = DefaultTransport)
- Repository tests use `bffclient.NewNetHttpAdapter` (public facade) to obtain a `BffClient`
- 12 tests in `internal/bffclient/adapters/http_client_test.go`
- 11 tests in `internal/ports/bff/repositories/bff_repository_impl_test.go`
- 22 tests in `internal/ports/bff/services/bff_service_impl_test.go`
