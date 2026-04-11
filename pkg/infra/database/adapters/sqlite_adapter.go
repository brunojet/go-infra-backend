package adapters

import (
	internaladapters "github.com/brunojet/go-infra-backend/internal/infra/database/adapters"
	db "github.com/brunojet/go-infra-backend/pkg/infra/database"
)

// NewSQLite delegates to the internal implementation.
func NewSQLite(databasePath string) (db.DatabaseAdapter, error) {
	return internaladapters.NewSQLite(databasePath)
}
