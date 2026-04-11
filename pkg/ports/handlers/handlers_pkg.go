package handlers

import (
	"github.com/brunojet/go-infra-backend/internal/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	"github.com/gin-gonic/gin"
)

type (
	GenericHandler[C, R, U any]       = hndcontracts.GenericHandler[C, R, U]
	NestedGenericHandler[C, R, U any] = hndcontracts.NestedGenericHandler[C, R, U]
	HandlerParameters                 = handlers.HandlerParameters
)

var (
	Int64GtZero   = handlers.Int64GtZero
	Int32GtZero   = handlers.Int32GtZero
	Int16GteZero  = handlers.Int16GteZero
	Base64UrlSafe = handlers.Base64UrlSafe
)

func NewGenericHandler[C, R, U any](hp HandlerParameters, s svccontracts.Service[C, R, U]) GenericHandler[C, R, U] {
	return handlers.NewGenericHandler(hp, s)
}

func NewGenericNestedHandler[C, R, U any](hp HandlerParameters, php HandlerParameters, s svccontracts.NestedService[C, R, U]) NestedGenericHandler[C, R, U] {
	return handlers.NewGenericNestedHandler(hp, php, s)
}

// Register is inherited from the aliased GinHandler type, but keeping gin import here
// helps discoverability for consumers.
var _ gin.HandlerFunc
