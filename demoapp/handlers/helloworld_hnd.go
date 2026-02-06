package handlers

import (
	helloworldRepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	helloworldSvc "github.com/brunojet/go-infra-backend/demoapp/services"
	hnd "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	svcContracts "github.com/brunojet/go-infra-backend/internal/ports/services/contracts"
)

type HelloWorldHandler struct {
	hnd.GinHandler[helloworldRepo.HelloWorld, helloworldSvc.HelloWorldDTO]
}

func NewHelloWorldHandler(s svcContracts.Service[helloworldSvc.HelloWorldDTO, helloworldRepo.HelloWorld]) *HelloWorldHandler {
	return &HelloWorldHandler{
		GinHandler: *hnd.NewGenericHandler[helloworldRepo.HelloWorld, helloworldSvc.HelloWorldDTO](s),
	}
}
