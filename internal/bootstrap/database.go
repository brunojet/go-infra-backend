package bootstrap

import (
	bootcontracts "github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/internal/database"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	gormobs "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
)

func NewSQLiteDatabaseWithObservability(databasePath string, sm bootcontracts.ShutdownManager) (dbcontracts.DatabaseAdapter, error) {
	db, err := database.NewSQLiteDatabase(databasePath, gormobs.NewOtelGormPlugin())
	if err != nil {
		return nil, err
	}
	sm.Register("sqlite", db)
	return db.DatabaseAdapter(), nil
}
