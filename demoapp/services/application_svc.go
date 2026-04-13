package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
)

type ApplicationService interface {
	contracts.Service[dtos.ApplicationPost, dtos.ApplicationGet, dtos.ApplicationPatch]
}

type applicationService struct {
	contracts.Service[dtos.ApplicationPost, dtos.ApplicationGet, dtos.ApplicationPatch]
}

func NewApplicationService(repo rpocts.Repository[models.Application]) ApplicationService {
	return &applicationService{
		Service: services.NewServiceImpl(repo, applicationMapper{}),
	}
}

type ApplicationConfigurationNestedService interface {
	contracts.NestedService[dtos.ApplicationConfigurationPost, dtos.ApplicationConfigurationGet, dtos.ApplicationConfigurationPatch]
}

type applicationConfigurationNestedService struct {
	contracts.NestedService[dtos.ApplicationConfigurationPost, dtos.ApplicationConfigurationGet, dtos.ApplicationConfigurationPatch]
}

func NewApplicationConfigurationNestedService(repo rpocts.Repository[models.ApplicationConfiguration]) ApplicationConfigurationNestedService {
	return &applicationConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, applicationConfigurationNestedMapper{}),
	}
}
