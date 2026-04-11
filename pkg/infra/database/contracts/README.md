# pkg/database/contracts

Minimal stable database interfaces and shared types.

Examples of what belongs here:

- repository/transaction interfaces consumed by external apps
- error types that callers need to check
- option/config structs used by public constructors

Avoid:

- GORM concrete types in the API, unless explicitly intended.
