package services

import (
	"context"
	"errors"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/internal/utils"
	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

const (
	nestedProfileParentIDArity                     = 1
	errTextNestedProfileApplicationIDScopeRequired = "application_id parent scope must be valid"
)

var errNestedProfileApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedProfileApplicationIDScopeRequired))

type applicationProfileNestedMapper struct{}

func (applicationProfileNestedMapper) DecodeParentID(parentID string) (map[string]any, error) {
	values, err := utils.DecodeCompactInt64s(parentID, nestedProfileParentIDArity)
	if err != nil {
		return nil, errNestedProfileApplicationIDRequired
	}

	return map[string]any{models.ColAppProfileApplicationID: values[0]}, nil
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

	applicationID, err := repo.RequireScopeInt64(parentScopes, models.ColAppProfileApplicationID, errNestedProfileApplicationIDRequired)
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
	profileID, err := utils.StringToInt64(id)
	if err != nil || profileID <= 0 {
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
