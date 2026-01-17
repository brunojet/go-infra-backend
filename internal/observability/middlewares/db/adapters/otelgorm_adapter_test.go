package adapters_test

import (
	"database/sql"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/middlewares/db/adapters"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func TestNewOtelGormPlugin_RegistersOnDB(t *testing.T) {
	p, err := adapters.NewOtelGormPlugin()
	if err != nil {
		t.Fatalf("unexpected error creating plugin: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}

	// open an in-memory sqlite DB using the modernc pure-Go driver
	sqlDB, err := sql.Open("sqlite", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("failed opening sql DB with modernc driver: %v", err)
	}
	defer sqlDB.Close()

	// create a gorm DB from the *sql.DB so GORM uses modernc under the hood
	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed opening gorm DB from sql.DB: %v", err)
	}

	if err := db.Use(p); err != nil {
		t.Fatalf("db.Use returned error: %v", err)
	}
}
