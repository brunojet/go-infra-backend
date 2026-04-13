package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
)

type ApplicationVersionNestedService interface {
	svccts.NestedService[dtos.ApplicationVersionPost, dtos.ApplicationVersionGet, dtos.ApplicationVersionPatch]
}

type applicationVersionNestedService struct {
	svccts.NestedService[dtos.ApplicationVersionPost, dtos.ApplicationVersionGet, dtos.ApplicationVersionPatch]
	vRepo  repositories.ApplicationVersionRepository
	pRepo  repositories.ApplicationProfileRepository
	cRepo  repositories.ApplicationCatalogRepository
	mapper applicationVersionNestedMapper
	zero   dtos.ApplicationVersionGet
}

func NewApplicationVersionNestedService(r repositories.ApplicationVersionRepository, p repositories.ApplicationProfileRepository, c repositories.ApplicationCatalogRepository) ApplicationVersionNestedService {
	return &applicationVersionNestedService{
		NestedService: services.NewNestedServiceImpl(r, applicationVersionNestedMapper{}),
		vRepo:         r,
		pRepo:         p,
		cRepo:         c,
		mapper:        applicationVersionNestedMapper{},
	}
}

func (s *applicationVersionNestedService) createAndArchive(txCtx context.Context, inOut *models.ApplicationVersion) error {
	if err := s.vRepo.Create(txCtx, inOut); err != nil {
		return err
	}
	return s.vRepo.ArchiveStageDuplicates(txCtx, inOut)
}

func (s *applicationVersionNestedService) CreateNested(ctx context.Context, parentID string, request dtos.ApplicationVersionPost, response *dtos.ApplicationVersionGet) error {
	var model models.ApplicationVersion
	if err := s.mapper.ToPostModel(request, &model); err != nil {
		return err
	}
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}
	if err := s.vRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.createAndArchive(txCtx, &model); err != nil {
			return err
		}
		return s.mapper.ToDTO(&model, response)
	}); err != nil {
		return err
	}
	return nil
}

func (s *applicationVersionNestedService) updateAndSyncVersion(ctx context.Context, scopes map[string]any, model *models.ApplicationVersion) error {
	if err := s.validateVersionStageTransition(ctx, scopes, model); err != nil {
		return err
	}
	if err := s.updateAndArchive(ctx, scopes, model); err != nil {
		return err
	}
	return s.syncCatalogFromVersion(ctx, *model)
}

func (s *applicationVersionNestedService) Update(ctx context.Context, id string, request dtos.ApplicationVersionPatch, response *dtos.ApplicationVersionGet) error {
	scopes, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	var model models.ApplicationVersion
	if err := s.mapper.ToPatchModel(request, &model); err != nil {
		return err
	}
	if err := s.vRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.updateAndSyncVersion(txCtx, scopes, &model); err != nil {
			return err
		}
		return s.mapper.ToDTO(&model, response)
	}); err != nil {
		return err
	}
	return nil
}

func (s *applicationVersionNestedService) validateVersionStageTransition(ctx context.Context, scopes map[string]any, model *models.ApplicationVersion) error {
	currentStage, err := s.vRepo.LoadCurrentStage(ctx, scopes)
	if err != nil {
		return err
	}

	return model.ValidateVersionStageTransition(currentStage)
}

func (s *applicationVersionNestedService) updateAndArchive(ctx context.Context, scopes map[string]any, model *models.ApplicationVersion) error {
	if err := s.vRepo.Update(ctx, scopes, model); err != nil {
		return err
	}

	return s.vRepo.ArchiveStageDuplicates(ctx, model)
}

func (s *applicationVersionNestedService) syncCatalogFromVersion(ctx context.Context, version models.ApplicationVersion) error {
	if _, valid := models.ValidVersionStagesCatalog[version.Stage.Int16]; !valid {
		return nil // if the version stage is not one that should be in catalog, we skip syncing to catalog
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
