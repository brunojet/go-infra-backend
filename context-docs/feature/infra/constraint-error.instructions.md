---
applyTo: "internal/ports/repositories/**,pkg/ports/repositories/**"
---

# Constraint Error — Architectural Context

## Purpose

Provide a single, driver-agnostic sentinel error for any database constraint violation
(UNIQUE, CHECK, FOREIGN KEY), normalizing raw database errors from different drivers
(SQLite, MySQL, …) into a consistent contract that application code can depend on.

## Problem Solved

GORM's `gorm.ErrDuplicatedKey` is only populated when the driver's `Translate()` method
can deserialize the raw error (e.g. via JSON marshal/unmarshal of the error code). The
`modernc/sqlite` driver (CGO-free) does not always do this, causing tests and service
code that checked `errors.Is(err, gorm.ErrDuplicatedKey)` to receive `false` even for
real UNIQUE violations.

## Implementation

### Sentinel

```go
// internal/ports/repositories/repository_errors.go
var ErrConstraintViolation = errors.New("constraint violation")
```

### MapDbError (classification point)

`MapDbError` is the single place that normalizes raw DB errors. The constraint branch
uses `fmt.Errorf` with multiple `%w` (Go 1.20+) to build a multi-error that satisfies
`errors.Is` for **both** sentinels simultaneously:

```go
} else if errors.Is(err, gorm.ErrDuplicatedKey) ||
    strings.Contains(err.Error(), "UNIQUE constraint failed") ||
    strings.Contains(err.Error(), "CHECK constraint failed") ||
    strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
    return fmt.Errorf("%w: %w", ErrConstraintViolation, err)
}
```

This means callers can use either:
- `errors.Is(err, ErrConstraintViolation)` — driver-agnostic, preferred
- `errors.Is(err, gorm.ErrDuplicatedKey)` — only if the driver translated it

### Public facade

```go
// pkg/ports/repositories/repositories_pkg.go
var ErrConstraintViolation = repositories.ErrConstraintViolation
```

## Usage Pattern

### In tests (model-level, raw GORM)
```go
// UPDATE hits the constraint directly — normalize before asserting
err := tx.Save(&entity).Error
require.ErrorIs(t, portsrepos.MapDbError(err), portsrepos.ErrConstraintViolation)
```

### In service/handler code
```go
if errors.Is(err, portsrepos.ErrConstraintViolation) {
    // return 409 Conflict
}
```

## Callsite Map

| File | Usage |
|---|---|
| `internal/ports/repositories/repository_errors.go` | definition + MapDbError |
| `pkg/ports/repositories/repositories_pkg.go` | public re-export |
| `demoapp/models/application_test.go` | UPDATE collision assertion |
| `demoapp/models/application_test.go` | AppConfiguration UPDATE collision assertion |
