package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	hndcts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ginHandler[C, R, U any] struct {
	hp      HandlerParameters
	service svccts.Service[C, R, U, rpocts.Entity]
}

func NewGenericHandler[C, R, U any](hp HandlerParameters, s svccts.Service[C, R, U, rpocts.Entity]) hndcts.GenericHandler[C, R, U] {
	return &ginHandler[C, R, U]{hp: hp, service: s}
}

func (h *ginHandler[C, R, U]) GetHandlerPath() string {
	return GetHandlerPath(h.hp)
}

func (h *ginHandler[C, R, U]) RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	handlerPath := GetHandlerPath(h.hp)
	rg.Handle(strings.ToUpper(method), handlerPath, handler)
}

func (h *ginHandler[C, R, U]) RegisterInstance(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	fullPath := GetHandlerPath(h.hp) + "/:id"
	rg.Handle(strings.ToUpper(method), fullPath, handler)
}

func (h *ginHandler[C, R, U]) Create(c *gin.Context) {
	var dto C
	if !BindJSONToDTO(c, &dto) {
		return
	}
	created, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ginHandler[C, R, U]) GetByID(c *gin.Context) {
	id, ok := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if !ok {
		return
	}
	dto, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

func (h *ginHandler[C, R, U]) List(c *gin.Context) {
	params := BuildListParamsFromRequest(c)
	list, _, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)

}

func (h *ginHandler[C, R, U]) Update(c *gin.Context) {
	id, ok := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if !ok {
		return
	}
	var dto U
	if !BindJSONToDTO(c, &dto) {
		return
	}
	updated, err := h.service.Update(c.Request.Context(), id, dto)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ginHandler[C, R, U]) Delete(c *gin.Context) {
	id, ok := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
