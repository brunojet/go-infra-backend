package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	svc "github.com/brunojet/go-infra-backend/demoapp/services"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	"github.com/gin-gonic/gin"
)

func NewFilterTypeHandler(rg *gin.RouterGroup, s svc.FilterTypeService) {
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "filter-types",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericHandler[models.FilterType](handlerParameters, s)
	handler.Register(rg, "POST", handler.Create)
	handler.Register(rg, "GET", handler.List)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewFiltersHandler(rg *gin.RouterGroup, s svc.FilterNestedService) {
	handlerNestedParameters := &hnd.HandlerParameters{
		HandlerPath:      "filter-types",
		IDValidationRule: hnd.Int64GtZero,
	}
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "filters",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericNestedHandler[models.Filter](handlerNestedParameters, handlerParameters, s)
	handler.RegisterNested(rg, "POST", handler.CreateNested)
	handler.RegisterNested(rg, "GET", handler.ListNested)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewTerminalModelHandler(rg *gin.RouterGroup, s svc.TerminalModelService) {
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "terminal-models",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericHandler[models.TerminalModel](handlerParameters, s)
	handler.Register(rg, "POST", handler.Create)
	handler.Register(rg, "GET", handler.List)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewTerminalModelConfigurationHandler(rg *gin.RouterGroup, s svc.TerminalModelConfigurationNestedService) {
	handlerNestedParameters := &hnd.HandlerParameters{
		HandlerPath:      "terminal-models",
		IDValidationRule: hnd.Int64GtZero,
	}
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "terminal-model-configurations",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericNestedHandler[models.TerminalModelConfiguration](handlerNestedParameters, handlerParameters, s)
	handler.RegisterNested(rg, "POST", handler.CreateNested)
	handler.RegisterNested(rg, "GET", handler.ListNested)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewApplicationHandler(rg *gin.RouterGroup, s svc.ApplicationService) {
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericHandler[models.Application](handlerParameters, s)
	handler.Register(rg, "POST", handler.Create)
	handler.Register(rg, "GET", handler.List)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewApplicationConfigurationHandler(rg *gin.RouterGroup, s svc.ApplicationConfigurationNestedService) {
	handlerNestedParameters := &hnd.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: hnd.Int64GtZero,
	}
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "application-configurations",
		IDValidationRule: hnd.Base64UrlSafe,
	}
	handler := hnd.NewGenericNestedHandler[models.ApplicationConfiguration](handlerNestedParameters, handlerParameters, s)
	handler.RegisterNested(rg, "POST", handler.CreateNested)
	handler.RegisterNested(rg, "GET", handler.ListNested)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewApplicationProfileHandler(rg *gin.RouterGroup, s svc.ApplicationProfileNestedService) {
	handlerNestedParameters := &hnd.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: hnd.Int64GtZero,
	}
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "application-profiles",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericNestedHandler[models.ApplicationProfile](handlerNestedParameters, handlerParameters, s)
	handler.RegisterNested(rg, "POST", handler.CreateNested)
	handler.RegisterNested(rg, "GET", handler.ListNested)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}

func NewApplicationVersionHandler(rg *gin.RouterGroup, s svc.ApplicationVersionNestedService) {
	handlerNestedParameters := &hnd.HandlerParameters{
		HandlerPath:      "application-configurations",
		IDValidationRule: hnd.Base64UrlSafe,
	}
	handlerParameters := &hnd.HandlerParameters{
		HandlerPath:      "application-versions",
		IDValidationRule: hnd.Int64GtZero,
	}
	handler := hnd.NewGenericNestedHandler[models.ApplicationVersion](handlerNestedParameters, handlerParameters, s)
	handler.RegisterNested(rg, "POST", handler.CreateNested)
	handler.RegisterNested(rg, "GET", handler.ListNested)
	handler.Register(rg, "GET", handler.GetByID)
	handler.Register(rg, "PUT", handler.Update)
	handler.Register(rg, "DELETE", handler.Delete)
}
