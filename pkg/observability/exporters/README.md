# pkg/observability/exporters

Exporter implementations that are intended to be reusable across applications.

Rules of thumb:

- Exporters here should be generic and configurable via small option structs.
- If an exporter is opinionated to this repo's demo app or environment, keep it in `internal/`.
