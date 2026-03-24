package handlers

import (
	internalhandlers "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccontracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/gin-gonic/gin"
)

type (
	GenericHandler[C, R, U any, E rpocontracts.Entity]       = hndcontracts.GenericHandler[C, R, U, E]
	NestedGenericHandler[C, R, U any, E rpocontracts.Entity] = hndcontracts.NestedGenericHandler[C, R, U, E]
	HandlerParameters                                        = internalhandlers.HandlerParameters
)

var (
	Int64GtZero   = internalhandlers.Int64GtZero
	Int32GtZero   = internalhandlers.Int32GtZero
	Int16GteZero  = internalhandlers.Int16GteZero
	Base64UrlSafe = internalhandlers.Base64UrlSafe
)

func NewGenericHandler[C, R, U any, E rpocontracts.Entity](hp HandlerParameters, s svccontracts.Service[C, R, U, E]) GenericHandler[C, R, U, E] {
	return internalhandlers.NewGenericHandler[C, R, U, E](hp, s)
}

func NewGenericNestedHandler[C, R, U any, E rpocontracts.Entity](hp HandlerParameters, php HandlerParameters, s svccontracts.NestedService[C, R, U, E]) NestedGenericHandler[C, R, U, E] {
	return internalhandlers.NewGenericNestedHandler[C, R, U, E](hp, php, s)
}

// Register is inherited from the aliased GinHandler type, but keeping gin import here
// helps discoverability for consumers.
var _ gin.HandlerFunc
