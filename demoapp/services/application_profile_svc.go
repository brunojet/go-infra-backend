package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	rpoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	utils "github.com/brunojet/go-infra-backend/pkg/utils"
	"gorm.io/gorm"
)

type applicationProfileNestedMapper struct{}

func (applicationProfileNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationProfileNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}
	applicationProfileID, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return nil, errNestedProfileApplicationIDRequired
	}
	mappedQueryScopes[models.ColAppProfileApplicationID] = applicationProfileID
	return mappedQueryScopes, nil
}

func (m applicationProfileNestedMapper) ApplyParentScopes(parentID string, model *models.ApplicationProfile) error {
	applicationProfileID, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.ApplicationId = applicationProfileID
	return nil
}

// Converte de DTO para Model
func (applicationProfileNestedMapper) ToModel(dto *dtos.ApplicationProfileDTO, model *models.ApplicationProfile) {
	model.ApplicationId = dto.ApplicationId
	model.Stage = utils.ToNullInt16(dto.Stage)
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	model.ApplicationImageId = dto.ApplicationImageId
}

// Converte de Model para DTO
func (applicationProfileNestedMapper) ToDTO(model *models.ApplicationProfile, dto *dtos.ApplicationProfileDTO) {
	dto.ApplicationProfileId = model.ApplicationProfileId
	dto.ApplicationId = model.ApplicationId
	dto.Stage = utils.FromNullInt16(model.Stage)
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	dto.ApplicationImageId = model.ApplicationImageId
	dto.ReviewAt = utils.FromNullTimeRFC3339(model.ReviewAt)
	dto.ProductionAt = utils.FromNullTimeRFC3339(model.ProductionAt)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	// TODO: Mapear relacionamentos aninhados se necessário
}

func (applicationProfileNestedMapper) GetModelKey(id string) (map[string]any, error) {
	profileID, err := services.ParseScopeIntFromString[int64](id, 1)
	if err != nil {
		return nil, errProfileScopeIDRequired
	}
	return map[string]any{models.ColAppProfileID: profileID}, nil
}

type ApplicationProfileNestedService interface {
	services.NestedService[dtos.ApplicationProfileDTO, models.ApplicationProfile]
}

type applicationProfileNestedService struct {
	services.NestedService[dtos.ApplicationProfileDTO, models.ApplicationProfile]
	pRepo  repo.ApplicationProfileRepository
	acRepo repo.ApplicationConfigurationRepository
	vRepo  repo.ApplicationVersionRepository
	cRepo  repo.ApplicationCatalogRepository
	mapper applicationProfileNestedMapper
}

func NewApplicationProfileNestedService(p repo.ApplicationProfileRepository, a repo.ApplicationConfigurationRepository, v repo.ApplicationVersionRepository, c repo.ApplicationCatalogRepository) ApplicationProfileNestedService {
	return &applicationProfileNestedService{
		NestedService: services.NewNestedServiceImpl(p, applicationProfileNestedMapper{}),
		pRepo:         p,
		acRepo:        a,
		vRepo:         v,
		cRepo:         c,
		mapper:        applicationProfileNestedMapper{},
	}
}

func (s *applicationProfileNestedService) createOneShot(ctx context.Context, inOut *models.ApplicationProfile) error {
	return s.pRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.pRepo.Create(txCtx, inOut); err != nil {
			return err
		}
		return s.pRepo.ArchiveStageDuplicates(txCtx, inOut)
	})
}

func (s *applicationProfileNestedService) CreateNested(ctx context.Context, parentID string, dto *dtos.ApplicationProfileDTO) error {
	var model models.ApplicationProfile
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}
	if err := s.createOneShot(ctx, &model); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, dto)
	return nil
}

func (s *applicationProfileNestedService) Update(ctx context.Context, id string, inOut *dtos.ApplicationProfileDTO) error {
	scopes, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	var model models.ApplicationProfile
	s.mapper.ToModel(inOut, &model)
	if err := s.pRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateProfileStageTransition(txCtx, scopes, &model); err != nil {
			return err
		}
		if err := s.updateProfile(txCtx, scopes, &model); err != nil {
			return err
		}
		return s.syncCatalogFromProfile(txCtx, model)
	}); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, inOut)
	return nil
}

func (s *applicationProfileNestedService) validateProfileStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	currentStage, err := s.pRepo.LoadCurrentStage(ctx, scopes)
	if err != nil {
		return err
	}
	return inOut.ValidateProfileStageTransition(currentStage)
}

func (s *applicationProfileNestedService) updateProfile(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	if err := s.pRepo.Update(ctx, scopes, inOut); err != nil {
		return err
	}
	return s.pRepo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationProfileNestedService) listApplicationConfigurationsPage(ctx context.Context, applicationID int64, page int) ([]models.ApplicationConfiguration, int64, error) {
	params := rpoContracts.ListParams{
		QueryParams: rpoContracts.QueryParams{
			Scopes: map[string]any{models.ColAppProfileApplicationID: applicationID},
		},
		Page:    page,
		Size:    profileSyncPageSize,
		OrderBy: profileSyncOrderBy,
		Order:   profileSyncOrder,
	}
	return s.acRepo.List(ctx, params)
}

func (s *applicationProfileNestedService) createCatalogFromProfile(ctx context.Context, config models.ApplicationConfiguration, profileId, versionId int64, stage int16) error {
	catalog := models.ApplicationCatalog{
		ApplicationId:                config.ApplicationId,
		TerminalModelConfigurationId: config.TerminalModelConfigurationId,
		Stage:                        stage,
		ApplicationProfileId:         profileId,
		ApplicationVersionId:         &versionId,
	}
	if err := s.cRepo.Create(ctx, &catalog); err != nil {
		return err
	}
	return nil
}

func (s *applicationProfileNestedService) createCatalogsFromProfile(ctx context.Context, configs []models.ApplicationConfiguration, profileId int64, stage int16) error {
	for _, config := range configs {
		versionId, err := s.vRepo.FindStageVersionID(ctx, config.ApplicationId, config.TerminalModelConfigurationId, stage)
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
	for {
		configs, total, err := s.listApplicationConfigurationsPage(ctx, profile.ApplicationId, page)
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
