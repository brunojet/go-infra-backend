package services

import (
	"log"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	rpoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	utils "github.com/brunojet/go-infra-backend/pkg/utils"
)

type applicationMapper struct{}

// Converte de DTO para Model
func (applicationMapper) ToModel(dto *dtos.ApplicationDTO, model *models.Application) {
	model.CustomerId = utils.ToNullString(dto.CustomerId)
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
}

// Converte de Model para DTO
func (applicationMapper) ToDTO(model *models.Application, dto *dtos.ApplicationDTO) {
	dto.ApplicationId = model.ApplicationId
	dto.CustomerId = utils.FromNullString(model.CustomerId)
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	// TODO: Mapear relacionamentos aninhados se necessário
}

func (applicationMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationID: id}, nil
}

func (applicationMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

type ApplicationService interface {
	svcContracts.Service[dtos.ApplicationDTO, models.Application]
}

type applicationService struct {
	svcContracts.Service[dtos.ApplicationDTO, models.Application]
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

// Converte de DTO para Model
func (applicationConfigurationNestedMapper) ToModel(dto *dtos.ApplicationConfigurationDTO, model *models.ApplicationConfiguration) {
	model.TerminalModelConfigurationId = dto.TerminalModelConfigurationId
	model.PackageName = utils.ToNullString(dto.PackageName)
}

// Converte de Model para DTO
func (applicationConfigurationNestedMapper) ToDTO(model *models.ApplicationConfiguration, dto *dtos.ApplicationConfigurationDTO) {
	var err error
	dto.ApplicationConfigurationId, err = utils.EncodeCompositeKey(model.ApplicationId, model.TerminalModelConfigurationId)
	if err != nil {
		log.Panicf("failed to encode composite key for ApplicationConfiguration: %v", err)
	}
	dto.PackageName = utils.FromNullString(model.PackageName)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
}

func (applicationConfigurationNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColApplicationConfigurationID: id}, nil
}

type ApplicationConfigurationNestedService interface {
	svcContracts.NestedService[dtos.ApplicationConfigurationDTO, models.ApplicationConfiguration]
}

type applicationConfigurationNestedService struct {
	svcContracts.NestedService[dtos.ApplicationConfigurationDTO, models.ApplicationConfiguration]
}

func NewApplicationConfigurationNestedService(repo rpoContracts.Repository[models.ApplicationConfiguration]) ApplicationConfigurationNestedService {
	return &applicationConfigurationNestedService{
		NestedService: services.NewNestedServiceImpl(repo, applicationConfigurationNestedMapper{}),
	}
}
