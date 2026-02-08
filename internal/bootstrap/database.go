package bootstrap

import (
	"github.com/brunojet/go-infra-backend/internal/database"
	gormobs "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	bootcontracts "github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
)

func NewSQLiteDatabaseWithObservability(databasePath string, sm bootcontracts.ShutdownManager) (dbcontracts.DatabaseAdapter, error) {
	db, err := database.NewSQLiteDatabase(databasePath, gormobs.NewOtelGormPlugin())
	if err != nil {
		return nil, err
	}
	sm.Register("sqlite", db)
	return db.DatabaseAdapter(), nil
}
