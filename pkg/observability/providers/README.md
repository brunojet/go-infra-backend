# pkg/observability/providers

Provider/factory layer.

This is usually where external apps want to start:

- `New...Provider(...)` / `Build...(...)` constructors
- accepts `contracts` interfaces and option structs

Avoid leaking implementation types from SDKs unless necessary.
