package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/brunojet/go-infra-backend/internal/database/contracts"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

const defaultInMemoryDSN = "file::memory:?cache=shared"

type Connector struct {
	DSN string

	mu    sync.Mutex
	gorm  *gorm.DB
	sqlDB *sql.DB
}

var _ contracts.DatabaseConnector = (*Connector)(nil)

func New(dsn string) *Connector {
	return &Connector{DSN: dsn}
}

func (c *Connector) Open() (*gorm.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.gorm != nil {
		return c.gorm, nil
	}

	dsn, pathForDir, okPath := normalizeDSN(c.DSN)
	if okPath {
		if err := ensureDirForPath(pathForDir); err != nil {
			return nil, err
		}
	}

	gormDB, err := gorm.Open(gormsqlite.Dialector{DriverName: "sqlite", DSN: dsn}, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// SQLite works best with a single open connection.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	c.gorm = gormDB
	c.sqlDB = sqlDB
	return gormDB, nil
}

func (c *Connector) DB() *gorm.DB {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gorm
}

func (c *Connector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sqlDB == nil {
		return nil
	}
	err := c.sqlDB.Close()
	c.sqlDB = nil
	c.gorm = nil
	return err
}

func normalizeDSN(raw string) (dsn string, pathForDir string, okPath bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return defaultInMemoryDSN, "", false
	}

	vl := strings.ToLower(v)
	if vl == ":memory:" || vl == "memory" || vl == "mem" || vl == "inmemory" || vl == "in-memory" {
		return defaultInMemoryDSN, "", false
	}

	// File URIs and explicit memory modes are passed through.
	if strings.HasPrefix(vl, "file:") || strings.Contains(vl, "mode=memory") {
		return v, "", false
	}

	base, _, cut := strings.Cut(v, "?")
	if cut {
		// Still a path, but without query params for dir creation.
		return v, base, true
	}
	return v, v, true
}

func ensureDirForPath(path string) error {
	// If path is relative or absolute, ensure its parent directory exists.
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
