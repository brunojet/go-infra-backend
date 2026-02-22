package database

import (
	"fmt"

	"github.com/brunojet/go-infra-backend/internal/config"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

// newDatabaseConfigFromEnv builds a DatabaseConfig using the standard env vars.
// It performs basic validation (endpoint format) but does not attempt to
// validate driver-specific semantics.
func newDatabaseConfigFromEnv() (*dbcontracts.DatabaseConfig, error) {
	cfg := &dbcontracts.DatabaseConfig{
		Driver:             dbcontracts.DatabaseDriver(config.GetEnv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverDefault))),
		Host:               config.GetEnv(dbcontracts.DB_HOST_ENV, ""),
		Port:               config.GetEnv(dbcontracts.DB_PORT_ENV, ""),
		Schema:             config.GetEnv(dbcontracts.DB_SCHEMA_ENV, ""),
		Name:               config.GetEnv(dbcontracts.DB_NAME_ENV, "app.db"),
		Mode:               dbcontracts.DatabaseMode(config.GetEnv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeDefault))),
		UserName:           config.GetEnv(dbcontracts.DB_USER_ENV, ""),
		Password:           config.GetEnv(dbcontracts.DB_PASS_ENV, ""),
		MaxOpenConnections: config.GetEnvAsInt(dbcontracts.DB_MAX_OPEN_CONNECTIONS_ENV, dbcontracts.DbMaxOpenConnectionsDefault),
		MaxIdleConnections: config.GetEnvAsInt(dbcontracts.DB_MAX_IDLE_CONNECTIONS_ENV, dbcontracts.DbMaxIdleConnectionsDefault),
		ConnMaxLifetimeSecs: config.GetEnvAsInt(dbcontracts.DB_CONN_MAX_LIFETIME_ENV,
			dbcontracts.DbConnMaxLifetimeSecsDefault),
		ConnMaxIdleTimeSecs: config.GetEnvAsInt(dbcontracts.DB_CONN_MAX_IDLE_TIME_ENV,
			dbcontracts.DbConnMaxIdleTimeSecsDefault),
	}

	switch cfg.Driver {
	case dbcontracts.DbDriverSQLite:
		if cfg.Mode == dbcontracts.DatabaseModeDisk && cfg.Name == "" {
			return nil, fmt.Errorf("invalid database name: cannot be empty for disk-based SQLite")
		}
	case dbcontracts.DbDriverPostgres, dbcontracts.DbDriverMySQL:
		if err := config.ValidateHost(cfg.Host); err != nil {
			return nil, fmt.Errorf("invalid database host: %w", err)
		}
		if err := config.ValidatePort(cfg.Port); err != nil {
			return nil, fmt.Errorf("invalid database port: %w", err)
		}
	default:
		return nil, fmt.Errorf("%w: %q", dbcontracts.ErrUnsupportedDriver, cfg.Driver)
	}

	return cfg, nil
}
