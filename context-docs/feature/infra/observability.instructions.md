---
applyTo: "internal/infra/observability/**,pkg/infra/observability/**"
---

# Observability — Architectural Context

## Purpose

Provide OpenTelemetry (OTel) tracing and metrics to all adapters through a consistent
**plugin-as-constructor-param** pattern. Cross-cutting concerns are never added inside
adapter method bodies — they are injected at construction time.

## Package Map

```
internal/infra/observability/
    adapters/otlp_tracer_wrapper.go         ← OTLP exporter setup (gRPC/HTTP)
    gorm_plugins/otel_gorm_plugin.go        ← GORM Plugin interface → db.Use(...)
    http_middlewares/otel_gin_middleware.go  ← Gin HandlerFunc → router.Use(...)
    http_transports/otelhttp_transport.go   ← http.RoundTripper wrapping otelhttp
    stats/stats.go                          ← Prometheus/OTel metric helpers

pkg/infra/observability/
    contracts/contracts.go                  ← TracerProvider, MeterProvider interfaces
    gormplugins/gormplugins_pkg.go          ← public facade: NewOtelGormPlugin()
    httpmiddlewares/httpmiddlewares_pkg.go  ← public facade: OtelGinMiddleware()
    httptransports/httptransports_pkg.go    ← public facade: NewOtelHttpTransport(inner)
```

## Plugin Pattern

Each adapter layer has a dedicated injection mechanism:

| Adapter | Mechanism | Example |
|---|---|---|
| GORM (database) | `gorm.Plugin` | `db.Use(gormplugins.NewOtelGormPlugin())` |
| Gin (HTTP server) | `gin.HandlerFunc` | `router.Use(httpmiddlewares.OtelGinMiddleware())` |
| net/http (BFF client) | `http.RoundTripper` | `bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))` |

All OTel plugins use the **global `TracerProvider`** registered during observability
bootstrap — no per-plugin configuration needed.

## Transport Composition Chain

```go
// OTel only (common case):
bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))

// OTel + token auth (planned):
bffclient.NewNetHttpAdapter(cfg,
    httptransports.NewOtelHttpTransport(
        bfftransports.NewTokenTransport(nil),
    ),
)
```

`nil` inner transport defaults to `http.DefaultTransport`.

## Bootstrap Integration

`internal/bootstrap/observability.go` wires the OTLP exporter and registers the global
`TracerProvider`. Application code never calls OTel SDK directly — always via the
plugin facades.

## Dependencies

| Module | Version | Role |
|---|---|---|
| `go.opentelemetry.io/otel` | current | SDK core |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | v0.68.0 | HTTP transport + middleware |
| `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin` | current | Gin middleware |
| `gorm.io/plugin/opentelemetry` | current | GORM plugin |
| `github.com/felixge/httpsnoop` | v1.0.4 | transitive (otelhttp) |

## Pending

- [ ] Token auth RoundTripper (`bfftransports.NewTokenTransport`) — see bff.instructions.md
