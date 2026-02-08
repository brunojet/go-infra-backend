package handlers

import (
	internalhandlers "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

type GinHandler[E contracts.Entity, D any] = internalhandlers.GinHandler[E, D]

var (
	ErrInvalidJSONBody   = internalhandlers.ErrInvalidJSONBody
	MapErrorToStatus     = internalhandlers.MapErrorToStatus
	SetResponseFromError = internalhandlers.SetResponseFromError
	BindJSONToDTOPtr     = internalhandlers.BindJSONToDTOPtr
)

func NewGenericHandler[E contracts.Entity, D any](s svccontracts.Service[D, E]) *GinHandler[E, D] {
	return internalhandlers.NewGenericHandler[E, D](s)
}

// Register is inherited from the aliased GinHandler type, but keeping gin import here
// helps discoverability for consumers.
var _ gin.HandlerFunc
