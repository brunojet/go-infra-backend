package services

import (
	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type applicationMapper struct{}

func (applicationMapper) ToPostModel(dto dtos.ApplicationPost, model *models.Application) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.CustomerId = utils.ToNullString(dto.CustomerId)
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (applicationMapper) ToPatchModel(dto dtos.ApplicationPatch, model *models.Application) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.CustomerId = utils.ToNullString(dto.CustomerId)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (applicationMapper) ToDTO(model *models.Application, dto *dtos.ApplicationGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	dto.ApplicationId = model.ApplicationId
	dto.CustomerId = utils.FromNullString(model.CustomerId)
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (applicationMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationID: id}, nil
}

func (applicationMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

type applicationConfigurationNestedMapper struct{}

func (applicationConfigurationNestedMapper) ToPostModel(dto dtos.ApplicationConfigurationPost, model *models.ApplicationConfiguration) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.TerminalModelConfigurationId = dto.TerminalModelConfigurationId
	model.PackageName = utils.ToNullString(dto.PackageName)
	return nil
}

func (applicationConfigurationNestedMapper) ToPatchModel(dto dtos.ApplicationConfigurationPatch, model *models.ApplicationConfiguration) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.PackageName = utils.ToNullString(dto.PackageName)
	return nil
}

func (applicationConfigurationNestedMapper) ToDTO(model *models.ApplicationConfiguration, dto *dtos.ApplicationConfigurationGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	applicationConfigurationId, err := utils.EncodeCompositeKey(model.ApplicationId, model.TerminalModelConfigurationId)
	if err != nil {
		return err
	}
	dto.ApplicationConfigurationId = applicationConfigurationId
	dto.TerminalModelConfigurationId = model.TerminalModelConfigurationId
	dto.PackageName = utils.FromNullString(model.PackageName)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (applicationConfigurationNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationConfigurationID: id}, nil
}

func (applicationConfigurationNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationConfigurationNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}
	parentScopeId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return nil, err
	}
	mappedQueryScopes[models.ColApplicationID] = parentScopeId
	return mappedQueryScopes, nil
}

func (m applicationConfigurationNestedMapper) ApplyParentScopes(parentID string, model *models.ApplicationConfiguration) error {
	parentScopeId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.ApplicationId = parentScopeId
	return nil
}
