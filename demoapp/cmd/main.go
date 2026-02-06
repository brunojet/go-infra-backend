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

	db := bootstrap.NewMySQLDatabaseWithObservability(databasePath, sm)

	if err := db.GormDB().AutoMigrate(&demoapprepo.HelloWorld{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	httpServer := bootstrap.NewHttpServerWithObservability(sm)

	api := httpServer.Router.Group("/")

	if err := demoapp.SetupHelloWorldModule(db.GormDB(), api); err != nil {
		log.Fatalf("failed to setup demoapp module: %v", err)
	}

	if err := httpServer.StartAndWaitTermination(); err != nil {
		log.Printf("http server error: %v", err)
	}
}
