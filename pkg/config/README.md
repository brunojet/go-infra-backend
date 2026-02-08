# pkg/config

Public configuration utilities.

Keep this small:

- env parsing/validation helpers that are useful across apps
- types that external apps need to reference

Do not expose app-specific environment variable names unless they are part of the library contract.

Currently:

- `GetEnv` and `GetEnvAsBool` are minimal wrappers over `internal/config` to keep external consumers away from `internal/` imports.
