package database

import (
	internaldb "github.com/brunojet/go-infra-backend/internal/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"gorm.io/gorm"
)

// ---- Contracts (re-exported for convenience) ----

type DatabaseAdapter = dbcontracts.DatabaseAdapter

type DatabaseManager = dbcontracts.DatabaseManager

type DatabaseDriver = dbcontracts.DatabaseDriver

type DatabaseConfig = dbcontracts.DatabaseConfig

const (
	DB_DRIVER_ENV               = dbcontracts.DB_DRIVER_ENV
	DB_ENDPOINT_ENV             = dbcontracts.DB_ENDPOINT_ENV
	DB_SCHEMA_ENV               = dbcontracts.DB_SCHEMA_ENV
	DB_NAME_ENV                 = dbcontracts.DB_NAME_ENV
	DB_MAX_OPEN_CONNECTIONS_ENV = dbcontracts.DB_MAX_OPEN_CONNECTIONS_ENV
	DB_MAX_IDLE_CONNECTIONS_ENV = dbcontracts.DB_MAX_IDLE_CONNECTIONS_ENV
	DB_CONN_MAX_LIFETIME_ENV    = dbcontracts.DB_CONN_MAX_LIFETIME_ENV
	DB_CONN_MAX_IDLE_TIME_ENV   = dbcontracts.DB_CONN_MAX_IDLE_TIME_ENV
)

const (
	DatabaseDriverSQLiteMemory = dbcontracts.DatabaseDriverSQLiteMemory
	DatabaseDriverSQLiteDisk   = dbcontracts.DatabaseDriverSQLiteDisk
	DatabaseDriverPostgres     = dbcontracts.DatabaseDriverPostgres
	DatabaseDriverMySQL        = dbcontracts.DatabaseDriverMySQL
)

// ---- Constructors (delegating to internal) ----

// NewDatabaseManager delegates to the internal implementation.
func NewDatabaseManager(adapter DatabaseAdapter, plugins ...gorm.Plugin) (DatabaseManager, error) {
	return internaldb.NewDatabaseManager(adapter, plugins...)
}

// NewSQLiteDatabase delegates to the internal implementation.
func NewSQLiteDatabase(databasePath string, plugins ...gorm.Plugin) (DatabaseManager, error) {
	return internaldb.NewSQLiteDatabase(databasePath, plugins...)
}
