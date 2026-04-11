---
applyTo: "internal/ports/bff/repositories/**,pkg/ports/bff/repositories/**"
---

# Ports/BFF — Repository Contracts

## Purpose

Define the BFF repository port contracts used by all upstream HTTP adapters.
Application code depends only on `pkg/ports/bff/repositories/contracts/` — never on
`internal/` directly.

## Package Map

```
pkg/ports/bff/repositories/
    contracts/contracts.go          ← BffRepository, BffNestedRepository, BffEntity, BffListParams
    bff_repositories_pkg.go         ← public facades: NewBffRepository, NewBffNestedRepository

internal/ports/bff/repositories/
    bff_repository_impl.go          ← concrete impl (BffClient → HTTP calls)
    bff_repository_impl_test.go     ← 11 tests using httptest.NewServer
```

## Type Parameter Conventions

Upstream APIs are asymmetric — request and response shapes differ per verb.

| Param | Role | Example |
|---|---|---|
| `CE` | upstream create DTO (POST body) | `ServiceNowCreateEntity` |
| `RE` | upstream read DTO (GET/List response) | `ServiceNowReadEntity` |
| `UE` | upstream update DTO (PATCH body) | `ServiceNowUpdateEntity` |

All `*E` types must implement `BffEntity` (provides `ResourceName() string`).

## Contracts

```go
// pkg/ports/bff/repositories/contracts/contracts.go

type BffEntity interface {
    ResourceName() string
}

type BffListParams struct {
    Offset int
    Limit  int
    Scopes map[string]any
}

type BffRepository[CE, RE, UE BffEntity] interface {
    Create(ctx context.Context, upstream CE, downstream *RE) error
    GetByID(ctx context.Context, id string, downstream *RE) error
    // List populates downstream with upstream page results.
    // Total count is NOT returned here — extracted by the mapper via ExtractUpstreamTotal.
    List(ctx context.Context, params BffListParams, downstream *[]RE) error
    Update(ctx context.Context, id string, upstream UE, downstream *RE) error
    Delete(ctx context.Context, id string) error
}

type BffNestedRepository[CE, RE, UE BffEntity] interface {
    BffRepository[CE, RE, UE]
    CreateNested(ctx context.Context, parentID string, upstream CE, downstream *RE) error
    // ListNested: total count extracted by mapper, not returned here.
    ListNested(ctx context.Context, parentID string, params BffListParams, downstream *[]RE) error
}
```

**Key decision**: `List` / `ListNested` return only `error`.
The total item count lives inside the upstream response body — only the mapper (`ExtractUpstreamTotal`)
knows how to extract it. See [ports/bff/service.instructions.md](service.instructions.md).

## Contract Hierarchy

```
BffClient (transport)
    ↓ used by
BffRepository[CE, RE, UE]       ← wire: pkg/ports/bff/repositories/bff_repositories_pkg.go
BffNestedRepository[CE, RE, UE]
    ↓ used by
BffService[C, R, U]             ← see ports/bff/service.instructions.md
```

Handlers are identical regardless of backing adapter (DB or BFF) — they always depend on
the domain-level `Service[C,R,U]` interface.

## Test Strategy

- All tests use `httptest.NewServer` — real HTTP round-trip, no mocks
- `BffClient` is obtained via `bffclient.NewNetHttpAdapter` (public facade)
- 11 tests in `internal/ports/bff/repositories/bff_repository_impl_test.go`
