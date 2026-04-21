package main

import (
	"log"
	"time"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/config"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient"
	"github.com/brunojet/go-infra-backend/pkg/infra/observability/httptransports"
)

func main() {
	sm, stop := bootstrap.NewShutdownManagerWithSignals(10 * time.Second)
	defer stop()

	if err := bootstrap.InitObservability(sm); err != nil {
		log.Fatalf("failed to setup observability: %v", err)
	}

	baseURL := config.GetEnv("DEMOAPP_BASE_URL", "http://localhost:8080")
	bffCfg := bffclient.BffClientConfig{
		BaseURL: baseURL,
		Timeout: 10 * time.Second,
	}
	client, _, err := bffclient.NewNetHttpAdapter(bffCfg, httptransports.NewOtelHttpTransport(nil))
	if err != nil {
		log.Fatalf("failed to create BFF client: %v", err)
	}
	_ = client

	httpServer := bootstrap.NewHttpServerWithObservability(sm,
		bffclient.BffHeadersMiddleware(bffCfg.HeadersProxy),
	)

	api := httpServer.Router.Group("/")
	_ = api

	// if err := demobff.SetupHelloWorldBffModule(client, api); err != nil {
	// 	log.Fatalf("failed to setup demobff module: %v", err)
	// }

	if err := httpServer.StartAndWaitTermination(); err != nil {
		log.Printf("http server error: %v", err)
	}
}
