package services

import (
	"context"
	"errors"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

var errNestedProfileApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New("application_id parent scope must be valid"))

type applicationProfileNestedMapper struct{}

func (applicationProfileNestedMapper) ApplyParentScopes(parentScopes map[string]any, _ *models.ApplicationProfile, model *models.ApplicationProfile) error {
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

type ApplicationProfileNestedService interface {
	svcContracts.NestedService[models.ApplicationProfile, models.ApplicationProfile]
}

type applicationProfileNestedService struct {
	nestedSvc svcContracts.NestedService[models.ApplicationProfile, models.ApplicationProfile]
	repo      repo.ApplicationProfileRepository
	mapper    applicationProfileNestedMapper
}

func NewApplicationProfileNestedService(r repo.ApplicationProfileRepository) ApplicationProfileNestedService {
	return &applicationProfileNestedService{
		nestedSvc: internalservices.NewNestedServiceImpl(r, applicationProfileNestedMapper{}),
		repo:      r,
		mapper:    applicationProfileNestedMapper{},
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

func (s *applicationProfileNestedService) CreateNested(ctx context.Context, parentScopes map[string]any, dto *models.ApplicationProfile) error {
	var model models.ApplicationProfile
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentScopes, dto, &model); err != nil {
		return err
	}

	return s.createOneShot(ctx, &model)
}

func (s *applicationProfileNestedService) ListNested(ctx context.Context, parentScopes map[string]any, params repoContracts.ListParams) ([]models.ApplicationProfile, int64, error) {
	return s.nestedSvc.ListNested(ctx, parentScopes, params)
}
