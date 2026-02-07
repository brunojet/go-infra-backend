package database

import (
	"context"

	dbadpt "github.com/brunojet/go-infra-backend/internal/database/adapters"
	"github.com/brunojet/go-infra-backend/internal/database/contracts"
	"gorm.io/gorm"
)

// NewSQLiteDatabase creates an in-memory SQLite database and registers
// optional GORM plugins. It returns a dbcontracts.Database which must be
// closed when no longer needed.

type databaseManagerImpl struct {
	db contracts.DatabaseAdapter
}

func (d *databaseManagerImpl) DatabaseAdapter() contracts.DatabaseAdapter {
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

func NewDatabaseManager(adapter contracts.DatabaseAdapter, plugins ...gorm.Plugin) (contracts.DatabaseManager, error) {
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

func NewSQLiteDatabase(databasePath string, plugins ...gorm.Plugin) (contracts.DatabaseManager, error) {
	db, err := dbadpt.NewSQLite(databasePath)
	if err != nil {
		return nil, err
	}
	return NewDatabaseManager(db, plugins...)
}
