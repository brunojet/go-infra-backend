package dbtest

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/pkg/infra/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"gorm.io/gorm"
)

func randomDBName() string {
	return fmt.Sprintf("memorytestdb_%d", rand.Int63())
}

func RunWithSQLiteRetry(fn func() error) error {
	var err error
	for i := 0; i < 8; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		msg := strings.ToLower(err.Error())
		if !strings.Contains(msg, "database is locked") && !strings.Contains(msg, "database table is locked") {
			return err
		}
		time.Sleep(5 * time.Millisecond)
	}
	return err
}

func OpenMemoryDB(t *testing.T, migrate ...any) *gorm.DB {
	t.Helper()
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))
	// Set random memory DB name using URI to avoid collisions
	t.Setenv(dbcontracts.DB_NAME_ENV, randomDBName())

	dbm, err := database.NewDatabaseManagerFromEnv()
	if err != nil {
		t.Fatalf("failed to create database manager: %v", err)
	}
	if dbm == nil {
		t.Fatalf("database manager is nil")
	}

	if len(migrate) > 0 {
		if err := dbm.Migrate(migrate...); err != nil {
			t.Fatalf("failed to migrate test models: %v", err)
		}
	}

	gdb, err := dbm.DatabaseAdapter().GormDB()
	if err != nil {
		t.Fatalf("failed to get gorm db: %v", err)
	}
	return gdb
}
