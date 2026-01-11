package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/core/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/helloworld/services"
	infra "github.com/brunojet/go-infra-backend/internal/http/ports"
)

type HelloWorldHandler struct {
	*infra.CRUDHandler[
		dtos.HelloWorldResponse,
		dtos.CreateHelloWorldRequest,
		dtos.PatchHelloWorldRequest,
		dtos.HelloWorldResponse,
	]
}

func NewHelloWorldHandler(service *services.HelloWorldService) *HelloWorldHandler {
	return &HelloWorldHandler{
		CRUDHandler: &infra.CRUDHandler[
			dtos.HelloWorldResponse,
			dtos.CreateHelloWorldRequest,
			dtos.PatchHelloWorldRequest,
			dtos.HelloWorldResponse,
		]{
			Service:    service,
			ToResponse: func(r dtos.HelloWorldResponse) dtos.HelloWorldResponse { return r },
		},
	}
}
