# pkg/

This folder is the **public API surface** intended to be imported by external applications.

Guiding principles:

- Keep the API **minimal and stable**.
- Prefer small facades (constructors + options) over exposing internal implementation details.
- Anything app-specific (composition/wiring, env bootstrap, lifecycle) stays under `internal/`.

Current module path:

- `github.com/brunojet/go-infra-backend/pkg/...`
