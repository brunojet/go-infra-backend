package bootstrap

import (
	"context"
	"log"

	bootcontracts "github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/internal/database"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	gormobs "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	"gorm.io/gorm"
)

func otelGormPlugin() gorm.Plugin {
	plugin, err := gormobs.NewOtelGormPlugin()
	if err != nil {
		log.Fatalf("failed to create otel gorm plugin: %v", err)
	}
	return plugin
}

func NewMySQLDatabaseWithObservability(databasePath string, sm bootcontracts.ShutdownManager) dbcontracts.Database {
	plugin := otelGormPlugin()

	db, err := database.NewSQLiteDatabase(databasePath, plugin)

	if err != nil {
		log.Fatalf("failed to open sqlite database (%s): %v", databasePath, err)
	}

	shutdown := func(ctx context.Context) error {
		db.Close()
		return nil
	}

	sm.RegisterFunc("sqlite", shutdown)

	return db
}
