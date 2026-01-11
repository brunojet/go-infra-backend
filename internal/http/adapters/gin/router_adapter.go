package ginadapter

import (
	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	group *gin.RouterGroup
}

func NewGinRouter(group *gin.RouterGroup) contracts.Router {
	return &GinRouter{group: group}
}

func (r *GinRouter) GET(path string, handler contracts.HandlerFunc) {
	r.group.GET(path, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (r *GinRouter) POST(path string, handler contracts.HandlerFunc) {
	r.group.POST(path, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (r *GinRouter) PATCH(path string, handler contracts.HandlerFunc) {
	r.group.PATCH(path, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (r *GinRouter) DELETE(path string, handler contracts.HandlerFunc) {
	r.group.DELETE(path, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (r *GinRouter) Group(relativePath string) contracts.Router {
	return &GinRouter{group: r.group.Group(relativePath)}
}
