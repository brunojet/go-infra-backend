package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type FilterTypeService interface {
	svcContracts.Service[dtos.FilterTypePost, dtos.FilterTypeGet, dtos.FilterTypePatch, models.FilterType]
}

type filterTypeService struct {
	svcContracts.Service[dtos.FilterTypePost, dtos.FilterTypeGet, dtos.FilterTypePatch, models.FilterType]
}

func NewFilterTypeService(repo repoContracts.Repository[models.FilterType]) FilterTypeService {
	return &filterTypeService{
		Service: services.NewServiceImpl(repo, filterTypeMapper{}),
	}
}

type FilterNestedService interface {
	svcContracts.NestedService[dtos.FilterPost, dtos.FilterGet, dtos.FilterPatch, models.Filter]
}

type filterNestedService struct {
	svcContracts.NestedService[dtos.FilterPost, dtos.FilterGet, dtos.FilterPatch, models.Filter]
}

func NewFilterNestedService(repo repoContracts.Repository[models.Filter]) FilterNestedService {
	return &filterNestedService{
		NestedService: services.NewNestedServiceImpl(repo, filterNestedMapper{}),
	}
}
