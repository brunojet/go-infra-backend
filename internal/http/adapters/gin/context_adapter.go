package ginadapter

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	c *gin.Context
}

func NewGinContext(c *gin.Context) contracts.Context {
	return &GinContext{c: c}
}

func (gc *GinContext) RequestContext() context.Context {
	return gc.c.Request.Context()
}

func (gc *GinContext) Param(name string) string {
	return gc.c.Param(name)
}

func (gc *GinContext) BindJSON(dst any) error {
	return gc.c.ShouldBindJSON(dst)
}

func (gc *GinContext) JSON(status int, v any) {
	gc.c.JSON(status, v)
}

func (gc *GinContext) String(status int, v string) {
	gc.c.String(status, v)
}

func (gc *GinContext) Status(status int) {
	gc.c.Status(status)
}
