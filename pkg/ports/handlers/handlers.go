package handlers

import (
	internalhandlers "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	repo "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	svc "github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/gin-gonic/gin"
)

type GinHandler[E repo.Entity, D any] = internalhandlers.GinHandler[E, D]

var (
	ErrInvalidJSONBody   = internalhandlers.ErrInvalidJSONBody
	MapErrorToStatus     = internalhandlers.MapErrorToStatus
	SetResponseFromError = internalhandlers.SetResponseFromError
)

func NewGenericHandler[E repo.Entity, D any](s svc.Service[D, E]) *GinHandler[E, D] {
	return internalhandlers.NewGenericHandler[E, D](s)
}

func BindJSONToDTOPtr[D any](c *gin.Context) (*D, error) {
	return internalhandlers.BindJSONToDTOPtr[D](c)
}

// Register is inherited from the aliased GinHandler type, but keeping gin import here
// helps discoverability for consumers.
var _ gin.HandlerFunc
