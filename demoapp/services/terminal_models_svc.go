package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
)

type TerminalModelService interface {
	svccts.Service[dtos.TerminalModelPost, dtos.TerminalModelGet, dtos.TerminalModelPatch]
}

type terminalModelService struct {
	svccts.Service[dtos.TerminalModelPost, dtos.TerminalModelGet, dtos.TerminalModelPatch]
}

func NewTerminalModelService(repo rpocts.Repository[models.TerminalModel]) TerminalModelService {
	return &terminalModelService{
		Service: services.NewServiceImpl(repo, terminalModelMapper{}),
	}
}

type TerminalModelConfigurationNestedService interface {
	svccts.NestedService[dtos.TerminalModelConfigurationPost, dtos.TerminalModelConfigurationGet, dtos.TerminalModelConfigurationPatch]
}

type terminalModelConfigurationNestedService struct {
	svccts.NestedService[dtos.TerminalModelConfigurationPost, dtos.TerminalModelConfigurationGet, dtos.TerminalModelConfigurationPatch]
}

func NewTerminalModelConfigurationNestedService(repo rpocts.Repository[models.TerminalModelConfiguration]) TerminalModelConfigurationNestedService {
	return &terminalModelConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, terminalModelConfigurationNestedMapper{}),
	}
}
