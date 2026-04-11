---
applyTo: "demobff/**"
---

# demobff — Application Context

## Purpose

`demobff` is a BFF (Backend For Frontend) application that proxies the `demoapp` backend
via HTTP using the `pkg/bffclient` + `pkg/ports/bff` infrastructure.

It demonstrates the full BFF stack end-to-end:
`BffClient → BffRepository → BffService → Handler`

Application profile: **app** layer — consumes `pkg/` contracts only, never imports `internal/`.

## Package Map

```
demobff/
    bootstrap/bootstrap.go      ← SetupHelloWorldBffModule(client, rg)
    cmd/main.go                 ← entry point: BffClient + HttpServer wiring
    handlers/
        helloworld_bff_hnd.go   ← registers POST /hello-worlds, GET/DELETE /hello-worlds/:id
    services/
        helloworld_bff_svc.go   ← HelloWorldBffE, HelloWorldBffDTO, HelloWorldBffMapper
```

## Wire-up

```
cmd/main.go
  ├─ config.GetEnv("DEMOAPP_BASE_URL", "http://localhost:8080")
  ├─ bffclient.NewNetHttpAdapter(cfg, httptransports.NewOtelHttpTransport(nil))
  ├─ bootstrap.NewHttpServerWithObservability(sm)
  └─ SetupHelloWorldBffModule(client, api)
       ├─ bffrepo.NewBffRepository[E, E, E](client)
       ├─ bffsvc.NewBffServiceImpl[DTO, DTO, DTO, E, E, E](repo, HelloWorldBffMapper{})
       └─ hnd.NewHelloWorldBffHandler(rg, service)
```

## Type Conventions

| Type | Role |
|---|---|
| `HelloWorldBffE` | Upstream DTO (CE = RE = UE, symmetric API) |
| `HelloWorldBffDTO` | Domain DTO exposed to handlers |
| `HelloWorldBffMapper` | Pass-through mapper (identical domain and upstream shapes) |

## Registered Routes

| Method | Path | Notes |
|---|---|---|
| POST | /hello-worlds | Create |
| GET | /hello-worlds/:id | GetByID |
| DELETE | /hello-worlds/:id | Delete |

## Known Limitations (Pending)

- **List not exposed**: demoapp wraps its list response in `{data, pagination}` — a JSON
  object, not an array. `BffRepository.List` decodes into `*[]RE` (expects a flat array)
  so this format mismatch causes a decode error. Resolution options:
  (a) add a flat-list endpoint to demoapp, or
  (b) introduce a `BffResponseUnwrapper` extension point in `BffRepository`.

- **Update not exposed**: `BffRepository.Update` issues `PATCH` but demoapp registers
  the Update route as `PUT`. Resolution: either align demoapp to PATCH or add a PUT
  variant to `BffRepository`.

## Configuration

| Env var | Default | Description |
|---|---|---|
| `DEMOAPP_BASE_URL` | `http://localhost:8080` | demoapp upstream base URL |
| `HTTP_ADDR` | `:8080` | demobff listen address (from `internal/bootstrap`) |
