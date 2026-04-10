package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"gorm.io/gorm"
)

type ApplicationProfileNestedService interface {
	services.NestedService[dtos.ApplicationProfilePost, dtos.ApplicationProfileGet, dtos.ApplicationProfilePatch]
}

// ApplicationProfileNestedService implements the lifecycle for ApplicationProfile entities.
//
// The entity is created in a single atomic step (CreateNested), including all associations and filters.
// After creation, the profile only advances through stages (e.g., Review, Pilot, Production) via updates.
// No further creation or association of filters is performed after the initial creation; only stage transitions are allowed.
// This design ensures data consistency and prevents duplicate associations.
type applicationProfileNestedService struct {
	services.NestedService[dtos.ApplicationProfilePost, dtos.ApplicationProfileGet, dtos.ApplicationProfilePatch]
	aprpo  repositories.ApplicationProfileRepository
	acrpo  repositories.ApplicationConfigurationRepository
	avrpo  repositories.ApplicationVersionRepository
	crpo   repositories.ApplicationCatalogRepository
	mapper applicationProfileNestedMapper
}

func NewApplicationProfileNestedService(p repositories.ApplicationProfileRepository, a repositories.ApplicationConfigurationRepository, v repositories.ApplicationVersionRepository, c repositories.ApplicationCatalogRepository) ApplicationProfileNestedService {
	return &applicationProfileNestedService{
		NestedService: services.NewNestedServiceImpl(p, applicationProfileNestedMapper{}),
		aprpo:         p,
		acrpo:         a,
		avrpo:         v,
		crpo:          c,
		mapper:        applicationProfileNestedMapper{},
	}
}

// createAndArchive creates an ApplicationProfile and ensures correct association of existing filters.
//
// Flow:
// 1. Saves the profile without filters to prevent GORM from creating new filters.
// 2. Restores the original filters slice and explicitly associates existing filters via many2many.
// 3. Archives previous profiles with the same stage (if any), so only the latest remains active.
//
// This approach guarantees only association (not creation) of filters, avoiding UNIQUE constraint violations or duplicates.
//
// Note: This method does not provide traditional idempotency protection. If the same or even different profiles are submitted multiple times, previous ones for the same stage will be archived, and only the latest remains active.
func (s *applicationProfileNestedService) createAndArchive(txCtx context.Context, inOut *models.ApplicationProfile) error {
	debugassert.Assert(inOut != nil, "input model cannot be nil")
	filtersLen := len(inOut.Filters)
	inOut.Filters = inOut.Filters[:0] // ensure GORM doesn't create new filters if the same filter is associated multiple times
	if err := s.aprpo.Create(txCtx, inOut); err != nil {
		return err
	}
	inOut.Filters = inOut.Filters[:filtersLen] // restore original filters for association after creation
	if err := s.aprpo.AssociateFilters(txCtx, inOut); err != nil {
		return err
	}
	return s.aprpo.ArchiveStageDuplicates(txCtx, inOut)
}

func (s *applicationProfileNestedService) CreateNested(ctx context.Context, parentID string, request dtos.ApplicationProfilePost, response *dtos.ApplicationProfileGet) error {
	var model models.ApplicationProfile
	if err := s.mapper.ToPostModel(request, &model); err != nil {
		return err
	}
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}
	return s.aprpo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.createAndArchive(txCtx, &model); err != nil {
			return err
		}
		if err := s.mapper.ToDTO(&model, response); err != nil {
			return err
		}
		return nil
	})
}

func (s *applicationProfileNestedService) Update(ctx context.Context, id string, request dtos.ApplicationProfilePatch, response *dtos.ApplicationProfileGet) error {
	scopes, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	var model models.ApplicationProfile
	if err := s.mapper.ToPatchModel(request, &model); err != nil {
		return err
	}
	return s.aprpo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateProfileStageTransition(txCtx, scopes, &model); err != nil {
			return err
		}
		if err := s.updateAndArchive(txCtx, scopes, &model); err != nil {
			return err
		}
		if err := s.syncCatalogFromProfile(txCtx, model); err != nil {
			return err
		}
		if err := s.mapper.ToDTO(&model, response); err != nil {
			return err
		}
		return nil
	})
}

func (s *applicationProfileNestedService) validateProfileStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	currentStage, err := s.aprpo.LoadCurrentStage(ctx, scopes)
	if err != nil {
		return err
	}
	return inOut.ValidateProfileStageTransition(currentStage)
}

func (s *applicationProfileNestedService) updateAndArchive(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	if err := s.aprpo.Update(ctx, scopes, inOut); err != nil {
		return err
	}
	return s.aprpo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationProfileNestedService) listApplicationConfigurationsPage(ctx context.Context, applicationID int64, page int, configs *[]models.ApplicationConfiguration) (int64, error) {
	params := contracts.ListParams{
		QueryParams: contracts.QueryParams{
			Scopes: map[string]any{models.ColApplicationID: applicationID},
		},
		Page:    page,
		OrderBy: profileSyncOrderBy,
		Order:   profileSyncOrder,
	}
	totalItems, err := s.acrpo.List(ctx, params, configs)
	if err != nil {
		return 0, err
	}
	return totalItems, nil
}

func (s *applicationProfileNestedService) createCatalogFromProfile(ctx context.Context, config models.ApplicationConfiguration, profileId, versionId int64, stage int16) error {
	catalog := models.ApplicationCatalog{
		ApplicationId:                config.ApplicationId,
		TerminalModelConfigurationId: config.TerminalModelConfigurationId,
		Stage:                        stage,
		ApplicationProfileId:         profileId,
		ApplicationVersionId:         &versionId,
	}
	if err := s.crpo.Create(ctx, &catalog); err != nil {
		return err
	}
	return nil
}

func (s *applicationProfileNestedService) createCatalogsFromProfile(ctx context.Context, configs []models.ApplicationConfiguration, profileId int64, stage int16) error {
	for _, config := range configs {
		versionId, err := s.avrpo.FindStageVersionID(ctx, config.ApplicationId, config.TerminalModelConfigurationId, stage)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err := s.createCatalogFromProfile(ctx, config, profileId, versionId, stage); err != nil {
			return err
		}
	}
	return nil
}

func (s *applicationProfileNestedService) syncCatalogsFromProfile(ctx context.Context, profile models.ApplicationProfile, stage int16) error {
	page := profileSyncInitialPage
	configs := make([]models.ApplicationConfiguration, 0, profileSyncPageSize)
	for {
		total, err := s.listApplicationConfigurationsPage(ctx, profile.ApplicationId, page, &configs)
		if err != nil {
			return err
		}
		if len(configs) == 0 {
			return nil
		}
		if err := s.createCatalogsFromProfile(ctx, configs, profile.ApplicationProfileId, stage); err != nil {
			return err
		}
		if int64(page*profileSyncPageSize) >= total {
			return nil
		}
		configs = configs[:0] // reset slice while keeping allocated memory
		page++
	}
}

func (s *applicationProfileNestedService) syncCatalogsFromProfileWithStages(ctx context.Context, profile models.ApplicationProfile, stages ...int16) error {
	for _, stage := range stages {
		if err := s.syncCatalogsFromProfile(ctx, profile, stage); err != nil {
			return err
		}
	}
	return nil
}

func (s *applicationProfileNestedService) syncCatalogFromProfile(ctx context.Context, profile models.ApplicationProfile) error {
	switch profile.Stage.Int16 {
	case models.ApplicationStageReview:
		return s.syncCatalogsFromProfileWithStages(ctx, profile, models.ApplicationStageReview)
	case models.ApplicationStageProduction:
		return s.syncCatalogsFromProfileWithStages(ctx, profile, models.ApplicationStagePilot, models.ApplicationStageProduction)
	default:
		return nil
	}
}
