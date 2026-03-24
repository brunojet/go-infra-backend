package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ginHandler[C any, R any, U any, E rpocontracts.Entity] struct {
	hp      *HandlerParameters
	service svccontracts.Service[C, R, U, E]
}

func NewGenericHandler[C any, R any, U any, E rpocontracts.Entity](hp *HandlerParameters, s svccontracts.Service[C, R, U, E]) hndcontracts.GenericHandler[E, C] {
	return &ginHandler[C, R, U, E]{hp: hp, service: s}
}

// RegisterCollection registers collection-level routes (e.g., /items)
func (h *ginHandler[C, R, U, E]) RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	handlerPath := strings.Trim(h.hp.HandlerPath, "/")
	if handlerPath == "" {
		log.Default().Panic("HandlerPath cannot be empty")
	}
	rg.Handle(strings.ToUpper(method), handlerPath, handler)
}

// RegisterInstance registers instance-level routes (e.g., /items/:id)
func (h *ginHandler[C, R, U, E]) RegisterInstance(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	handlerPath := strings.Trim(h.hp.HandlerPath, "/")
	if handlerPath == "" {
		log.Default().Panic("HandlerPath cannot be empty")
	}
	fullPath := handlerPath + "/:id"
	rg.Handle(strings.ToUpper(method), fullPath, handler)
}

func (h *ginHandler[C, R, U, E]) Create(c *gin.Context) {
	dtoPtr, err := BindJSONToDTOPtr[C](c)
	if err != nil {
		return
	}
	dto := *dtoPtr
	created, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ginHandler[C, R, U, E]) GetByID(c *gin.Context) {
	id, err := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	dto, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

func (h *ginHandler[C, R, U, E]) List(c *gin.Context) {
	params := BuildListParamsFromRequest(c)
	list, _, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)

}

func (h *ginHandler[C, R, U, E]) Update(c *gin.Context) {
	id, err := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	dtoPtr, err := BindJSONToDTOPtr[U](c)
	if err != nil {
		return
	}
	dto := *dtoPtr
	updated, err := h.service.Update(c.Request.Context(), id, dto)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ginHandler[C, R, U, E]) Delete(c *gin.Context) {
	id, err := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
