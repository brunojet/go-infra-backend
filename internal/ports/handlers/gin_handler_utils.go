package handlers

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

var (
	paginationIgnoredKeys = map[string]struct{}{
		QueryParamPage:    {},
		QueryParamSize:    {},
		QueryParamOrderBy: {},
		QueryParamOrder:   {},
	}
)

func MapErrorToStatus(err error) int {
	switch err {
	case repositories.ErrConflictValidationFailed:
		return http.StatusConflict
	case ErrInvalidJSONBody, ErrInvalidIDFormat:
		return http.StatusBadRequest
	case repositories.ErrNotFound:
		return http.StatusNotFound
	case repositories.ErrDBUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func GetHandlerPath(handlerParameters HandlerParameters) string {
	handlerPath := strings.Trim(handlerParameters.HandlerPath, "/")
	if handlerPath == "" {
		log.Default().Panic("HandlerPath cannot be empty")
	}
	return handlerPath
}

func GetValidatedIDFromParam(c *gin.Context, paramName string, validationRule *regexp.Regexp) (string, bool) {
	id := c.Param(paramName)
	if validationRule != nil && !validationRule.MatchString(id) {
		SetResponseFromError(c, ErrInvalidIDFormat)
		return "", false
	}
	return id, true
}

func SetResponseFromError(c *gin.Context, err error) {
	// Record the error on the gin.Context so otelgin can span.RecordError(...) if enabled.
	_ = c.Error(err)
	status := MapErrorToStatus(err)
	c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}

func BindJSONToDTO(c *gin.Context, dto any) bool {
	if err := c.ShouldBindJSON(dto); err != nil {
		SetResponseFromError(c, ErrInvalidJSONBody)
		return false
	}
	return true
}

func ParsePaginationParams(values url.Values) (int, int) {
	page, _ := strconv.Atoi(values.Get(QueryParamPage))
	size, _ := strconv.Atoi(values.Get(QueryParamSize))
	if page < DefaultPage {
		page = DefaultPage
	}
	if size < DefaultPageSize {
		size = DefaultPageSize
	} else if size > MaxPageSize {
		size = MaxPageSize
	}
	return page, size
}

func ParseOrderParams(values url.Values) (string, string) {
	orderBy := values.Get(QueryParamOrderBy)
	order := strings.ToLower(values.Get(QueryParamOrder))
	if order != OrderAsc && order != OrderDesc {
		order = OrderAsc
	}
	return orderBy, order
}

func ExtractQueryScopes(values url.Values) map[string]any {
	scopes := make(map[string]any, len(values))
	for key, vals := range values {
		if _, ignored := paginationIgnoredKeys[key]; ignored {
			continue
		}
		if len(vals) == 0 {
			continue
		}
		scopes[key] = vals[0]
	}
	return scopes
}

func BuildListParamsFromRequest(c *gin.Context) contracts.ListParams {
	query := c.Request.URL.Query()
	page, size := ParsePaginationParams(query)
	orderBy, order := ParseOrderParams(query)
	return contracts.ListParams{
		QueryParams: contracts.QueryParams{Scopes: ExtractQueryScopes(query)},
		Page:        page,
		Size:        size,
		OrderBy:     orderBy,
		Order:       order,
	}
}
