package services

import (
	"database/sql"
	"time"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// filterTypeService: service simples (não nested)
type filterTypeMapper struct{}

func (filterTypeMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterTypeID: id}, nil
}
func (filterTypeMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

// DTO -> Model
func (filterTypeMapper) ToModel(dto *dtos.FilterTypeDTO, model *models.FilterType) {
	model.Name = sql.NullString{String: dto.Name, Valid: dto.Name != ""}
	model.Description = sql.NullString{String: dto.Description, Valid: dto.Description != ""}
}

// Model -> DTO
func (filterTypeMapper) ToDTO(model *models.FilterType, dto *dtos.FilterTypeDTO) {
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = model.Name.String
	if model.Description.Valid {
		dto.Description = model.Description.String
	} else {
		dto.Description = ""
	}
	if model.CreatedAt.Valid {
		dto.CreatedAt = model.CreatedAt.Time.Format(time.RFC3339)
	} else {
		dto.CreatedAt = ""
	}
	if model.UpdatedAt.Valid {
		dto.UpdatedAt = model.UpdatedAt.Time.Format(time.RFC3339)
	} else {
		dto.UpdatedAt = ""
	}
	if model.DeletedAt.Valid {
		dto.DeletedAt = model.DeletedAt.Time.Format(time.RFC3339)
	} else {
		dto.DeletedAt = ""
	}
}

type FilterTypeService interface {
	svcContracts.Service[dtos.FilterTypeDTO, models.FilterType]
}

type filterTypeService struct {
	svcContracts.Service[dtos.FilterTypeDTO, models.FilterType]
}

func NewFilterTypeService(repo repoContracts.Repository[models.FilterType]) FilterTypeService {
	return &filterTypeService{
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

// DTO -> Model
func (filterNestedMapper) ToModel(dto *dtos.FilterDTO, model *models.Filter) {
	model.FilterId = dto.FilterId
	model.FilterTypeId = dto.FilterTypeId
	model.Name = sql.NullString{String: dto.Name, Valid: dto.Name != ""}
	model.Description = sql.NullString{String: dto.Description, Valid: dto.Description != ""}
	// Datas não são preenchidas no ToModel (normalmente gerenciadas pelo banco)
}

// Model -> DTO
func (filterNestedMapper) ToDTO(model *models.Filter, dto *dtos.FilterDTO) {
	dto.FilterId = model.FilterId
	dto.FilterTypeId = model.FilterTypeId
	dto.Name = model.Name.String
	if model.Description.Valid {
		dto.Description = model.Description.String
	} else {
		dto.Description = ""
	}
	if model.CreatedAt.Valid {
		dto.CreatedAt = model.CreatedAt.Time.Format(time.RFC3339)
	} else {
		dto.CreatedAt = ""
	}
	if model.UpdatedAt.Valid {
		dto.UpdatedAt = model.UpdatedAt.Time.Format(time.RFC3339)
	} else {
		dto.UpdatedAt = ""
	}
	if model.DeletedAt.Valid {
		dto.DeletedAt = model.DeletedAt.Time.Format(time.RFC3339)
	} else {
		dto.DeletedAt = ""
	}
	// Mapear FilterType aninhado, se presente
	if model.FilterType != nil {
		var filterTypeDTO dtos.FilterTypeDTO
		filterTypeMapper{}.ToDTO(model.FilterType, &filterTypeDTO)
		dto.FilterType = filterTypeDTO
	} else {
		dto.FilterType = dtos.FilterTypeDTO{}
	}
}

func (filterNestedMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{models.ColFilterID: id}, nil
}

type FilterNestedService interface {
	svcContracts.NestedService[dtos.FilterDTO, models.Filter]
}

type filterNestedService struct {
	svcContracts.NestedService[dtos.FilterDTO, models.Filter]
}

func NewFilterNestedService(repo repoContracts.Repository[models.Filter]) FilterNestedService {
	return &filterNestedService{
		NestedService: services.NewNestedServiceImpl(repo, filterNestedMapper{}),
	}
}
