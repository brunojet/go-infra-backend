package bootstrap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMySQLDatabaseWithObservability_MemoryDBRegistersShutdown(t *testing.T) {
	assert := assert.New(t)

	sm := NewShutdownManager(context.Background())
	db, err := NewSQLiteDatabaseWithObservability("memory", sm)
	assert.NoError(err)
	assert.NotNil(db)
	assert.NotNil(db.GormDB())

	// Shutdown should execute the registered sqlite shutdown and return no error
	assert.NoError(sm.ShutdownWithTimeout(500))
}

func TestNewSQLiteDatabaseWithObservability_InvalidPathReturnsError(t *testing.T) {
	sm := NewShutdownManager(context.Background())

	// choose a path that adapter buildDSNFromPath rejects
	db, err := NewSQLiteDatabaseWithObservability("not_a_valid_dsn", sm)
	require.Error(t, err)
	require.Nil(t, db)
}
