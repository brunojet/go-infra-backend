package database

import (
	"testing"

	plugins "github.com/brunojet/go-infra-backend/internal/database/plugins"
)

// TestNewSQLiteDatabase_WithOtelPlugin ensures that `NewSQLiteDatabase` accepts
// the OpenTelemetry GORM plugin returned by `adapters.NewOtelGormPlugin`
// and that a usable DB is returned.
func TestNewSQLiteDatabase_WithOtelPlugin(t *testing.T) {
	p, err := plugins.NewOtelGormPlugin()
	if err != nil {
		t.Fatalf("failed to create otel gorm plugin: %v", err)
	}

	db, err := NewSQLiteDatabase(p)
	if err != nil {
		t.Fatalf("NewSQLiteDatabase returned error: %v", err)
	}
	if db == nil {
		t.Fatalf("expected non-nil database")
	}

	// Close and ensure no panic / error on Close.
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close database: %v", err)
	}
}
