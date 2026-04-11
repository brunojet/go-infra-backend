package handlers

import (
	"github.com/brunojet/go-infra-backend/demobff/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	"github.com/gin-gonic/gin"
)

// NewHelloWorldBffHandler registers the BFF-backed helloworld routes.
//
// Only Create, GetByID, and Delete are exposed:
//   - List is not registered: demoapp wraps its list response in {data, pagination},
//     which requires a custom decoder outside the standard BffRepository.List path.
//   - Update is not registered: BffRepository.Update uses PATCH but demoapp expects PUT.
func NewHelloWorldBffHandler(
	rg *gin.RouterGroup,
	s svccts.Service[services.HelloWorldBffDTO, services.HelloWorldBffDTO, services.HelloWorldBffDTO],
) {
	hp := handlers.HandlerParameters{
		HandlerPath:      "hello-worlds",
		IDValidationRule: handlers.Int64GtZero,
	}
	h := handlers.NewGenericHandler(hp, s)
	h.RegisterCollection(rg, "POST", h.Create)
	h.RegisterInstance(rg, "GET", h.GetByID)
	h.RegisterInstance(rg, "DELETE", h.Delete)
}
