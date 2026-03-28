package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type nestedGinHandler[C, R, U any] struct {
	hndcontracts.GenericHandler[C, R, U]
	php HandlerParameters
	svc svccontracts.NestedService[C, R, U, rpocontracts.Entity]
}

func NewGenericNestedHandler[C, R, U any](php HandlerParameters, hp HandlerParameters, s svccontracts.NestedService[C, R, U, rpocontracts.Entity]) hndcontracts.NestedGenericHandler[C, R, U] {
	var baseHandler hndcontracts.GenericHandler[C, R, U]
	if svc, ok := s.(svccontracts.Service[C, R, U, rpocontracts.Entity]); ok {
		baseHandler = NewGenericHandler(hp, svc)
	} else {
		log.Default().Panic("provided service does not implement Service[D, D, D, E]")
	}
	return &nestedGinHandler[C, R, U]{
		GenericHandler: baseHandler, // Use the base service for non-nested operations
		php:            php,
		svc:            s,
	}
}

func (h *nestedGinHandler[C, R, U]) GetHandlerPath() string {
	nestedPath := GetHandlerPath(h.php)
	path := h.GenericHandler.GetHandlerPath()
	return fmt.Sprintf("%s/:id/%s", nestedPath, path)
}

// RegisterCollection registers collection-level routes for nested resources (e.g., /parent/:parentId/items)
func (h *nestedGinHandler[C, R, U]) RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	rg.Handle(strings.ToUpper(method), h.GetHandlerPath(), handler)
}

func (h *nestedGinHandler[C, R, U]) CreateNested(c *gin.Context) {
	parentID, ok := GetValidatedIDFromParam(c, "id", h.php.IDValidationRule)
	if !ok {
		return
	}
	var dto C
	if !BindJSONToDTO(c, &dto) {
		return
	}
	created, err := h.svc.CreateNested(c.Request.Context(), parentID, dto)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *nestedGinHandler[C, R, U]) ListNested(c *gin.Context) {
	parentID, ok := GetValidatedIDFromParam(c, "id", h.php.IDValidationRule)
	if !ok {
		return
	}
	params := BuildListParamsFromRequest(c)
	list, _, err := h.svc.ListNested(c.Request.Context(), parentID, params)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
