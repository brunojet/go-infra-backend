package database

import (
	"context"
	"fmt"

	dbadpt "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"gorm.io/gorm"
)

type databaseManagerImpl struct {
	db dbcontracts.DatabaseAdapter
}

func (d *databaseManagerImpl) DatabaseAdapter() dbcontracts.DatabaseAdapter {
	return d.db
}

func (d *databaseManagerImpl) HealthCheck(ctx context.Context) error {
	sqlDB, err := d.db.SqlDB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (d *databaseManagerImpl) Shutdown(ctx context.Context) error {
	return d.db.Close()
}

func NewDatabaseManager(adapter dbcontracts.DatabaseAdapter, plugins ...gorm.Plugin) (dbcontracts.DatabaseManager, error) {
	if adapter == nil {
		return nil, fmt.Errorf("adapter is nil")
	}
	if len(plugins) > 0 {
		gormDb, err := adapter.GormDB()
		if err != nil {
			adapter.Close()
			return nil, err
		}

		for _, p := range plugins {
			if p == nil {
				continue
			}
			if err := gormDb.Use(p); err != nil {
				adapter.Close()
				return nil, err
			}
		}
	}

	return &databaseManagerImpl{db: adapter}, nil
}

func newDatabaseAdapterFromConfig(cfg *dbcontracts.DatabaseConfig) (dbcontracts.DatabaseAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	switch cfg.Driver {
	case dbcontracts.DbDriverSQLiteMemory:
		return dbadpt.NewSQLite(":memory:")
	case dbcontracts.DbDriverSQLiteDisk:
		return dbadpt.NewSQLite(cfg.Name)
	default:
		return nil, dbcontracts.ErrUnsupportedDriver
	}
}

func newDatabaseManagerFromConfig(cfg *dbcontracts.DatabaseConfig, plugins ...gorm.Plugin) (dbcontracts.DatabaseManager, error) {
	adapter, err := newDatabaseAdapterFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return NewDatabaseManager(adapter, plugins...)
}

func NewDatabaseManagerFromEnv(plugins ...gorm.Plugin) (dbcontracts.DatabaseManager, error) {
	cfg, err := newDatabaseConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return newDatabaseManagerFromConfig(cfg, plugins...)
}
