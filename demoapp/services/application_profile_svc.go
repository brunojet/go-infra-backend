package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

// constants and errors moved to consts.go and errors.go

type ApplicationProfileService interface {
	UpdateAndSyncCatalog(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error
}

type applicationProfileService struct {
	profileRepo repo.ApplicationProfileRepository
	appCfgRepo  repoContracts.Repository[models.ApplicationConfiguration]
	catalogRepo repo.ApplicationCatalogRepository
}

func NewApplicationProfileService(
	profileRepo repo.ApplicationProfileRepository,
	appCfgRepo repoContracts.Repository[models.ApplicationConfiguration],
	catalogRepo repo.ApplicationCatalogRepository,
) ApplicationProfileService {
	return &applicationProfileService{profileRepo: profileRepo, appCfgRepo: appCfgRepo, catalogRepo: catalogRepo}
}

func (s *applicationProfileService) UpdateAndSyncCatalog(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	return s.profileRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateProfileStageTransition(txCtx, scopes, inOut); err != nil {
			return err
		}
		if err := s.updateProfile(txCtx, scopes, inOut); err != nil {
			return err
		}
		return s.syncCatalogFromProfile(txCtx, *inOut)
	})
}

func (s *applicationProfileService) validateProfileStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	profileID, err := services.ParseScopeInt[int64](scopes, models.ColAppProfileID, 1)
	if err != nil {
		return err
	}
	currentStage, err := s.profileRepo.LoadCurrentStage(ctx, profileID)
	if err != nil {
		return err
	}

	return inOut.ValidateProfileStageTransition(currentStage)
}

func (s *applicationProfileService) updateProfile(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	if err := s.profileRepo.Update(ctx, scopes, inOut); err != nil {
		return err
	}

	return s.profileRepo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationProfileService) listApplicationConfigurationsPage(ctx context.Context, applicationID int64, page int) ([]models.ApplicationConfiguration, int64, error) {
	params := repoContracts.ListParams{
		QueryParams: repoContracts.QueryParams{
			Scopes: map[string]any{models.ColAppProfileApplicationID: applicationID},
		},
		Page:    page,
		Size:    profileSyncPageSize,
		OrderBy: profileSyncOrderBy,
		Order:   profileSyncOrder,
	}

	return s.appCfgRepo.List(ctx, params)
}

func (s *applicationProfileService) createCatalogsFromProfile(ctx context.Context, profile models.ApplicationProfile, configs []models.ApplicationConfiguration, stage int16) error {
	for _, config := range configs {
		catalog := models.ApplicationCatalog{
			ApplicationId:                config.ApplicationId,
			TerminalModelConfigurationId: config.TerminalModelConfigurationId,
			Stage:                        stage,
			ApplicationProfileId:         profile.ApplicationProfileId,
		}

		if err := s.catalogRepo.Create(ctx, &catalog); err != nil {
			return err
		}
	}

	return nil
}

func (s *applicationProfileService) syncCatalogsFromProfile(ctx context.Context, profile models.ApplicationProfile, stage int16) error {
	page := profileSyncInitialPage
	for {
		configs, total, err := s.listApplicationConfigurationsPage(ctx, profile.ApplicationId, page)
		if err != nil {
			return err
		}
		if len(configs) == 0 {
			return nil
		}
		if err := s.createCatalogsFromProfile(ctx, profile, configs, stage); err != nil {
			return err
		}
		if int64(page*profileSyncPageSize) >= total {
			return nil
		}
		page++
	}
}

func (s *applicationProfileService) syncCatalogsFromProfileWithStages(ctx context.Context, profile models.ApplicationProfile, stages ...int16) error {
	for _, stage := range stages {
		if err := s.syncCatalogsFromProfile(ctx, profile, stage); err != nil {
			return err
		}
	}

	return nil
}

func (s *applicationProfileService) syncCatalogFromProfile(ctx context.Context, profile models.ApplicationProfile) error {
	switch profile.Stage.Int16 {
	case models.ApplicationStageReview:
		return s.syncCatalogsFromProfileWithStages(ctx, profile, models.ApplicationStageReview)
	case models.ApplicationStageProduction:
		if err := s.syncCatalogsFromProfileWithStages(ctx, profile, models.ApplicationStageReview); err != nil {
			return err
		}
		return s.relinkVersionCatalogsToProductionProfile(ctx, profile)
	default:
		return nil
	}
}

func (s *applicationProfileService) relinkVersionCatalogsToProductionProfile(ctx context.Context, profile models.ApplicationProfile) error {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return err
	}

	return tx.Model(&models.ApplicationCatalog{}).
		Where(models.ColAppProfileApplicationID+" = ?", profile.ApplicationId).
		Where(models.ColAppProfileStage+" IN ?", []int16{models.ApplicationStagePilot, models.ApplicationStageProduction}).
		Updates(map[string]any{models.ColAppProfileID: profile.ApplicationProfileId}).Error
}
