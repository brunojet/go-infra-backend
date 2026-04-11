package main

import (
	"log"
	"time"

	demobff "github.com/brunojet/go-infra-backend/demobff/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/config"
	"github.com/brunojet/go-infra-backend/pkg/infra/observability/httptransports"
)

func main() {
	sm, stop := bootstrap.NewShutdownManagerWithSignals(10 * time.Second)
	defer stop()

	baseURL := config.GetEnv("DEMOAPP_BASE_URL", "http://localhost:8080")
	bffCfg := bffclient.BffClientConfig{
		BaseURL: baseURL,
		Timeout: 10 * time.Second,
	}
	client, _, err := bffclient.NewNetHttpAdapter(bffCfg, httptransports.NewOtelHttpTransport(nil))
	if err != nil {
		log.Fatalf("failed to create BFF client: %v", err)
	}

	httpServer := bootstrap.NewHttpServerWithObservability(sm)

	api := httpServer.Router.Group("/")

	if err := demobff.SetupHelloWorldBffModule(client, api); err != nil {
		log.Fatalf("failed to setup demobff module: %v", err)
	}

	if err := httpServer.StartAndWaitTermination(); err != nil {
		log.Printf("http server error: %v", err)
	}
}
