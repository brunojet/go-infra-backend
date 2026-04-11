package adapters

import (
	"database/sql"
	"fmt"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlAdapter struct {
	gormDB *gorm.DB
}

// NewMySQL opens a MySQL DB using the provided DSN and returns a DatabaseAdapter.
// The DSN should follow the go-sql-driver/mysql format, for example:
//
//	user:password@tcp(host:3306)/dbname?parseTime=true
func NewMySQL(dsn string) (dbcontracts.DatabaseAdapter, error) {
	sqlDb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		_ = sqlDb.Close()
		return nil, err
	}

	return &mysqlAdapter{gormDB: db}, nil
}

func (m *mysqlAdapter) GormDB() (*gorm.DB, error) {
	if m.gormDB == nil {
		return nil, fmt.Errorf("database connection is closed")
	}
	return m.gormDB, nil
}

func (m *mysqlAdapter) SqlDB() (*sql.DB, error) {
	if m.gormDB == nil {
		return nil, fmt.Errorf("database connection is closed")
	}
	return m.gormDB.DB()
}

func (m *mysqlAdapter) Close() error {
	if m.gormDB != nil {
		if sqlDb, err := m.gormDB.DB(); err != nil {
			return err
		} else if err := sqlDb.Close(); err != nil {
			return err
		}
		m.gormDB = nil
	}
	return nil
}
