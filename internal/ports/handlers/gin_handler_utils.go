package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	repoerrs "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
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
	if errors.Is(err, ErrInvalidJSONBody) ||
		errors.Is(err, ErrInvalidIDFormat) {
		return http.StatusBadRequest
	}
	if errors.Is(err, repoerrs.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, repoerrs.ErrDBUnavailable) {
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}

func GetHandlerPath(handlerParameters *HandlerParameters) string {
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

// BindJSONToDTOPtr lê o corpo JSON do request e vincula em um DTO genérico `D`.
// Retorna ponteiro para `D` ou erro já tratado na resposta HTTP.
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
	return page, size
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

func BuildListParamsFromRequest(c *gin.Context) svccontracts.ListParams {
	query := c.Request.URL.Query()
	page, size := ParsePaginationParams(query)
	orderBy := query.Get(QueryParamOrderBy)
	return svccontracts.ListParams{
		QueryParams: svccontracts.QueryParams{Scopes: ExtractQueryScopes(query)},
		Page:        page,
		Size:        size,
		OrderBy:     orderBy,
		Order:       query.Get(QueryParamOrder),
	}
}
