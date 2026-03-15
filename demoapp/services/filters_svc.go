package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// FilterTypeService: service simples (não nested)
type filterTypeMapper struct{}

func (filterTypeMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterTypeID: id}, nil
}
func (filterTypeMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (filterTypeMapper) ToModel(dto *models.FilterType, model *models.FilterType) {
	*model = *dto
}

func (filterTypeMapper) ToDTO(model *models.FilterType, dto *models.FilterType) {
	*dto = *model
}

type FilterTypeService struct {
	svcContracts.Service[models.FilterType, models.FilterType]
}

func NewFilterTypeService(repo repoContracts.Repository[models.FilterType]) *FilterTypeService {
	return &FilterTypeService{
		Service: services.NewServiceImpl(repo, filterTypeMapper{}),
	}
}

// FilterService: nested em FilterType
type filterNestedMapper struct{}

func (filterNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m filterNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}
	parentScopeId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return nil, err
	}
	mappedQueryScopes[models.ColFilterTypeID] = parentScopeId
	return mappedQueryScopes, nil
}

func (m filterNestedMapper) ApplyParentScopes(parentID string, model *models.Filter) error {
	parentScopeId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.FilterTypeId = parentScopeId
	return nil
}

func (filterNestedMapper) ToModel(dto *models.Filter, model *models.Filter) {
	*model = *dto
}

func (filterNestedMapper) ToDTO(model *models.Filter, dto *models.Filter) {
	*dto = *model
}

func (filterNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterID: id}, nil
}

type FilterNestedService struct {
	svcContracts.NestedService[models.Filter, models.Filter]
}

func NewFilterNestedService(repo repoContracts.Repository[models.Filter]) *FilterNestedService {
	return &FilterNestedService{
		NestedService: services.NewNestedServiceImpl(repo, filterNestedMapper{}),
	}
}
