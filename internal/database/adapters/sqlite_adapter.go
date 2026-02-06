package adapters

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type sqliteAdapter struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func buildDSNFromPath(databasePath string) (string, error) {
	lp := strings.TrimSpace(databasePath)

	// memory by default if empty or contains "memory"
	if lp == "" || strings.Contains(lp, "memory") {
		lp = "memory"
	}

	if lp == "memory" || lp == ":memory:" {
		return "file::memory:?mode=memory&cache=shared", nil
	} else if strings.HasSuffix(strings.ToLower(lp), ".db") && !strings.Contains(lp, "file:") {
		absPath, err := filepath.Abs(lp)

		if err != nil {
			return "", fmt.Errorf("unable to get absolute path for database file %s: %w", lp, err)
		}

		dir := filepath.Dir(absPath)

		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("database directory does not exist: %s", dir)
			}
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
func NewSQLite(databasePath string, plugins ...gorm.Plugin) (dbcontracts.Database, error) {
	dsn, err := buildDSNFromPath(databasePath)
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open("sqlite", dsn)
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
