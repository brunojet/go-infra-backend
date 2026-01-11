package database

import (
	"github.com/brunojet/go-infra-backend/internal/database"
	databasetypes "github.com/brunojet/go-infra-backend/internal/database/types"
)

type DatabaseParams = database.DatabaseParams

type DatabaseManager = database.DatabaseManager

var NewDatabaseManager = database.NewDatabaseManager

// Export NormalizeDBDriver function
var NormalizeDBDriver = databasetypes.NormalizeDBDriver

// Export DBDriver type and constants
type DBDriver = databasetypes.DBDriver

const (
	DBDriverMySQL  = databasetypes.DBDriverMySQL
	DBDriverSQLite = databasetypes.DBDriverSQLite
)
