package database

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

type fakeConnector struct {
	openErr  error
	closeErr error
	opened   bool
	closed   bool
	returnDB *gorm.DB
}

func (f *fakeConnector) Open() (*gorm.DB, error) {
	f.opened = true
	return f.returnDB, f.openErr
}
func (f *fakeConnector) Close() error {
	f.closed = true
	return f.closeErr
}

func TestSelectFromParams_UnsupportedDriver(t *testing.T) {
	p := DatabaseParams{Driver: "invalid"}
	_, err := SelectFromParams(p)
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestNewDatabaseManager_Error(t *testing.T) {
	p := DatabaseParams{Driver: "invalid"}
	_, err := NewDatabaseManager(p)
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestDatabaseManager_Open_Close(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	db, err := mgr.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db != nil {
		// always nil in fake
	}
	if err := mgr.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}

func TestDatabaseManager_Close_Error(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{closeErr: errors.New("fail")}}
	err := mgr.Close()
	if err == nil {
		t.Fatal("expected error from Connector.Close")
	}
}

func TestDatabaseManager_openUnlocked_Cache(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	mgr.db = &gorm.DB{} // simulate already opened
	db, err := mgr.openUnlocked()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db != mgr.db {
		t.Fatal("should return cached db instance")
	}
}

func TestDatabaseManager_OpenAndMigrate_OpenUnlockedError(t *testing.T) {
	mgr := &DatabaseManager{Connector: nil}
	mgr.Params.Migrate = true
	_, err := mgr.OpenAndMigrate(func(db *gorm.DB) error { return nil })
	if err == nil {
		t.Fatal("expected error from openUnlocked")
	}
}

func TestDatabaseManager_OpenAndMigrate_MigratorError(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	mgr.Params.Migrate = true
	migrator := func(db *gorm.DB) error { return errors.New("migrate fail") }
	_, err := mgr.OpenAndMigrate(migrator)
	if err == nil || err.Error() != "migrate fail" {
		t.Fatalf("expected migrate fail error, got %v", err)
	}
}

func TestDatabaseManager_Open_Error(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{openErr: errors.New("fail")}}
	_, err := mgr.Open()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDatabaseManager_OpenAndMigrate(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	mgr.Params.Migrate = true
	called := false
	migrator := func(db *gorm.DB) error { called = true; return nil }
	_, err := mgr.OpenAndMigrate(migrator)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected migrator to be called")
	}
}

func TestDatabaseManager_OpenAndMigrate_NoMigrator(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	mgr.Params.Migrate = true
	_, err := mgr.OpenAndMigrate(nil)
	if !errors.Is(err, ErrMissingMigrator) {
		t.Fatalf("expected ErrMissingMigrator, got %v", err)
	}
}

func TestDatabaseManager_OpenAndMigrate_NoMigrateFlag(t *testing.T) {
	mgr := &DatabaseManager{Connector: &fakeConnector{}}
	mgr.Params.Migrate = false
	called := false
	migrator := func(db *gorm.DB) error { called = true; return nil }
	_, err := mgr.OpenAndMigrate(migrator)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("migrator should not be called when Migrate is false")
	}
}

func TestDatabaseManager_Open_NilConnector(t *testing.T) {
	mgr := &DatabaseManager{Connector: nil}
	_, err := mgr.Open()
	if err == nil {
		t.Fatal("expected error for nil connector")
	}
}

func TestSelect_SQLite(t *testing.T) {
	cfg := DatabaseConfig{Driver: "sqlite", DSN: "file::memory:"}
	conn, err := Select(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connector for sqlite")
	}
}

func TestSelect_MySQL(t *testing.T) {
	cfg := DatabaseConfig{Driver: "mysql", DSN: "some-dsn"}
	conn, err := Select(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connector for mysql")
	}
}

func TestSelect_Unsupported(t *testing.T) {
	cfg := DatabaseConfig{Driver: "invalid", DSN: ""}
	conn, err := Select(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
	if conn != nil {
		t.Fatal("expected nil connector for unsupported driver")
	}
}

func TestDatabaseManager_Close_NilConnector(t *testing.T) {
	mgr := &DatabaseManager{Connector: nil, db: &gorm.DB{}}
	err := mgr.Close()
	if err != nil {
		t.Fatalf("expected nil error when Connector is nil, got %v", err)
	}
}

func TestSelectFromParams_Success(t *testing.T) {
	p := DatabaseParams{Driver: "sqlite", DSN: "file::memory:"}
	conn, err := SelectFromParams(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connector for valid params")
	}
}

func TestNewDatabaseManager_Success(t *testing.T) {
	p := DatabaseParams{Driver: "sqlite", DSN: "file::memory:"}
	mgr, err := NewDatabaseManager(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mgr == nil || mgr.Connector == nil {
		t.Fatal("expected non-nil manager and connector")
	}
}
