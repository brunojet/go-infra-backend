---
applyTo: "internal/ports/bff/services/**,pkg/ports/bff/services/**"
---

# Ports/BFF — Service Contracts

## Purpose

Define the BFF service and mapper contracts that bridge domain DTOs (`C`, `R`, `U`) with
upstream API DTOs (`CE`, `RE`, `UE`). The service interface is **identical** to the
database-backed `Service[C,R,U]` — handlers cannot tell the difference.

## Package Map

```
pkg/ports/bff/services/
    contracts/contracts.go          ← BffServiceMapper, BffService, BffNestedService
    bff_services_pkg.go             ← public facades: NewBffServiceImpl, NewBffNestedServiceImpl

internal/ports/bff/services/
    bff_service_impl.go             ← concrete impl
    bff_service_impl_test.go        ← 22 tests using stubs
```

## Type Parameters

| Param | Role |
|---|---|
| `C` | domain create DTO |
| `R` | domain read/response DTO |
| `U` | domain update DTO |
| `CE` | upstream create DTO |
| `RE` | upstream read DTO |
| `UE` | upstream update DTO |

## BffServiceMapper Contract

```go
// pkg/ports/bff/services/contracts/contracts.go

type BffServiceMapper[C, R, U, CE, RE, UE any] interface {
    ToUpstreamPost(dto C, upstream *CE) error
    ToUpstreamPatch(dto U, upstream *UE) error
    ToDomainDTO(upstream *RE, dto *R) error
    GetUpstreamID(id string) (string, error)
    ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
    // ExtractUpstreamTotal reads the total item count from a slice of upstream read DTOs.
    // Only the mapper knows the RE shape and where the count lives (e.g. a metadata
    // field in the first element), so extraction is delegated here rather than to
    // the repository or transport layers.
    ExtractUpstreamTotal(upstream []RE) int64
}
```

**Decision**: `ExtractUpstreamTotal(upstream []RE) int64` lives on the mapper because
`BffRepositoryImpl` is generic and has no knowledge of `RE`'s structure.
`BffRepository.List` and `BffNestedRepository.ListNested` return only `error` — the total
is surfaced by the service after calling `ExtractUpstreamTotal`.

## BffService / BffNestedService Contracts

```go
// Same interface shape as backend Service[C,R,U] — handlers are interchangeable

type BffService[C, R, U any] interface {
    Create(ctx context.Context, request C, response *R) error
    GetByID(ctx context.Context, id string, response *R) error
    List(ctx context.Context, queryScopes map[string]any, response *[]R) (int64, error)
    Update(ctx context.Context, id string, request U, response *R) error
    Delete(ctx context.Context, id string) error
}

type BffNestedService[C, R, U any] interface {
    BffService[C, R, U]
    CreateNested(ctx context.Context, parentID string, request C, response *R) error
    ListNested(ctx context.Context, parentID string, queryScopes map[string]any, response *[]R) (int64, error)
}
```

## Contract Chain

```
BffServiceMapper[C,R,U,CE,RE,UE]   ← app implements this
    +
BffRepository[CE,RE,UE]            ← see ports/bff/repository.instructions.md
    ↓ wired by
BffServiceImpl                     ← internal/ports/bff/services/bff_service_impl.go
    ↓ exposes
BffService[C,R,U]                  ← consumed by handlers (same as DB-backed Service)
```

## Test Strategy

- 22 tests in `internal/ports/bff/services/bff_service_impl_test.go`
- Use hand-written stubs (no gomock) — same pattern as backend service tests
- Stub implements `BffRepository` and `BffServiceMapper` interfaces directly
