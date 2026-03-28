package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

// filterTypeService: service simples (não nested)
type filterTypeMapper struct{}

func (filterTypeMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterTypeID: id}, nil
}
func (filterTypeMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (filterTypeMapper) ToPostModel(dto dtos.FilterTypePost, model *models.FilterType) error {
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (filterTypeMapper) ToPatchModel(dto dtos.FilterTypePatch, model *models.FilterType) error {
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

// Model -> DTO
func (filterTypeMapper) ToDTO(model *models.FilterType, dto *dtos.FilterTypeGet) error {
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

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

// FilterService: nested em FilterType
type filterNestedMapper struct {
	ftm filterTypeMapper
}

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

func (filterNestedMapper) ToPostModel(dto dtos.FilterPost, model *models.Filter) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (filterNestedMapper) ToPatchModel(dto dtos.FilterPatch, model *models.Filter) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

// Model -> DTO
func (m filterNestedMapper) ToDTO(model *models.Filter, dto *dtos.FilterGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
	dto.FilterId = model.FilterId
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	// Mapear FilterType aninhado, se presente
	if model.FilterType != nil {
		dto.FilterType = &dtos.FilterTypeGet{}
		if err := m.ftm.ToDTO(model.FilterType, dto.FilterType); err != nil {
			return err
		}
	}
	return nil
}

func (filterNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterID: id}, nil
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
