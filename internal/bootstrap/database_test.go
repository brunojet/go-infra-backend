package bootstrap

import (
	"context"
	"path/filepath"
	"testing"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabaseWithObservability_RegistersShutdown(t *testing.T) {
	assert := assert.New(t)

	tmpDir := t.TempDir()
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLiteDisk))
	t.Setenv(dbcontracts.DB_NAME_ENV, filepath.Join(tmpDir, "app.db"))

	sm := NewShutdownManager(context.Background())
	db, err := NewDatabaseWithObservability(sm)
	assert.NoError(err)
	assert.NotNil(db)
	gdb, err := db.GormDB()
	assert.NoError(err)
	assert.NotNil(gdb)

	// Shutdown should execute the registered sqlite shutdown and return no error
	assert.NoError(sm.ShutdownWithTimeout(500))
}

func TestNewDatabaseWithObservability_InvalidEnvReturnsError(t *testing.T) {
	sm := NewShutdownManager(context.Background())

	// Use an unsupported driver to guarantee NewDatabaseManagerFromEnv fails.
	t.Setenv(dbcontracts.DB_DRIVER_ENV, "unsupported")

	db, err := NewDatabaseWithObservability(sm)
	require.Error(t, err)
	require.Nil(t, db)
}
