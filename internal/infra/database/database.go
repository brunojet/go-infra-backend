package database

import (
	"context"
	"fmt"
	"strings"

	dbadpt "github.com/brunojet/go-infra-backend/internal/infra/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type databaseManagerImpl struct {
	db dbcontracts.DatabaseAdapter
}

func (d *databaseManagerImpl) DatabaseAdapter() dbcontracts.DatabaseAdapter {
	return d.db
}

func (d *databaseManagerImpl) Migrate(dst ...interface{}) error {
	gormDB, err := d.db.GormDB()
	if err != nil {
		return err
	}
	return gormDB.AutoMigrate(dst...)
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
	// Determine driver (default to configured default when empty)
	driver := cfg.Driver
	if driver == "" {
		driver = dbcontracts.DbDriverDefault
	}

	switch driver {
	case dbcontracts.DbDriverMySQL:
		// Build DSN using github.com/go-sql-driver/mysql Config helper.
		addr := strings.TrimSpace(cfg.Host)
		if cfg.Port != "" {
			addr = addr + ":" + strings.TrimSpace(cfg.Port)
		}
		mysqlCfg := mysql.Config{
			User:   cfg.UserName,
			Passwd: cfg.Password,
			Net:    "tcp",
			Addr:   addr,
			DBName: cfg.Name,
			Params: map[string]string{
				"parseTime": "true",
				"charset":   "utf8mb4",
				"collation": "utf8mb4_unicode_ci",
				"loc":       "Local",
			},
			InterpolateParams: true,
		}
		dsn := mysqlCfg.FormatDSN()
		return dbadpt.NewMySQL(dsn)
	case dbcontracts.DbDriverSQLite:
		// Prefer configured Mode for SQLite selection. This allows DB driver to
		// remain 'sqlite' while Mode selects between memory and disk.
		switch cfg.Mode {
		case dbcontracts.DatabaseModeMemory:
			if cfg.Name != "" {
				// allow user-provided name to be included in the in-memory DSN
				return dbadpt.NewSQLite("file:" + cfg.Name + "?mode=memory&cache=shared")
			}
			return dbadpt.NewSQLite(":memory:")
		case dbcontracts.DatabaseModeDisk:
			return dbadpt.NewSQLite(cfg.Name)
		default:
			return nil, dbcontracts.ErrUnsupportedDriver
		}
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
