# pkg/observability/httpmiddlewares

HTTP middleware helpers intended for external apps.

Only place things here if:

- the middleware is a supported part of the library product
- its config surface is stable and does not depend on internal bootstrap

Otherwise, keep middleware in `internal/` and expose only configuration hooks.
