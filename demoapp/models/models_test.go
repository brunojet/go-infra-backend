package models

import (
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLiteMemory))

	db, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.Migrate(&Application{}, &TerminalModel{}, &TerminalModelConfiguration{})
	require.NoError(t, err)
}
