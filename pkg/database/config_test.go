package database

import "testing"

func TestNewDatabaseConfigFromEnv_HappyPath(t *testing.T) {
	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverSQLiteMemory))
	t.Setenv(DB_ENDPOINT_ENV, ":0")
	t.Setenv(DB_SCHEMA_ENV, "")
	t.Setenv(DB_NAME_ENV, "app.db")

	cfg, err := NewDatabaseConfigFromEnv()
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if cfg == nil {
		t.Fatalf("expected non-nil cfg")
	}
}

func TestNewDatabaseConfigFromEnv_InvalidEndpoint(t *testing.T) {
	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverSQLiteMemory))
	t.Setenv(DB_ENDPOINT_ENV, ":99999")

	cfg, err := NewDatabaseConfigFromEnv()
	if err == nil {
		t.Fatalf("expected error")
	}
	if cfg != nil {
		t.Fatalf("expected nil cfg")
	}
}
