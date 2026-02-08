package database

import (
	"fmt"
	"strconv"

	pubcfg "github.com/brunojet/go-infra-backend/pkg/config"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

// NewDatabaseConfigFromEnv builds a DatabaseConfig using the standard env vars.
// It performs basic validation (endpoint format) but does not attempt to
// validate driver-specific semantics.
func NewDatabaseConfigFromEnv() (*dbcontracts.DatabaseConfig, error) {
	cfg := &dbcontracts.DatabaseConfig{
		Driver:   dbcontracts.DatabaseDriver(pubcfg.GetEnv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DatabaseDriverSQLiteMemory))),
		Endpoint: pubcfg.GetEnv(dbcontracts.DB_ENDPOINT_ENV, ":0"),
		Schema:   pubcfg.GetEnv(dbcontracts.DB_SCHEMA_ENV, ""),
		Name:     pubcfg.GetEnv(dbcontracts.DB_NAME_ENV, "app.db"),
	}

	if v := pubcfg.GetEnv(dbcontracts.DB_MAX_OPEN_CONNECTIONS_ENV, ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s=%q: %w", dbcontracts.DB_MAX_OPEN_CONNECTIONS_ENV, v, err)
		}
		cfg.MaxOpenConnections = n
	}
	if v := pubcfg.GetEnv(dbcontracts.DB_MAX_IDLE_CONNECTIONS_ENV, ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s=%q: %w", dbcontracts.DB_MAX_IDLE_CONNECTIONS_ENV, v, err)
		}
		cfg.MaxIdleConnections = n
	}
	if v := pubcfg.GetEnv(dbcontracts.DB_CONN_MAX_LIFETIME_ENV, ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s=%q: %w", dbcontracts.DB_CONN_MAX_LIFETIME_ENV, v, err)
		}
		cfg.ConnMaxLifetimeSecs = n
	}
	if v := pubcfg.GetEnv(dbcontracts.DB_CONN_MAX_IDLE_TIME_ENV, ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s=%q: %w", dbcontracts.DB_CONN_MAX_IDLE_TIME_ENV, v, err)
		}
		cfg.ConnMaxIdleTimeSecs = n
	}

	if err := pubcfg.ValidateEndpoint(cfg.Endpoint); err != nil {
		return nil, fmt.Errorf("invalid database endpoint: %w", err)
	}

	return cfg, nil
}
