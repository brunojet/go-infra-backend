package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type TerminalModelService interface {
	svcContracts.Service[dtos.TerminalModelPost, dtos.TerminalModelGet, dtos.TerminalModelPatch, models.TerminalModel]
}

type terminalModelService struct {
	svcContracts.Service[dtos.TerminalModelPost, dtos.TerminalModelGet, dtos.TerminalModelPatch, models.TerminalModel]
}

func NewTerminalModelService(repo repoContracts.Repository[models.TerminalModel]) TerminalModelService {
	return &terminalModelService{
		Service: services.NewServiceImpl(repo, terminalModelMapper{}),
	}
}

type TerminalModelConfigurationNestedService interface {
	svcContracts.NestedService[dtos.TerminalModelConfigurationPost, dtos.TerminalModelConfigurationGet, dtos.TerminalModelConfigurationPatch, models.TerminalModelConfiguration]
}

type terminalModelConfigurationNestedService struct {
	svcContracts.NestedService[dtos.TerminalModelConfigurationPost, dtos.TerminalModelConfigurationGet, dtos.TerminalModelConfigurationPatch, models.TerminalModelConfiguration]
}

func NewTerminalModelConfigurationNestedService(repo repoContracts.Repository[models.TerminalModelConfiguration]) TerminalModelConfigurationNestedService {
	return &terminalModelConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, terminalModelConfigurationNestedMapper{}),
	}
}
