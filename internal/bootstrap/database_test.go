package bootstrap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMySQLDatabaseWithObservability_MemoryDBRegistersShutdown(t *testing.T) {
	assert := assert.New(t)

	sm := NewShutdownManager(context.Background())
	db := NewMySQLDatabaseWithObservability("memory", sm)
	assert.NotNil(db)
	assert.NotNil(db.GormDB())

	// Shutdown should execute the registered sqlite shutdown and return no error
	assert.NoError(sm.ShutdownWithTimeout(500))
}
