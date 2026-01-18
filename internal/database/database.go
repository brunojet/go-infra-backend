package database

import (
	dbadpt "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	"gorm.io/gorm"
)

// NewSQLiteDatabase creates an in-memory SQLite database and registers
// optional GORM plugins. It returns a dbcontracts.Database which must be
// closed when no longer needed.
func NewSQLiteDatabase(plugins ...gorm.Plugin) (dbcontracts.Database, error) {
	return dbadpt.NewInMemory(plugins...)
}
