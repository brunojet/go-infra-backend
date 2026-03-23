package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type terminalModelMapper struct{}

func (terminalModelMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColTerminalModelID: id}, nil
}

func (terminalModelMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (terminalModelMapper) ToModel(dto *models.TerminalModel, model *models.TerminalModel) {
	*model = *dto
}

func (terminalModelMapper) ToDTO(model *models.TerminalModel, dto *models.TerminalModel) {
	*dto = *model
}

type TerminalModelService interface {
	svcContracts.Service[models.TerminalModel, models.TerminalModel]
}

type terminalModelService struct {
	svcContracts.Service[models.TerminalModel, models.TerminalModel]
}

func NewTerminalModelService(repo repoContracts.Repository[models.TerminalModel]) TerminalModelService {
	return &terminalModelService{
		Service: services.NewServiceImpl(repo, terminalModelMapper{}),
	}
}

// terminalModelConfigurationNestedService: nested em TerminalModel
type terminalModelConfigurationNestedMapper struct{}

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
func (terminalModelConfigurationNestedMapper) ToModel(dto *models.TerminalModelConfiguration, model *models.TerminalModelConfiguration) {
	*model = *dto
}
func (terminalModelConfigurationNestedMapper) ToDTO(model *models.TerminalModelConfiguration, dto *models.TerminalModelConfiguration) {
	*dto = *model
}
func (terminalModelConfigurationNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColTerminalModelConfigurationID: id}, nil
}

type TerminalModelConfigurationNestedService interface {
	svcContracts.NestedService[models.TerminalModelConfiguration, models.TerminalModelConfiguration]
}

type terminalModelConfigurationNestedService struct {
	svcContracts.NestedService[models.TerminalModelConfiguration, models.TerminalModelConfiguration]
}

func NewTerminalModelConfigurationNestedService(repo repoContracts.Repository[models.TerminalModelConfiguration]) TerminalModelConfigurationNestedService {
	return &terminalModelConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, terminalModelConfigurationNestedMapper{}),
	}
}
