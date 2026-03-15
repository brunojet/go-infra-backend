package contracts

import (
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/gin-gonic/gin"
)

type BaseHandlerMethods[D any] interface {
	Register(rg *gin.RouterGroup, method, path string, handler gin.HandlerFunc)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type GenericHandler[E rpocontracts.Entity, D any] interface {
	BaseHandlerMethods[D]
	Create(c *gin.Context)
	List(c *gin.Context)
}

type NestedGenericHandler[E rpocontracts.Entity, D any] interface {
	BaseHandlerMethods[D]
	RegisterNested(rg *gin.RouterGroup, method, parentPath, path string, handler gin.HandlerFunc)
	CreateNested(c *gin.Context)
	ListNested(c *gin.Context)
}
