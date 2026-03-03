package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	repocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svc "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type GinHandler[E repocontracts.Entity, D any] struct {
	service svc.Service[D, E]
}

func NewGenericHandler[E repocontracts.Entity, D any](s svc.Service[D, E]) *GinHandler[E, D] {
	return &GinHandler[E, D]{service: s}
}

func (h *GinHandler[E, D]) Register(rg *gin.RouterGroup, method, path string, handler gin.HandlerFunc) {
	rg.Handle(strings.ToUpper(method), path, handler)
}

func (h *GinHandler[E, D]) Create(c *gin.Context) {
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

func (h *GinHandler[E, D]) GetByID(c *gin.Context) {
	id := c.Param("id")
	dto, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

func (h *GinHandler[E, D]) List(c *gin.Context) {
	params := BuildListParamsFromRequest(c)
	list, _, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)

}

func (h *GinHandler[E, D]) Update(c *gin.Context) {
	id := c.Param("id")
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

func (h *GinHandler[E, D]) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		SetResponseFromError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
