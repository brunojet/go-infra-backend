package contracts

import (
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/gin-gonic/gin"
)

// GenericHandler is a generic contract for handlers operating on a resource.
type GenericHandler[E rpocontracts.Entity, D any] interface {
	RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc)
	RegisterInstance(rg *gin.RouterGroup, method string, handler gin.HandlerFunc)
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// NestedGenericHandler is a contract for handlers that are nested under another resource.
// Only collection-level registration is supported for nested handlers.
type NestedGenericHandler[E rpocontracts.Entity, D any] interface {
	GenericHandler[E, D]
	CreateNested(c *gin.Context)
	ListNested(c *gin.Context)
}
