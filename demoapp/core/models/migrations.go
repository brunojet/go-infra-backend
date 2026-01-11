package models

// MigratableModels returns the list of GORM models that should be auto-migrated.
// Keep this list in core so the server wiring can remain simple.
func MigratableModels() []any {
	return []any{
		&HelloWorld{},
	}
}
