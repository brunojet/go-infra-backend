---
applyTo: "internal/ports/backend/services/**,pkg/ports/backend/services/**"
---

# Ports — Backend Service Contracts

## Purpose

Generic GORM-backed `Service` and `NestedService` implementations that bridge repository
operations with service-layer business rules. Application services (`demoapp/services/`)
depend only on the `pkg/` contracts — never on `internal/` directly.

## Package Map

```
internal/ports/backend/services/
    service_impl.go                  ← Service[C,R,U,E] concrete impl
    nested_service_impl.go           ← NestedService[C,R,U,E] concrete impl
    service_utils.go                 ← shared helpers (IsSubSetInterface, etc.)
    consts.go                        ← default page sizes, sort constants
    errors.go                        ← service-layer error sentinels re-exported
    service_impl_test.go             ← unit tests for Service
    nested_service_impl_test.go      ← unit tests for NestedService
    repository_impl_mock_helper_test.go ← mock repository builder
    repository_impl_mock_test.go     ← mock repository type

pkg/ports/backend/services/
    services_pkg.go                  ← public facades: NewServiceImpl, NewNestedServiceImpl
    contracts/contracts.go           ← Service[C,R,U], NestedService[C,R,U], ServiceMapper interfaces
```

## Type Parameters

| Param | Role |
|---|---|
| `C` | Create DTO (request) |
| `R` | Read DTO (response) |
| `U` | Update DTO (request) |
| `E` | Domain entity / model (must implement `rpocts.Entity`) |

## Contract

```go
// pkg/ports/backend/services/contracts/contracts.go
type Service[C, R, U any] interface {
    Create(ctx context.Context, request C, response *R) error
    GetByID(ctx context.Context, id string, response *R) error
    List(ctx context.Context, queryScopes map[string]any, response *[]R) (int64, error)
    Update(ctx context.Context, id string, request U, response *R) error
    Delete(ctx context.Context, id string) error
}

type ServiceMapper[C, R, U, E any] interface {
    ToPostModel(dto C, model *E) error
    ToPatchModel(dto U, model *E) error
    ToDTO(model *E, dto *R) error
    GetModelKey(id string) (any, error)
    ApplyQueryScopes(queryScopes map[string]any) ([]func(*gorm.DB) *gorm.DB, error)
}
```

## Conflict Validation Flow

`service_impl.Create` handles `ErrConflictValidationRequired` from the repository:

```
repo.Create(ctx, &model)
  → ErrConflictValidationRequired (INSERT ignored, existing found)
      → mapper.ToDTO(existingModel, &response)
      → utils.IsSubSetInterface(request, response)
          → true  → transparent idempotent success
          → false → ErrConflictValidationFailed (caller returns 409)
```

The handler layer maps `ErrConflictValidationFailed` → HTTP 409.

## Callsite Map

| File | Role |
|---|---|
| `internal/ports/backend/services/service_impl.go` | Service[C,R,U,E] implementation |
| `internal/ports/backend/services/nested_service_impl.go` | NestedService[C,R,U,E] implementation |
| `pkg/ports/backend/services/services_pkg.go` | Public facades |
| `pkg/ports/backend/services/contracts/contracts.go` | Interface definitions |
| `demoapp/services/*.go` | All application service implementations |
