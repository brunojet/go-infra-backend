package demoapp

import (
	hnd "github.com/brunojet/go-infra-backend/demoapp/handlers"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	svc "github.com/brunojet/go-infra-backend/demoapp/services"
	"github.com/brunojet/go-infra-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func SetupHelloWorldModule(db database.DatabaseAdapter, rg *gin.RouterGroup) error {
	helloWorldRepo := repo.NewHelloWorldRepo(db)
	helloWorldService := svc.NewHelloWorldService(helloWorldRepo)
	helloWorldHandler := hnd.NewHelloWorldHandler(helloWorldService)
	helloWorldHandler.Register(rg, "POST", "/helloworlds/", helloWorldHandler.Create)
	helloWorldHandler.Register(rg, "GET", "/helloworlds/", helloWorldHandler.List)
	helloWorldHandler.Register(rg, "GET", "/helloworld/:id", helloWorldHandler.GetByID)
	helloWorldHandler.Register(rg, "PUT", "/helloworld/:id", helloWorldHandler.Update)
	helloWorldHandler.Register(rg, "DELETE", "/helloworld/:id", helloWorldHandler.Delete)
	return nil
}
