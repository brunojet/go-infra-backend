# pkg/database/adapters/sqlite

SQLite adapter placeholder.

Whether this becomes public depends on your target consumers:

- If external apps should be able to use SQLite quickly (local/dev/test), we can promote a supported adapter here.
- If SQLite is only for this repo's tests/demos, it should remain in `internal/`.
