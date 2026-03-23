package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/utils"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type applicationVersionNestedMapper struct{}

func (applicationVersionNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationVersionNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	var applicationID, terminalModelConfigurationID int64
	if err := utils.DecodeCompositeKey(parentID, &applicationID, &terminalModelConfigurationID); err != nil {
		return nil, errNestedVersionApplicationIDRequired
	}
	queryScopes[models.ColAppVersionApplicationID] = applicationID
	queryScopes[models.ColAppVersionTerminalModelConfigurationID] = terminalModelConfigurationID
	return queryScopes, nil
}

func (m applicationVersionNestedMapper) ApplyParentScopes(parentID string, model *models.ApplicationVersion) error {
	var applicationID, terminalModelConfigurationID int64
	if err := utils.DecodeCompositeKey(parentID, &applicationID, &terminalModelConfigurationID); err != nil {
		return errNestedVersionApplicationIDRequired
	}
	model.ApplicationId = applicationID
	model.TerminalModelConfigurationId = terminalModelConfigurationID
	return nil
}

func (applicationVersionNestedMapper) ToModel(dto *models.ApplicationVersion, model *models.ApplicationVersion) {
	*model = *dto
}

func (applicationVersionNestedMapper) ToDTO(model *models.ApplicationVersion, dto *models.ApplicationVersion) {
	*dto = *model
}

func (applicationVersionNestedMapper) GetModelKey(id string) (map[string]any, error) {
	versionID, err := services.ParseScopeIntFromString[int64](id, 1)
	if err != nil {
		return nil, errNestedVersionIDRequired
	}
	return map[string]any{models.ColAppVersionID: versionID}, nil
}

type ApplicationVersionNestedService interface {
	svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
}

type applicationVersionNestedService struct {
	svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
	vRepo  repo.ApplicationVersionRepository
	pRepo  repo.ApplicationProfileRepository
	cRepo  repo.ApplicationCatalogRepository
	mapper applicationVersionNestedMapper
}

func NewApplicationVersionNestedService(r repo.ApplicationVersionRepository, p repo.ApplicationProfileRepository, c repo.ApplicationCatalogRepository) ApplicationVersionNestedService {
	return &applicationVersionNestedService{
		NestedService: services.NewNestedServiceImpl(r, applicationVersionNestedMapper{}),
		vRepo:         r,
		pRepo:         p,
		cRepo:         c,
		mapper:        applicationVersionNestedMapper{},
	}
}

func (s *applicationVersionNestedService) createOneShot(ctx context.Context, inOut *models.ApplicationVersion) error {
	return s.vRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.vRepo.Create(txCtx, inOut); err != nil {
			return err
		}
		return s.vRepo.ArchiveStageDuplicates(txCtx, inOut)
	})
}

func (s *applicationVersionNestedService) CreateNested(ctx context.Context, parentID string, dto *models.ApplicationVersion) error {
	var model models.ApplicationVersion
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}

	return s.createOneShot(ctx, &model)
}

func (s *applicationVersionNestedService) Update(ctx context.Context, id string, inOut *models.ApplicationVersion) error {
	scopes, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	return s.vRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateVersionStageTransition(txCtx, scopes, inOut); err != nil {
			return err
		}
		if err := s.updateVersion(txCtx, scopes, inOut); err != nil {
			return err
		}
		return s.syncCatalogFromVersion(txCtx, *inOut)
	})
}

func (s *applicationVersionNestedService) validateVersionStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	versionID, err := services.ParseScopeInt[int64](scopes, models.ColAppVersionID, 1)
	if err != nil {
		return err
	}

	currentStage, err := s.vRepo.LoadCurrentStage(ctx, versionID)
	if err != nil {
		return err
	}

	return inOut.ValidateVersionStageTransition(currentStage)
}

func (s *applicationVersionNestedService) updateVersion(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	if err := s.vRepo.Update(ctx, scopes, inOut); err != nil {
		return err
	}

	return s.vRepo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationVersionNestedService) syncCatalogFromVersion(ctx context.Context, version models.ApplicationVersion) error {
	if version.Stage.Int16 != models.ApplicationStagePilot && version.Stage.Int16 != models.ApplicationStageProduction {
		return nil
	}

	profileID, err := s.pRepo.FindCurrentProductionProfileID(ctx, version.ApplicationId)
	if err != nil {
		return err
	}

	catalog := models.ApplicationCatalog{
		ApplicationId:                version.ApplicationId,
		TerminalModelConfigurationId: version.TerminalModelConfigurationId,
		Stage:                        version.Stage.Int16,
		ApplicationProfileId:         profileID,
		ApplicationVersionId:         &version.ApplicationVersionId,
	}

	return s.cRepo.Create(ctx, &catalog)
}
