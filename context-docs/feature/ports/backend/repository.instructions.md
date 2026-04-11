---
applyTo: "internal/ports/backend/repositories/**,pkg/ports/backend/repositories/**"
---

# Ports — Repository Contracts

## Purpose

Define and document the repository port contracts and shared error sentinels used by all
adapter implementations (GORM, BFF, …). Application code depends only on these contracts
via the `pkg/` public facade — never on `internal/` directly.

## Package Map

```
internal/ports/backend/repositories/
    repository_impl.go          ← GenericRepository + GenericNestedRepository concrete impl (GORM)
    repository_errors.go        ← all error sentinels + MapDbError + MapTxError
    repository_impl_test.go     ← GORM repository integration tests
    repository_tx_test.go       ← transaction helper tests

pkg/ports/backend/repositories/
    repositories_pkg.go         ← public re-exports: NewRepository, NewNestedRepository,
                                    error sentinels, MapDbError, MapTxError
    repositories_pkg_test.go    ← facade smoke tests
    contracts/contracts.go      ← Repository[C,R,U], NestedRepository[C,R,U] interfaces
```

## Error Sentinels

All sentinels live in `internal/ports/backend/repositories/repository_errors.go` and are
re-exported via `pkg/ports/backend/repositories/repositories_pkg.go`.

| Sentinel | When emitted |
|---|---|
| `ErrNotFound` | Record not found (GORM `ErrRecordNotFound`) |
| `ErrConflictValidationRequired` | INSERT ignored (rows==0) but a conflicting record exists |
| `ErrConflictValidationFailed` | Service determined the conflict is a real duplicate (409) |
| `ErrCheckConstraintViolated` | INSERT ignored but no conflicting record found (edge case) |
| `ErrConstraintViolation` | Any UNIQUE / CHECK / FK constraint hit at DB level |

## ErrConstraintViolation — Constraint Normalization

**Problem**: GORM's `gorm.ErrDuplicatedKey` is only in the error chain when the driver
translates the raw error. The `modernc/sqlite` (CGO-free) driver sometimes does not,
causing `errors.Is(err, gorm.ErrDuplicatedKey) == false` for real UNIQUE violations.

**Solution**: `MapDbError` normalizes all constraint strings into a single sentinel using
Go 1.20+ multi-wrap:

```go
var ErrConstraintViolation = errors.New("constraint violation")

// In MapDbError:
} else if errors.Is(err, gorm.ErrDuplicatedKey) ||
    strings.Contains(err.Error(), "UNIQUE constraint failed") ||
    strings.Contains(err.Error(), "CHECK constraint failed") ||
    strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
    return fmt.Errorf("%w: %w", ErrConstraintViolation, err)
}
```

This means callers can use:
- `errors.Is(err, ErrConstraintViolation)` — driver-agnostic, preferred
- `errors.Is(err, gorm.ErrDuplicatedKey)` — only if the driver translated it

### Usage: UPDATE collision in tests
```go
err := tx.Save(&entity).Error
require.ErrorIs(t, portsrepos.MapDbError(err), portsrepos.ErrConstraintViolation)
```

### Usage: service/handler code
```go
if errors.Is(err, portsrepos.ErrConstraintViolation) {
    // return 409 Conflict
}
```

## CREATE Conflict Flow (MapTxError)

```
tx.Create(inOut)
  → RowsAffected == 0 → MapTxError returns ErrCheckConstraintViolated
  → if WhereOnConflict(tx).First(&existing) finds a record:
      → returns ErrConflictValidationRequired
  → service.Create receives ErrConflictValidationRequired
      → mapper.ToDTO(existing, &response)
      → IsSubSetInterface(request, response)
          → true  → transparent idempotent "created"
          → false → ErrConflictValidationFailed (409)
```

## Callsite Map

| File | Role |
|---|---|
| `internal/ports/backend/repositories/repository_errors.go` | All sentinel definitions + MapDbError + MapTxError |
| `pkg/ports/backend/repositories/repositories_pkg.go` | Public re-export of all sentinels + helpers |
| `demoapp/models/application_test.go` | UPDATE collision assertion |
| `demoapp/models/application_test.go` | AppConfiguration UPDATE collision assertion |
