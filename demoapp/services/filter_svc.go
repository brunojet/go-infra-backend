package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
)

type FilterTypeService interface {
	svccts.Service[dtos.FilterTypePost, dtos.FilterTypeGet, dtos.FilterTypePatch]
}

type filterTypeService struct {
	svccts.Service[dtos.FilterTypePost, dtos.FilterTypeGet, dtos.FilterTypePatch]
}

func NewFilterTypeService(repo rpocts.Repository[models.FilterType]) FilterTypeService {
	return &filterTypeService{
		Service: services.NewServiceImpl(repo, filterTypeMapper{}),
	}
}

type FilterNestedService interface {
	svccts.NestedService[dtos.FilterPost, dtos.FilterGet, dtos.FilterPatch]
}

type filterNestedService struct {
	svccts.NestedService[dtos.FilterPost, dtos.FilterGet, dtos.FilterPatch]
}

func NewFilterNestedService(repo rpocts.Repository[models.Filter]) FilterNestedService {
	return &filterNestedService{
		NestedService: services.NewNestedServiceImpl(repo, filterNestedMapper{}),
	}
}
