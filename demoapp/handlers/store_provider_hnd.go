package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	"github.com/gin-gonic/gin"
)

func NewFilterTypeHandler(rg *gin.RouterGroup, s services.FilterTypeService) {
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "filter-types",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericHandler(handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.Create)
	handler.RegisterCollection(rg, "GET", handler.List)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewFiltersHandler(rg *gin.RouterGroup, s services.FilterNestedService) {
	handlerNestedParameters := handlers.HandlerParameters{
		HandlerPath:      "filter-types",
		IDValidationRule: handlers.Int64GtZero,
	}
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "filters",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericNestedHandler(handlerNestedParameters, handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.CreateNested)
	handler.RegisterCollection(rg, "GET", handler.ListNested)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewTerminalModelHandler(rg *gin.RouterGroup, s services.TerminalModelService) {
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "terminal-models",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericHandler(handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.Create)
	handler.RegisterCollection(rg, "GET", handler.List)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewTerminalModelConfigurationHandler(rg *gin.RouterGroup, s services.TerminalModelConfigurationNestedService) {
	handlerNestedParameters := handlers.HandlerParameters{
		HandlerPath:      "terminal-models",
		IDValidationRule: handlers.Int64GtZero,
	}
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "terminal-model-configurations",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericNestedHandler(handlerNestedParameters, handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.CreateNested)
	handler.RegisterCollection(rg, "GET", handler.ListNested)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewApplicationHandler(rg *gin.RouterGroup, s services.ApplicationService) {
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericHandler(handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.Create)
	handler.RegisterCollection(rg, "GET", handler.List)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewApplicationConfigurationHandler(rg *gin.RouterGroup, s services.ApplicationConfigurationNestedService) {
	handlerNestedParameters := handlers.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: handlers.Int64GtZero,
	}
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "application-configurations",
		IDValidationRule: handlers.Base64UrlSafe,
	}
	handler := handlers.NewGenericNestedHandler(handlerNestedParameters, handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.CreateNested)
	handler.RegisterCollection(rg, "GET", handler.ListNested)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewApplicationProfileHandler(rg *gin.RouterGroup, s services.ApplicationProfileNestedService) {
	handlerNestedParameters := handlers.HandlerParameters{
		HandlerPath:      "applications",
		IDValidationRule: handlers.Int64GtZero,
	}
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "application-profiles",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericNestedHandler(handlerNestedParameters, handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.CreateNested)
	handler.RegisterCollection(rg, "GET", handler.ListNested)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}

func NewApplicationVersionHandler(rg *gin.RouterGroup, s services.ApplicationVersionNestedService) {
	handlerNestedParameters := handlers.HandlerParameters{
		HandlerPath:      "application-configurations",
		IDValidationRule: handlers.Base64UrlSafe,
	}
	handlerParameters := handlers.HandlerParameters{
		HandlerPath:      "application-versions",
		IDValidationRule: handlers.Int64GtZero,
	}
	handler := handlers.NewGenericNestedHandler(handlerNestedParameters, handlerParameters, s)
	handler.RegisterCollection(rg, "POST", handler.CreateNested)
	handler.RegisterCollection(rg, "GET", handler.ListNested)
	handler.RegisterInstance(rg, "GET", handler.GetByID)
	handler.RegisterInstance(rg, "PUT", handler.Update)
	handler.RegisterInstance(rg, "DELETE", handler.Delete)
}
