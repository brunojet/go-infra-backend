package handlers

import (
	"errors"
	"net/http"

	repoerrs "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSONBody = errors.New("invalid JSON body")
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
	c.JSON(MapErrorToStatus(err), gin.H{"error": err.Error()})
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
