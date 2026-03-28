package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	"github.com/gin-gonic/gin"
)

func NewHelloWorldHandler(rg *gin.RouterGroup, s services.HelloWorldService) {
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "hello-worlds",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericHandler(handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.Create)
	handler.RegisterCollection(rg, "GET", handler.List)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}
