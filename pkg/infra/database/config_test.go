package database

import (
	"context"
	"path/filepath"
	"testing"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/stretchr/testify/require"
)

func TestNewDatabaseManagerFromEnv_HappyPath_SQLiteDisk(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeDisk))
	t.Setenv(DB_NAME_ENV, filepath.Join(tmpDir, "app.db"))
	t.Setenv(DB_SCHEMA_ENV, "")

	db, err := NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)

	require.NoError(t, db.HealthCheck(context.Background()))
	require.NoError(t, db.Shutdown(context.Background()))
}

func TestNewDatabaseManagerFromEnv_InvalidPort_Postgres(t *testing.T) {
	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverPostgres))
	t.Setenv(DB_HOST_ENV, "localhost")
	t.Setenv(DB_PORT_ENV, "99999")

	db, err := NewDatabaseManagerFromEnv()
	require.Error(t, err)
	require.Nil(t, db)
}

func TestNewDatabaseManagerFromEnv_UnsupportedDriver(t *testing.T) {
	t.Setenv(DB_DRIVER_ENV, "unsupported")

	db, err := NewDatabaseManagerFromEnv()
	require.Error(t, err)
	require.Nil(t, db)
}

func TestNewDatabaseManagerFromEnv_SQLiteMemory_Unsupported(t *testing.T) {
	// New design: use DB_DRIVER=sqlite and DB_MODE=memory to select in-memory sqlite.
	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	db, err := NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)
}

func TestNewDatabaseManagerFromEnv_SQLiteDisk_MemoryDSN(t *testing.T) {
	// Current implementation supports in-memory sqlite via sqlite_disk + Name containing "memory".
	t.Setenv(DB_DRIVER_ENV, string(DatabaseDriverSQLite))
	t.Setenv(DB_NAME_ENV, "memory")

	db, err := NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NoError(t, db.Shutdown(context.Background()))
}

func TestEnvConstants_AreFromContracts(t *testing.T) {
	// Guardrail: make sure the pkg/database re-exports match contracts.
	require.Equal(t, dbcontracts.DB_DRIVER_ENV, DB_DRIVER_ENV)
	require.Equal(t, dbcontracts.DB_HOST_ENV, DB_HOST_ENV)
	require.Equal(t, dbcontracts.DB_PORT_ENV, DB_PORT_ENV)
}
