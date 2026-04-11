package adapters

import (
	internaladapters "github.com/brunojet/go-infra-backend/internal/infra/database/adapters"
	db "github.com/brunojet/go-infra-backend/pkg/infra/database"
)

// NewMySQL delegates to the internal implementation.
func NewMySQL(dsn string) (db.DatabaseAdapter, error) {
	return internaladapters.NewMySQL(dsn)
}
