package database

import (
	internaldb "github.com/brunojet/go-infra-backend/internal/database"

	"github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"gorm.io/gorm"
)

// NewDatabaseManager delegates to the internal implementation.
func NewDatabaseManager(adapter contracts.DatabaseAdapter, plugins ...gorm.Plugin) (contracts.DatabaseManager, error) {
	return internaldb.NewDatabaseManager(adapter, plugins...)
}

// NewSQLiteDatabase delegates to the internal implementation.
func NewSQLiteDatabase(databasePath string, plugins ...gorm.Plugin) (contracts.DatabaseManager, error) {
	return internaldb.NewSQLiteDatabase(databasePath, plugins...)
}
