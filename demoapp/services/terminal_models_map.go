package services

import (
	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type terminalModelMapper struct{}

func (terminalModelMapper) ToPostModel(dto dtos.TerminalModelPost, model *models.TerminalModel) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (terminalModelMapper) ToPatchModel(dto dtos.TerminalModelPatch, model *models.TerminalModel) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (terminalModelMapper) ToDTO(model *models.TerminalModel, dto *dtos.TerminalModelGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	dto.TerminalModelId = model.TerminalModelId
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (terminalModelMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColTerminalModelID: id}, nil
}

func (terminalModelMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

type terminalModelConfigurationNestedMapper struct {
	tmm terminalModelMapper
}

func (terminalModelConfigurationNestedMapper) ToPostModel(dto dtos.TerminalModelConfigurationPost, model *models.TerminalModelConfiguration) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.IntegrationType = utils.ToNullInt16(integrationTypeMapToModel[dto.IntegrationType])
	return nil
}

func (terminalModelConfigurationNestedMapper) ToPatchModel(dto dtos.TerminalModelConfigurationPatch, model *models.TerminalModelConfiguration) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.IntegrationType = utils.ToNullInt16(integrationTypeMapToModel[dto.IntegrationType])
	return nil
}

func (m terminalModelConfigurationNestedMapper) ToDTO(model *models.TerminalModelConfiguration, dto *dtos.TerminalModelConfigurationGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	dto.TerminalModelConfigurationId = model.TerminalModelConfigurationId
	dto.TerminalModelId = model.TerminalModelId
	dto.IntegrationType = integrationTypeMapFromModel[utils.FromNullInt16(model.IntegrationType)]
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	if model.TerminalModel != nil {
		dto.TerminalModel = &dtos.TerminalModelGet{}
		if err := m.tmm.ToDTO(model.TerminalModel, dto.TerminalModel); err != nil {
			return err
		}
	}
	return nil
}

func (terminalModelConfigurationNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColTerminalModelConfigurationID: id}, nil
}

func (terminalModelConfigurationNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m terminalModelConfigurationNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}
	parentId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return nil, err
	}
	mappedQueryScopes[models.ColTerminalModelID] = parentId
	return mappedQueryScopes, nil
}

func (m terminalModelConfigurationNestedMapper) ApplyParentScopes(parentID string, model *models.TerminalModelConfiguration) error {
	parentId, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.TerminalModelId = parentId
	return nil
}
