package handlers

import (
	helloworldRepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	helloworldSvc "github.com/brunojet/go-infra-backend/demoapp/services"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

func NewHelloWorldHandler(rg *gin.RouterGroup, s svcContracts.Service[helloworldSvc.HelloWorldDTO, helloworldRepo.HelloWorld]) {
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "hello-worlds",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericHandler[helloworldRepo.HelloWorld](handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.Create)
	handler.RegisterCollection(rg, "GET", handler.List)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}
