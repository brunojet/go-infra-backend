# pkg/observability

Public observability building blocks.

Intent:

- External apps can import this package tree to configure logging/metrics/tracing.
- The public surface should be narrow: contracts + a few constructors/options.

Notes:

- Existing internal layout uses folders like `internal/observability/http_middlewares` and `internal/observability/gorm_plugins`.
- In `pkg/`, we prefer idiomatic package folder names without underscores:
  - `httpmiddlewares` (instead of `http_middlewares`)
  - `gormplugins` (instead of `gorm_plugins`)

During the transition we can duplicate code temporarily; later we consolidate and update imports.
