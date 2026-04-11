package main

import (
	"log"
	"time"

	demoapp "github.com/brunojet/go-infra-backend/demoapp/bootstrap"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	demorepos "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap"
)

func main() {
	sm, stop := bootstrap.NewShutdownManagerWithSignals(10 * time.Second)
	defer stop()

	db, err := bootstrap.NewDatabaseWithObservability(sm)

	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	gormDb, err := db.GormDB()
	if err != nil {
		log.Fatalf("failed to get gorm DB: %v", err)
	}

	// AutoMigrate all domain models in dependency order (parents before children).
	if err := gormDb.AutoMigrate(
		&demorepos.HelloWorld{},
		&models.FilterType{},
		&models.Filter{},
		&models.TerminalModel{},
		&models.TerminalModelConfiguration{},
		&models.Application{},
		&models.ApplicationImage{},
		&models.ApplicationProfile{},
		&models.ApplicationProfileScreenshot{},
		&models.ApplicationConfiguration{},
		&models.ApplicationVersion{},
		&models.ApplicationCatalog{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	httpServer := bootstrap.NewHttpServerWithObservability(sm)

	api := httpServer.Router.Group("/")

	if err := demoapp.SetupHelloWorldModule(db, api); err != nil {
		log.Fatalf("failed to setup demoapp module: %v", err)
	}

	if err := demoapp.SetupStoreProviderModule(db, api); err != nil {
		log.Fatalf("failed to setup demoapp module: %v", err)
	}

	if err := httpServer.StartAndWaitTermination(); err != nil {
		log.Printf("http server error: %v", err)
	}
}
