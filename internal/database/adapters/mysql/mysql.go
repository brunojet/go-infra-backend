package mysql

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/brunojet/go-infra-backend/internal/database/contracts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrMissingDSN = errors.New("DB_DSN is required for mysql")

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

	if c.DSN == "" {
		return nil, ErrMissingDSN
	}

	gormDB, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Reasonable defaults for MySQL.
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

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
