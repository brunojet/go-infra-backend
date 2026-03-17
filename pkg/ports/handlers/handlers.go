package handlers

import (
	internalhandlers "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

type (
	GenericHandler[E rpocontracts.Entity, D any]       = hndcontracts.GenericHandler[E, D]
	NestedGenericHandler[E rpocontracts.Entity, D any] = hndcontracts.NestedGenericHandler[E, D]
	HandlerParameters                                  = internalhandlers.HandlerParameters
)

var (
	ErrInvalidJSONBody   = internalhandlers.ErrInvalidJSONBody
	MapErrorToStatus     = internalhandlers.MapErrorToStatus
	SetResponseFromError = internalhandlers.SetResponseFromError
	Int64GtZero          = internalhandlers.Int64GtZero
	Int32GtZero          = internalhandlers.Int32GtZero
	Int16GteZero         = internalhandlers.Int16GteZero
	Base64UrlSafe        = internalhandlers.Base64UrlSafe
)

func NewGenericHandler[E rpocontracts.Entity, D any](hp *HandlerParameters, s svccontracts.Service[D, E]) GenericHandler[E, D] {
	return internalhandlers.NewGenericHandler[E](hp, s)
}

func NewGenericNestedHandler[E rpocontracts.Entity, D any](hp *HandlerParameters, php *HandlerParameters, s svccontracts.NestedService[D, E]) NestedGenericHandler[E, D] {
	return internalhandlers.NewGenericNestedHandler[E](hp, php, s)
}

func BindJSONToDTOPtr[D any](c *gin.Context) (*D, error) {
	return internalhandlers.BindJSONToDTOPtr[D](c)
}

// Register is inherited from the aliased GinHandler type, but keeping gin import here
// helps discoverability for consumers.
var _ gin.HandlerFunc
