package services

import (
	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type filterTypeMapper struct{}

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

func (filterTypeMapper) ToDTO(model *models.FilterType, dto *dtos.FilterTypeGet) error {
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (filterTypeMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterTypeID: id}, nil
}

func (filterTypeMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

type filterNestedMapper struct {
	ftm filterTypeMapper
}

func (filterNestedMapper) ToPostModel(dto dtos.FilterPost, model *models.Filter) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (filterNestedMapper) toAssociativeModel(filterIds []int64, models *[]models.Filter) error {
	debugassert.Assert(models != nil, "models cannot be nil")
	for id, filterId := range filterIds {
		(*models)[id].FilterId = filterId
	}
	return nil
}

func (filterNestedMapper) ToPatchModel(dto dtos.FilterPatch, model *models.Filter) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (m filterNestedMapper) ToDTO(model *models.Filter, dto *dtos.FilterGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	dto.FilterId = model.FilterId
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
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
	debugassert.Assert(model != nil, "model cannot be nil")
	parentScopeId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.FilterTypeId = parentScopeId
	return nil
}
