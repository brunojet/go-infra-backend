package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	rpoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ApplicationService interface {
	contracts.Service[dtos.ApplicationPost, dtos.ApplicationGet, dtos.ApplicationPatch, models.Application]
}

type applicationService struct {
	contracts.Service[dtos.ApplicationPost, dtos.ApplicationGet, dtos.ApplicationPatch, models.Application]
}

func NewApplicationService(repo rpoContracts.Repository[models.Application]) ApplicationService {
	return &applicationService{
		Service: services.NewServiceImpl(repo, applicationMapper{}),
	}
}

type ApplicationConfigurationNestedService interface {
	contracts.NestedService[dtos.ApplicationConfigurationPost, dtos.ApplicationConfigurationGet, dtos.ApplicationConfigurationPatch, models.ApplicationConfiguration]
}

type applicationConfigurationNestedService struct {
	contracts.NestedService[dtos.ApplicationConfigurationPost, dtos.ApplicationConfigurationGet, dtos.ApplicationConfigurationPatch, models.ApplicationConfiguration]
}

func NewApplicationConfigurationNestedService(repo rpoContracts.Repository[models.ApplicationConfiguration]) ApplicationConfigurationNestedService {
	return &applicationConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, applicationConfigurationNestedMapper{}),
	}
}
