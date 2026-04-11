---
applyTo: "demoapp/models/**"
---

# Application Model Rules — Architectural Context

## Purpose

Document the model-layer business rules, hook behavior, and corresponding test
conventions for `demoapp/models`. This is a living reference — update as new
models or rules are added.

## Core Pattern: BeforeCreate + INSERT OR IGNORE

All models that support idempotent creation use:

```go
func (a MyModel) BeforeCreate(tx *gorm.DB) error {
    return repositories.AddOnConflictDoNothing(tx, ColA, ColB)
}
```

`AddOnConflictDoNothing` adds `ON CONFLICT (colA, colB) DO NOTHING` to the INSERT.
This means **a conflicting INSERT is silently ignored** — `err == nil`, `RowsAffected == 0`.

**The model layer never rejects a conflict with an error.** Rejection is responsibility
of the service layer via `ErrConflictValidationFailed`.

### Conflict resolution flow (repository layer)

```
tx.Create(inOut)
  → RowsAffected == 0 → MapTxError returns ErrNotFound
  → repo.Create calls getExistingWhenConflict
      → WhereOnConflict(tx).First(out)
      → found   → returns ErrConflictValidationRequired
      → not found → returns ErrCheckConstraintViolated (edge case)

service.Create receives ErrConflictValidationRequired
  → mapper.ToDTO(found_model, &response)
  → IsSubSetInterface(request, response)
      → true  → transparent "created" (idempotent)
      → false → ErrConflictValidationFailed (409)
```

### UPDATE behavior

`BeforeCreate` is **never called on UPDATE** (`tx.Save` with non-zero PK executes UPDATE).
UPDATE hits the UNIQUE constraint directly and returns the raw database error.
Callers must normalize via `MapDbError(err)` to get `ErrConstraintViolation`.

## Model Rules by Entity

### Application

| Rule | Implementation |
|---|---|
| `name` globally unique | `uniqueIndex:ux_application_name` |
| Same name + same customer → idempotent | `BeforeCreate` + `INSERT OR IGNORE ON CONFLICT (name)` |
| Same name + different customer → service rejects | `WhereOnConflict` returns existing; service compares customer_id |
| Rename to existing name → DB rejects | `ErrConstraintViolation` on UPDATE |

### ApplicationImage

| Rule | Implementation |
|---|---|
| `(application_id, file_hash, image_type)` unique | `uniqueIndex:idx_application_image_application` |
| `image_type` required | `BeforeCreate` validates `ImageType.Valid` before AddOnConflictDoNothing |
| `file_size` required | `NOT NULL` on column; **must be set in all test helpers** |
| Same hash + same type → idempotent | `BeforeCreate` + `INSERT OR IGNORE` |
| `GetOrCreate` | Model method; always call inside active transaction |

### ApplicationConfiguration

| Rule | Implementation |
|---|---|
| `(terminal_model_configuration_id, package_name)` unique | `uniqueIndex:ux_app_cfg_terminal_app` |
| Same package for same cfgTerm across apps → idempotent at model | `BeforeCreate` + `INSERT OR IGNORE` |
| Same package name can repeat for same app with different cfgTerms | Allowed by index design |
| Rename to conflicting package → DB rejects | `ErrConstraintViolation` on UPDATE |

## Test Conventions

### CREATE conflict → assert silent ignore
```go
err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
    return tx.Create(&second).Error, nil
})
require.NoError(t, err)                    // INSERT OR IGNORE → nil
require.Zero(t, second.PrimaryKey)         // struct not populated on ignore
// assert original record preserved via count or SELECT
```

### UPDATE conflict → normalize then assert
```go
err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
    return tx.Save(&entity).Error, nil
})
require.ErrorIs(t, portsrepos.MapDbError(err), portsrepos.ErrConstraintViolation)
```

### Test helper requirement
All `newApplicationImage(...)` calls and inline `&ApplicationImage{}` structs
**must include `FileSize: sql.NullInt64{Int64: 1024, Valid: true}`** — the field
is NOT NULL and SQLite will reject inserts without it.

## Pending

- [ ] Validate that the `Update`/`Patch` service path checks for name collision before
      hitting the DB (current gap: constraint only caught at DB level, not pre-validated)
