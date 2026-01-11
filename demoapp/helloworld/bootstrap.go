package helloworld

import (
	"log"

	"github.com/brunojet/go-infra-backend/demoapp/core/repository"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/handlers"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/routers"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/services"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/testdata"
	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"gorm.io/gorm"
)

type HelloworldModule struct {
	Router *routers.HelloWorldRouter
}

func NewHelloworldModule(db *gorm.DB) *HelloworldModule {
	// Seed automático da tabela helloworld
	if err := testdata.SeedHelloWorld(db, "."); err != nil {
		log.Printf("seed helloworld: %v", err)
	}
	repo := repository.NewHelloWorldRepository(db)
	service := services.NewHelloWorldService(repo)
	handler := handlers.NewHelloWorldHandler(service)
	router := routers.NewHelloWorldRouter(handler)
	return &HelloworldModule{Router: router}
}

func (m *HelloworldModule) Register(r contracts.Router) {
	m.Router.Register(r)
}
