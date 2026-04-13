package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
)

type nestedGinHandler[C, R, U any] struct {
	hndcontracts.GenericHandler[C, R, U]
	php HandlerParameters
	svc svccontracts.NestedService[C, R, U]
}

func NewGenericNestedHandler[C, R, U any](php HandlerParameters, hp HandlerParameters, s svccontracts.NestedService[C, R, U]) hndcontracts.NestedGenericHandler[C, R, U] {
	baseHandler := NewGenericHandler(hp, s)
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
	var created R
	err := h.svc.CreateNested(c.Request.Context(), parentID, dto, &created)
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
	responses := make([]R, 0, params.Size)
	totalItems, err := h.svc.ListNested(c.Request.Context(), parentID, params, &responses)
	if err != nil {
		SetResponseFromError(c, err)
		return
	}
	listResponse := ListResponses[R]{
		Data: responses,
		Pagination: PaginationResponse{
			Page:       params.Page,
			Size:       len(responses),
			TotalItems: totalItems,
			OrderBy:    params.OrderBy,
			Order:      params.Order,
		},
	}
	c.JSON(http.StatusPartialContent, listResponse)
}
