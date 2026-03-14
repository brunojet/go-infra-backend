package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
)

// constants and errors moved to consts.go and errors.go

type ApplicationVersionService interface {
	UpdateAndSyncCatalog(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error
}

type applicationVersionService struct {
	versionRepo repo.ApplicationVersionRepository
	profileRepo repo.ApplicationProfileRepository
	catalogRepo repo.ApplicationCatalogRepository
}

func NewApplicationVersionService(
	versionRepo repo.ApplicationVersionRepository,
	profileRepo repo.ApplicationProfileRepository,
	catalogRepo repo.ApplicationCatalogRepository,
) ApplicationVersionService {
	return &applicationVersionService{versionRepo: versionRepo, profileRepo: profileRepo, catalogRepo: catalogRepo}
}

func (s *applicationVersionService) UpdateAndSyncCatalog(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	return s.versionRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateVersionStageTransition(txCtx, scopes, inOut); err != nil {
			return err
		}
		if err := s.updateVersion(txCtx, scopes, inOut); err != nil {
			return err
		}
		return s.syncCatalogFromVersion(txCtx, *inOut)
	})
}

func (s *applicationVersionService) validateVersionStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	versionID, err := services.ParseScopeInt[int64](scopes, models.ColAppVersionID, 1)
	if err != nil {
		return err
	}

	currentStage, err := s.versionRepo.LoadCurrentStage(ctx, versionID)
	if err != nil {
		return err
	}

	return inOut.ValidateVersionStageTransition(currentStage)
}

func (s *applicationVersionService) updateVersion(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	if err := s.versionRepo.Update(ctx, scopes, inOut); err != nil {
		return err
	}

	return s.versionRepo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationVersionService) syncCatalogFromVersion(ctx context.Context, version models.ApplicationVersion) error {
	if version.Stage.Int16 != models.ApplicationStagePilot && version.Stage.Int16 != models.ApplicationStageProduction {
		return nil
	}

	profileID, err := s.profileRepo.FindCurrentProductionProfileID(ctx, version.ApplicationId)
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

	return s.catalogRepo.Create(ctx, &catalog)
}
