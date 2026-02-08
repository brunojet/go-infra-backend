package adapters

import (
	internaladapters "github.com/brunojet/go-infra-backend/internal/database/adapters"
	"github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

// NewSQLite delegates to the internal implementation.
func NewSQLite(databasePath string) (contracts.DatabaseAdapter, error) {
	return internaladapters.NewSQLite(databasePath)
}
