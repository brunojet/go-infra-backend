package bootstrap

import (
	gormobs "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	bootcontracts "github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/pkg/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

func NewDatabaseWithObservability(sm bootcontracts.ShutdownManager) (dbcontracts.DatabaseAdapter, error) {
	db, err := database.NewDatabaseManagerFromEnv(gormobs.NewOtelGormPlugin())
	if err != nil {
		return nil, err
	}
	sm.Register("database", db)
	return db.DatabaseAdapter(), nil
}
