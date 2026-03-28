package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	utils "github.com/brunojet/go-infra-backend/pkg/utils"
)

type terminalModelMapper struct{}

func (terminalModelMapper) ToPostModel(dto dtos.TerminalModelPost, model *models.TerminalModel) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (terminalModelMapper) ToPatchModel(dto dtos.TerminalModelPatch, model *models.TerminalModel) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	return nil
}

func (terminalModelMapper) ToDTO(model *models.TerminalModel, dto *dtos.TerminalModelGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
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

// terminalModelConfigurationNestedService: nested em TerminalModel
type terminalModelConfigurationNestedMapper struct {
	tmm terminalModelMapper
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

// Converte de DTO para Model
func (terminalModelConfigurationNestedMapper) ToPostModel(dto dtos.TerminalModelConfigurationPost, model *models.TerminalModelConfiguration) error {
	if model == nil {
		return errMapperNilModel
	}
	model.IntegrationType = utils.ToNullInt16(dto.IntegrationType)
	return nil
}

func (terminalModelConfigurationNestedMapper) ToPatchModel(dto dtos.TerminalModelConfigurationPatch, model *models.TerminalModelConfiguration) error {
	if model == nil {
		return errMapperNilModel
	}
	model.IntegrationType = utils.ToNullInt16(dto.IntegrationType)
	return nil
}

// Converte de Model para DTO
func (m terminalModelConfigurationNestedMapper) ToDTO(model *models.TerminalModelConfiguration, dto *dtos.TerminalModelConfigurationGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
	dto.TerminalModelConfigurationId = model.TerminalModelConfigurationId
	dto.TerminalModelId = model.TerminalModelId
	dto.IntegrationType = utils.FromNullInt16(model.IntegrationType)
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
