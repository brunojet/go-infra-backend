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

type nestedGinHandler[E rpocontracts.Entity, D any] struct {
	php         *HandlerParameters
	svc         svccontracts.NestedService[D, E]
	basehandler hndcontracts.GenericHandler[E, D]
}

func NewGenericNestedHandler[E rpocontracts.Entity, D any](php *HandlerParameters, hp *HandlerParameters, s svccontracts.NestedService[D, E]) hndcontracts.NestedGenericHandler[E, D] {
	var baseHandler hndcontracts.GenericHandler[E, D]
	if svc, ok := s.(svccontracts.Service[D, E]); ok {
		baseHandler = NewGenericHandler[E](hp, svc)
	} else {
		panic("provided service does not implement Service[D, E]")
	}
	return &nestedGinHandler[E, D]{
		php:         php,
		svc:         s,
		basehandler: baseHandler, // Use the base service for non-nested operations
	}
}

// RegisterCollection registers collection-level routes for nested resources (e.g., /parent/:parentId/items)
func (h *nestedGinHandler[E, D]) RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	handlerNestedPath := strings.Trim(h.php.HandlerPath, "/")
	handlerPath := strings.Trim(h.basehandler.(*ginHandler[E, D]).hp.HandlerPath, "/")
	if handlerNestedPath == "" || handlerPath == "" {
		log.Default().Panic("HandlerPath cannot be empty")
	}
	fullPath := fmt.Sprintf("%s/:id/%s", handlerNestedPath, handlerPath)
	rg.Handle(strings.ToUpper(method), fullPath, handler)
}

func (h *nestedGinHandler[E, D]) RegisterInstance(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	h.basehandler.RegisterInstance(rg, method, handler)
}

func (h *nestedGinHandler[E, D]) CreateNested(c *gin.Context) {
	parentID, err := GetValidatedIDFromParam(c, "id", h.php.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	dto, err := BindJSONToDTOPtr[D](c)
	if err != nil {
		return
	}
	if err := h.svc.CreateNested(c.Request.Context(), parentID, dto); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto)
}

func (h *nestedGinHandler[E, D]) ListNested(c *gin.Context) {
	parentID, err := GetValidatedIDFromParam(c, "id", h.php.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
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

func (h *nestedGinHandler[E, D]) GetByID(c *gin.Context) {
	h.basehandler.GetByID(c)
}

func (h *nestedGinHandler[E, D]) Update(c *gin.Context) {
	h.basehandler.Update(c)
}

func (h *nestedGinHandler[E, D]) Delete(c *gin.Context) {
	h.basehandler.Delete(c)
}
