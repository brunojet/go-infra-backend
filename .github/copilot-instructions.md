# GitHub Copilot — Project Instructions

## Architecture

This project follows the **Ports & Adapters (Hexagonal)** pattern with three distinct layers:

| Layer | Location | Responsibility |
|---|---|---|
| **infra** | `internal/` | Adapter implementations (net/http, GORM, OTel, etc.) — never imported directly by app code |
| **ports** | `pkg/ports/`, `pkg/` | Public port contracts (interfaces) and public facades re-exporting internal constructors |
| **app** | `demoapp/`, `demobff/`, … | Application layer: handlers, services, repositories, models, DTOs — consumes only `pkg/` contracts |

External dependencies are **never imported directly** from application packages (`demoapp/`, `demobff/`, etc.). Application code depends only on `pkg/` contracts.

## Context Documents

Detailed architectural decisions, patterns, and current implementation state are maintained in `context-docs/`.
Context documents are organized to mirror the three-layer architecture:

| Category | Folder | Contents |
|---|---|---|
| **infra** | `context-docs/feature/infra/` | Infrastructure adapters: transport, database, observability |
| **ports** | `context-docs/feature/ports/` | Reusable port contracts and public facades |
| **app** | `context-docs/feature/app/` | Application-level features consuming infra and ports |

When working in a specific subsystem, **read the corresponding context document before making changes**:

| Area | Context file |
|---|---|
| BFF / upstream HTTP layer (infra) | [context-docs/feature/infra/bff.instructions.md](../context-docs/feature/infra/bff.instructions.md) |
| Observability (OTel, metrics, logging) | [context-docs/observability-contracts.md](../context-docs/observability-contracts.md) |
| Coverage directives | [context-docs/diretivas-cobertura-go.md](../context-docs/diretivas-cobertura-go.md) |

## Observability plugin pattern

The project uses a consistent **plugin-as-constructor-param** pattern across all adapters:

| Layer | Plugin mechanism |
|---|---|
| GORM (database) | `gorm.Plugin` → `db.Use(gormplugins.NewOtelGormPlugin())` |
| Gin (HTTP server) | `gin.HandlerFunc` → `router.Use(httpmiddlewares.OtelGinMiddleware())` |
| net/http (BFF client) | `http.RoundTripper` → `bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))` |

Never add cross-cutting concerns (tracing, auth, logging) inside adapter method bodies. Always inject via the plugin mechanism.

## Naming conventions

- `Upstream*` / `*E` suffix → upstream/legacy API DTOs (e.g. `CE`, `RE`, `UE`)
- `Domain*` / plain suffix → internal domain types (e.g. `C`, `R`, `U`)
- `*Impl` → concrete service/repository implementations
- `*Pkg` or `*_pkg.go` → public facade files in `pkg/`

## Testing

- Use `-tags=debug` for all test runs (build tag required for `debugassert`)
- `httptest.NewServer` for BFF adapter tests — real HTTP round-trip, no mocks
- Test files live alongside the code they test (`*_test.go` in same package)

## Context Document Lifecycle

Context documents in `context-docs/` are living references — keep them accurate.

- **At conversation start**: Agree on a feature name. Read the corresponding context document before making any changes.
- **During implementation**: Update the context document as contracts, file locations, and structural decisions are finalized — do not defer to the end.
- **At implementation completion**: Sanitize the document — correct stale file paths, remove historical "Completed" blocks (completed items are already reflected in the Package Map and code), and ensure the Pending list matches actual remaining work.

Never leave a context document with paths or contracts that no longer match the codebase.
