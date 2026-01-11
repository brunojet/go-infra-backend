package main

import (
	"log"
	"net/http"

	"github.com/brunojet/go-infra-backend/demoapp/core"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld"

	infradatabase "github.com/brunojet/go-infra-backend/infra/database"
	infrahttp "github.com/brunojet/go-infra-backend/infra/http"
	infraobs "github.com/brunojet/go-infra-backend/infra/observability"
)

func main() {
	cfg := configFromEnv()

	dbMgr, err := infradatabase.NewDatabaseManager(cfg.Database)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	db, err := dbMgr.OpenAndMigrate(core.Register)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}

	defer func() {
		if err := dbMgr.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	obsMgr := infraobs.NewObservabilityManager(cfg.Observability)
	mws, err := obsMgr.Open()
	if err != nil {
		log.Fatalf("obs init: %v", err)
	}
	defer func() {
		if err := obsMgr.Close(); err != nil {
			log.Printf("close obs: %v", err)
		}
	}()

	httpParams := cfg.HTTP
	httpParams.ObservabilityMiddlewares = mws
	httpMgr, err := infrahttp.NewHTTPManager(httpParams)
	if err != nil {
		log.Fatalf("http init: %v", err)
	}

	// --- Helloworld Module ---
	hwModule := helloworld.NewHelloworldModule(db)

	httpRuntime, err := httpMgr.OpenAndRegister(func(r infrahttp.Router) error {
		r.GET("/health", func(c infrahttp.Context) {
			c.String(200, "ok")
		})
		hwModule.Register(r)
		return nil
	})
	if err != nil {
		log.Fatalf("http register: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           httpRuntime.Handler,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}
	log.Printf("server listening on %s", cfg.Server.Addr)
	log.Fatal(srv.ListenAndServe())
}
