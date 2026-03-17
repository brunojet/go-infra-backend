package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	rpoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type applicationMapper struct{}

func (applicationMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationID: id}, nil
}

func (applicationMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (applicationMapper) ToModel(dto *models.Application, model *models.Application) {
	*model = *dto
}

func (applicationMapper) ToDTO(model *models.Application, dto *models.Application) {
	*dto = *model
}

type ApplicationService interface {
	svcContracts.Service[models.Application, models.Application]
}

type applicationService struct {
	svcContracts.Service[models.Application, models.Application]
}

func NewApplicationService(repo rpoContracts.Repository[models.Application]) ApplicationService {
	return &applicationService{
		Service: services.NewServiceImpl(repo, applicationMapper{}),
	}
}

type applicationConfigurationNestedMapper struct{}

func (applicationConfigurationNestedMapper) DecodeParentID(parentID string) (map[string]any, error) {
	return map[string]any{models.ColApplicationID: parentID}, nil
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

func (applicationConfigurationNestedMapper) ToModel(dto *models.ApplicationConfiguration, model *models.ApplicationConfiguration) {
	*model = *dto
}

func (applicationConfigurationNestedMapper) ToDTO(model *models.ApplicationConfiguration, dto *models.ApplicationConfiguration) {
	*dto = *model
}

func (applicationConfigurationNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationConfigurationID: id}, nil
}

type ApplicationConfigurationNestedService interface {
	svcContracts.NestedService[models.ApplicationConfiguration, models.ApplicationConfiguration]
}

type applicationConfigurationNestedService struct {
	svcContracts.NestedService[models.ApplicationConfiguration, models.ApplicationConfiguration]
}

func NewApplicationConfigurationNestedService(repo rpoContracts.Repository[models.ApplicationConfiguration]) ApplicationConfigurationNestedService {
	return &applicationConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, applicationConfigurationNestedMapper{}),
	}
}
