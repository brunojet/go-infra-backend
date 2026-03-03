package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	repoerrs "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSONBody = errors.New("invalid JSON body")
)

const (
	queryParamPage    = "page"
	queryParamSize    = "size"
	queryParamOrderBy = "orderBy"
	queryParamOrder   = "order"
)

func MapErrorToStatus(err error) int {
	if errors.Is(err, ErrInvalidJSONBody) {
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

func SetResponseFromError(c *gin.Context, err error) {
	// Record the error on the gin.Context so otelgin can span.RecordError(...) if enabled.
	_ = c.Error(err)
	status := MapErrorToStatus(err)
	c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}

// BindJSONToDTOPtr lê o corpo JSON do request e vincula em um DTO genérico `D`.
// Retorna ponteiro para `D` ou erro já tratado na resposta HTTP.
func BindJSONToDTOPtr[D any](c *gin.Context) (*D, error) {
	var dto D
	if err := c.ShouldBindJSON(&dto); err != nil {
		SetResponseFromError(c, ErrInvalidJSONBody)
		return nil, ErrInvalidJSONBody
	}
	return &dto, nil
}

func ParsePaginationParams(values url.Values) (int, int) {
	page, _ := strconv.Atoi(values.Get(queryParamPage))
	size, _ := strconv.Atoi(values.Get(queryParamSize))
	return page, size
}

func ExtractQueryScopes(values url.Values, ignoredKeys map[string]struct{}) map[string]any {
	scopes := make(map[string]any)
	for key, vals := range values {
		if _, ignored := ignoredKeys[key]; ignored {
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

	ignored := map[string]struct{}{
		queryParamPage:    {},
		queryParamSize:    {},
		queryParamOrderBy: {},
		queryParamOrder:   {},
	}

	orderBy := query.Get(queryParamOrderBy)

	return svccontracts.ListParams{
		QueryParams: svccontracts.QueryParams{Scopes: ExtractQueryScopes(query, ignored)},
		Page:        page,
		Size:        size,
		OrderBy:     orderBy,
		Order:       query.Get(queryParamOrder),
	}
}
