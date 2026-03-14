package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type applicationProfileNestedMapper struct{}

func (applicationProfileNestedMapper) DecodeParentID(parentID string) (map[string]any, error) {
	id, err := services.ParseScopeIntFromString[int64](parentID, 0)
	if err != nil {
		return nil, errNestedProfileApplicationIDRequired
	}
	return map[string]any{models.ColAppProfileApplicationID: id}, nil
}

func (applicationProfileNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationProfileNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}

	parentScopes, err := m.DecodeParentID(parentID)
	if err != nil {
		return nil, err
	}

	mergedScopes := make(map[string]any, len(parentScopes)+len(mappedQueryScopes))
	for key, value := range mappedQueryScopes {
		mergedScopes[key] = value
	}
	for key, value := range parentScopes {
		mergedScopes[key] = value
	}

	return mergedScopes, nil
}

func (m applicationProfileNestedMapper) ApplyParentScopes(parentID string, model *models.ApplicationProfile) error {
	parentScopes, err := m.DecodeParentID(parentID)
	if err != nil {
		return err
	}

	applicationID, err := services.ParseScopeInt[int64](parentScopes, models.ColAppProfileApplicationID, 1)
	if err != nil {
		return err
	}

	model.ApplicationId = applicationID
	return nil
}

func (applicationProfileNestedMapper) ToModel(dto *models.ApplicationProfile, model *models.ApplicationProfile) {
	*model = *dto
}

func (applicationProfileNestedMapper) ToDTO(model *models.ApplicationProfile, dto *models.ApplicationProfile) {
	*dto = *model
}

func (applicationProfileNestedMapper) GetModelKey(id string) (map[string]any, error) {
	profileID, err := services.ParseScopeIntFromString[int64](id, 0)
	if err != nil {
		return nil, errProfileScopeIDRequired
	}
	return map[string]any{models.ColAppProfileID: profileID}, nil
}

type ApplicationProfileNestedService interface {
	svcContracts.NestedService[models.ApplicationProfile, models.ApplicationProfile]
}

type applicationProfileNestedService struct {
	svcContracts.NestedService[models.ApplicationProfile, models.ApplicationProfile]
	repo   repo.ApplicationProfileRepository
	mapper applicationProfileNestedMapper
}

func NewApplicationProfileNestedService(r repo.ApplicationProfileRepository) ApplicationProfileNestedService {
	return &applicationProfileNestedService{
		NestedService: internalservices.NewNestedServiceImpl(r, applicationProfileNestedMapper{}),
		repo:          r,
		mapper:        applicationProfileNestedMapper{},
	}
}

func (s *applicationProfileNestedService) createOneShot(ctx context.Context, inOut *models.ApplicationProfile) error {
	return s.repo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, inOut); err != nil {
			return err
		}
		return s.repo.ArchiveStageDuplicates(txCtx, inOut)
	})
}

func (s *applicationProfileNestedService) CreateNested(ctx context.Context, parentID string, dto *models.ApplicationProfile) error {
	var model models.ApplicationProfile
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}

	return s.createOneShot(ctx, &model)
}
