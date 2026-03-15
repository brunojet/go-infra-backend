package handlers

import (
	helloworldRepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	helloworldSvc "github.com/brunojet/go-infra-backend/demoapp/services"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type HelloWorldHandler struct {
	hndcontracts.GenericHandler[helloworldRepo.HelloWorld, helloworldSvc.HelloWorldDTO]
}

func NewHelloWorldHandler(s svcContracts.Service[helloworldSvc.HelloWorldDTO, helloworldRepo.HelloWorld]) *HelloWorldHandler {
	return &HelloWorldHandler{
		GenericHandler: hnd.NewGenericHandler[helloworldRepo.HelloWorld](IDInt64Parameters, s),
	}
}
