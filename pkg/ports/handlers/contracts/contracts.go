package contracts

import (
	"github.com/gin-gonic/gin"
)

// GenericHandler is a generic contract for handlers operating on a resource.
type GenericHandler[C, R, U any] interface {
	GetHandlerPath() string
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
type NestedGenericHandler[C, R, U any] interface {
	GenericHandler[C, R, U]
	CreateNested(c *gin.Context)
	ListNested(c *gin.Context)
}
