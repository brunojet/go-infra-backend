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

type ginHandler[E rpocontracts.Entity, D any] struct {
	hp      *HandlerParameters
	service svccontracts.Service[D, E]
}

func NewGenericHandler[E rpocontracts.Entity, D any](hp *HandlerParameters, s svccontracts.Service[D, E]) hndcontracts.GenericHandler[E, D] {
	return &ginHandler[E, D]{hp: hp, service: s}
}

func (h *ginHandler[E, D]) Register(rg *gin.RouterGroup, method string, handler gin.HandlerFunc) {
	handlerPath := strings.Trim(h.hp.HandlerPath, "/")
	if handlerPath == "" {
		log.Default().Panic("HandlerPath cannot be empty")
	}
	rg.Handle(strings.ToUpper(method), handlerPath, handler)
}

func (h *ginHandler[E, D]) Create(c *gin.Context) {
	dto, err := BindJSONToDTOPtr[D](c)
	if err != nil {
		return
	}
	if err := h.service.Create(c.Request.Context(), dto); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto)
}

func (h *ginHandler[E, D]) GetByID(c *gin.Context) {
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

func (h *ginHandler[E, D]) List(c *gin.Context) {
	params := BuildListParamsFromRequest(c)
	list, _, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)

}

func (h *ginHandler[E, D]) Update(c *gin.Context) {
	id, err := GetValidatedIDFromParam(c, "id", h.hp.IDValidationRule)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	dto, err := BindJSONToDTOPtr[D](c)
	if err != nil {
		return
	}
	if err := h.service.Update(c.Request.Context(), id, dto); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

func (h *ginHandler[E, D]) Delete(c *gin.Context) {
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
