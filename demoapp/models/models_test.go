package models

import (
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	db, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.Migrate(&Application{}, &TerminalModel{}, &TerminalModelConfiguration{})
	require.NoError(t, err)

	err = db.Migrate(&AuditEvent{}, &AuditFieldChange{})
	require.NoError(t, err)
}
