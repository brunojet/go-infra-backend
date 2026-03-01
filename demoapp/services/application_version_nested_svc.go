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

var (
	errNestedVersionApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New("application_id parent scope must be valid"))
	errNestedVersionTerminalIDRequired    = porterrors.NewBusinessRuleError(errors.New("terminal_model_configuration_id parent scope must be valid"))
)

type applicationVersionNestedMapper struct{}

func (applicationVersionNestedMapper) ApplyParentScopes(parentScopes map[string]any, _ *models.ApplicationVersion, model *models.ApplicationVersion) error {
	applicationID, err := repo.RequireScopeInt64(parentScopes, models.ColAppVersionApplicationID, errNestedVersionApplicationIDRequired)
	if err != nil {
		return err
	}

	terminalID, err := repo.RequireScopeInt64(parentScopes, models.ColAppVersionTerminalModelConfigurationID, errNestedVersionTerminalIDRequired)
	if err != nil {
		return err
	}

	model.ApplicationId = applicationID
	model.TerminalModelConfigurationId = terminalID
	return nil
}

func (applicationVersionNestedMapper) ToModel(dto *models.ApplicationVersion, model *models.ApplicationVersion) {
	*model = *dto
}

func (applicationVersionNestedMapper) ToDTO(model *models.ApplicationVersion, dto *models.ApplicationVersion) {
	*dto = *model
}

type ApplicationVersionNestedService interface {
	svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
}

type applicationVersionNestedService struct {
	nestedSvc svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
	repo      repo.ApplicationVersionRepository
	mapper    applicationVersionNestedMapper
}

func NewApplicationVersionNestedService(r repo.ApplicationVersionRepository) ApplicationVersionNestedService {
	return &applicationVersionNestedService{
		nestedSvc: internalservices.NewNestedServiceImpl(r, applicationVersionNestedMapper{}),
		repo:      r,
		mapper:    applicationVersionNestedMapper{},
	}
}

func (s *applicationVersionNestedService) createOneShot(ctx context.Context, inOut *models.ApplicationVersion) error {
	return s.repo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, inOut); err != nil {
			return err
		}
		return s.repo.ArchiveStageDuplicates(txCtx, inOut)
	})
}

func (s *applicationVersionNestedService) CreateNested(ctx context.Context, parentScopes map[string]any, dto *models.ApplicationVersion) error {
	var model models.ApplicationVersion
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentScopes, dto, &model); err != nil {
		return err
	}

	return s.createOneShot(ctx, &model)
}

func (s *applicationVersionNestedService) ListNested(ctx context.Context, parentScopes map[string]any, params repoContracts.ListParams) ([]models.ApplicationVersion, int64, error) {
	return s.nestedSvc.ListNested(ctx, parentScopes, params)
}
