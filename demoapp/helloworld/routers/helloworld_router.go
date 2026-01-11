package routers

import (
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/handlers"
	"github.com/brunojet/go-infra-backend/internal/http/contracts"
)

type HelloWorldRouter struct {
	handler *handlers.HelloWorldHandler
}

func NewHelloWorldRouter(handler *handlers.HelloWorldHandler) *HelloWorldRouter {
	return &HelloWorldRouter{handler: handler}
}

func (r *HelloWorldRouter) Register(router contracts.Router) {
	router.GET("/helloworld", r.handler.List)
	router.POST("/helloworld", r.handler.Create)
	router.GET("/helloworld/:id", r.handler.Get)
	router.PATCH("/helloworld/:id", r.handler.Patch)
	router.DELETE("/helloworld/:id", r.handler.Delete)
}
