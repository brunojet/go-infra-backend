package adapters

import (
	"database/sql"

	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// NewInMemory opens an in-memory sqlite DB and registers provided GORM plugins.
// Plugins are optional; callers can pass zero or more plugin instances.
func NewInMemory(plugins ...gorm.Plugin) (dbcontracts.Database, error) {
	sqlDB, err := sql.Open("sqlite", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	for _, p := range plugins {
		if p == nil {
			continue
		}
		if err := db.Use(p); err != nil {
			if sq, err2 := db.DB(); err2 == nil {
				_ = sq.Close()
			}
			return nil, err
		}
	}

	return &sqliteAdapter{sqlDB: sqlDB, gormDB: db}, nil
}

type sqliteAdapter struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func (s *sqliteAdapter) GormDB() *gorm.DB { return s.gormDB }

func (s *sqliteAdapter) Close() error {
	if s.gormDB != nil {
		if db, err := s.gormDB.DB(); err == nil && db != nil {
			return db.Close()
		}
	}
	if s.sqlDB != nil {
		return s.sqlDB.Close()
	}
	return nil
}
