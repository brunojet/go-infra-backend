package handlers

import (
	"fmt"
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

func NewNestedGenericHandler[E rpocontracts.Entity, D any](hp *HandlerParameters, php *HandlerParameters, s svccontracts.NestedService[D, E]) hndcontracts.NestedGenericHandler[E, D] {
	return &nestedGinHandler[E, D]{
		php:         php,
		svc:         s,
		basehandler: NewGenericHandler[E](hp, s.(svccontracts.Service[D, E])), // Use the base service for non-nested operations
	}
}

func (h *nestedGinHandler[E, D]) RegisterNested(rg *gin.RouterGroup, method, parentPath, path string, handler gin.HandlerFunc) {
	fullPath := fmt.Sprintf("%s/:parentID/%s", strings.Trim(parentPath, "/"), strings.Trim(path, "/"))
	rg.Handle(strings.ToUpper(method), fullPath, handler)
}

func (h *nestedGinHandler[E, D]) CreateNested(c *gin.Context) {
	parentID, err := GetValidatedIDFromParam(c, "parentID", h.php.IDValidationRule)
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
	parentID, err := GetValidatedIDFromParam(c, "parentID", h.php.IDValidationRule)
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

func (h *nestedGinHandler[E, D]) Register(rg *gin.RouterGroup, method, path string, handler gin.HandlerFunc) {
	h.basehandler.Register(rg, method, path, handler)
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
