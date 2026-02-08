package bootstrap

import (
	internalbootstrap "github.com/brunojet/go-infra-backend/internal/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

// NewSQLiteDatabaseWithObservability delegates to the internal bootstrap helper.
func NewSQLiteDatabaseWithObservability(databasePath string, sm contracts.ShutdownManager) (contracts.DatabaseAdapter, error) {
	// NOTE: contracts.DatabaseAdapter is an alias to internal/database/contracts.DatabaseAdapter.
	return internalbootstrap.NewSQLiteDatabaseWithObservability(databasePath, sm)
}
