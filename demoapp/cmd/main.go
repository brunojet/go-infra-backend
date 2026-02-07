package main

import (
	"log"
	"os"
	"strings"
	"time"

	demoapp "github.com/brunojet/go-infra-backend/demoapp/bootstrap"
	demoapprepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/bootstrap"
)

func main() {
	sm, stop := bootstrap.NewShutdownManagerWithSignals(10 * time.Second)
	defer stop()

	bootstrap.InitObservability(sm)

	databasePath := strings.TrimSpace(os.Getenv("DEMOAPP_SQLITE_PATH"))

	db, err := bootstrap.NewSQLiteDatabaseWithObservability(databasePath, sm)

	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	gormDb, err := db.GormDB()
	if err != nil {
		log.Fatalf("failed to get gorm DB: %v", err)
	}

	if err := gormDb.AutoMigrate(&demoapprepo.HelloWorld{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	httpServer := bootstrap.NewHttpServerWithObservability(sm)

	api := httpServer.Router.Group("/")

	if err := demoapp.SetupHelloWorldModule(db, api); err != nil {
		log.Fatalf("failed to setup demoapp module: %v", err)
	}

	if err := httpServer.StartAndWaitTermination(); err != nil {
		log.Printf("http server error: %v", err)
	}
}
