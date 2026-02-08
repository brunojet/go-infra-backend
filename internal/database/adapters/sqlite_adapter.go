package adapters

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type sqliteAdapter struct {
	gormDB *gorm.DB
}

func buildDSNFromPath(databasePath string) (string, error) {
	lp := strings.TrimSpace(databasePath)

	// memory by default if empty or contains "memory"
	if lp == "" || strings.Contains(lp, "memory") {
		return "file::memory:?mode=memory&cache=shared", nil
	} else if strings.HasSuffix(strings.ToLower(lp), ".db") && !strings.Contains(lp, "file:") {
		absPath, err := filepath.Abs(lp)

		if err != nil {
			return "", fmt.Errorf("unable to get absolute path for database file %s: %w", lp, err)
		}

		dir := filepath.Dir(absPath)

		if _, err := os.Stat(dir); err != nil {
			return "", fmt.Errorf("unable to stat database directory %s: %w", dir, err)
		}

		return "file:" + filepath.ToSlash(absPath), nil
	} else {
		return "", fmt.Errorf("unable to configure dsn with databasePath: %s", lp)
	}
}

// NewSQLite opens a sqlite DB at the provided path and registers provided GORM plugins.
// `dataSourceName` can be a filesystem path or a DSN like "file:my.db".
// Plugins are optional; callers can pass zero or more plugin instances.
func NewSQLite(databasePath string) (dbcontracts.DatabaseAdapter, error) {
	dsn, err := buildDSNFromPath(databasePath)
	if err != nil {
		return nil, err
	}

	sqlDb, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDb}, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &sqliteAdapter{gormDB: db}, nil
}

func (s *sqliteAdapter) GormDB() (*gorm.DB, error) {
	if s.gormDB == nil {
		return nil, fmt.Errorf("database connection is closed")
	}
	return s.gormDB, nil
}

func (s *sqliteAdapter) SqlDB() (*sql.DB, error) {
	if s.gormDB == nil {
		return nil, fmt.Errorf("database connection is closed")
	}
	return s.gormDB.DB()
}

func (s *sqliteAdapter) Close() error {
	if s.gormDB != nil {
		if sqlDb, err := s.gormDB.DB(); err != nil {
			return err
		} else if err := sqlDb.Close(); err != nil {
			return err
		}
		s.gormDB = nil
	}
	return nil
}
