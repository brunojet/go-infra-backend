package database

import (
	"github.com/brunojet/go-infra-backend/internal/database"
	"github.com/brunojet/go-infra-backend/internal/database/adapters"
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
	DB_HOST_ENV                 = dbcontracts.DB_HOST_ENV
	DB_PORT_ENV                 = dbcontracts.DB_PORT_ENV
	DB_SCHEMA_ENV               = dbcontracts.DB_SCHEMA_ENV
	DB_NAME_ENV                 = dbcontracts.DB_NAME_ENV
	DB_MAX_OPEN_CONNECTIONS_ENV = dbcontracts.DB_MAX_OPEN_CONNECTIONS_ENV
	DB_MAX_IDLE_CONNECTIONS_ENV = dbcontracts.DB_MAX_IDLE_CONNECTIONS_ENV
	DB_CONN_MAX_LIFETIME_ENV    = dbcontracts.DB_CONN_MAX_LIFETIME_ENV
	DB_CONN_MAX_IDLE_TIME_ENV   = dbcontracts.DB_CONN_MAX_IDLE_TIME_ENV
)

const (
	DatabaseDriverSQLite   = dbcontracts.DbDriverSQLite
	DatabaseDriverPostgres = dbcontracts.DbDriverPostgres
	DatabaseDriverMySQL    = dbcontracts.DbDriverMySQL
)

// ---- Constructors (delegating to internal) ----
func NewSQLite(databasePath string) (dbcontracts.DatabaseAdapter, error) {
	return adapters.NewSQLite(databasePath)
}

func NewDatabaseManager(adapter DatabaseAdapter, plugins ...gorm.Plugin) (DatabaseManager, error) {
	return database.NewDatabaseManager(adapter, plugins...)
}

func NewDatabaseManagerFromEnv(plugins ...gorm.Plugin) (DatabaseManager, error) {
	return database.NewDatabaseManagerFromEnv(plugins...)
}
