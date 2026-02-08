package adapters

import (
	internaladapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	db "github.com/brunojet/go-infra-backend/pkg/database"
)

// NewSQLite delegates to the internal implementation.
func NewSQLite(databasePath string) (db.DatabaseAdapter, error) {
	return internaladapters.NewSQLite(databasePath)
}
