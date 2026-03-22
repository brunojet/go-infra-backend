package contracts

import (
	rpocontracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/gin-gonic/gin"
)

// BaseHandlerMethods defines the contract for handler registration.
type BaseHandlerMethods[D any] interface {
	// RegisterCollection registers collection-level routes (e.g., /items).
	RegisterCollection(rg *gin.RouterGroup, method string, handler gin.HandlerFunc)
	// RegisterInstance registers instance-level routes (e.g., /items/:id).
	RegisterInstance(rg *gin.RouterGroup, method string, handler gin.HandlerFunc)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// GenericHandler is a generic contract for handlers operating on a resource.
type GenericHandler[E rpocontracts.Entity, D any] interface {
	BaseHandlerMethods[D]
	Create(c *gin.Context)
	List(c *gin.Context)
}

// NestedGenericHandler is a contract for handlers that are nested under another resource.
// Only collection-level registration is supported for nested handlers.
type NestedGenericHandler[E rpocontracts.Entity, D any] interface {
	// RegisterCollection registers collection-level routes for nested resources (e.g., /parent/:parentId/items).
	BaseHandlerMethods[D]
	CreateNested(c *gin.Context)
	ListNested(c *gin.Context)
}
