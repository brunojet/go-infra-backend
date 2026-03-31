package services

import (
	"context"
	"strconv"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/utils"
	"gorm.io/gorm"
)

type applicationProfileNestedMapper struct {
	im applicationImageMapper
}

func (m applicationProfileNestedMapper) toScreenshotsModel(dtos []dtos.ApplicationProfileScreenshotPostDTO, modelsPtr *[]models.ApplicationProfileScreenshot) error {
	debugassert.Assert(modelsPtr != nil, "modelsPtr cannot be nil")
	screenshots := make([]models.ApplicationProfileScreenshot, len(dtos))
	for i, dto := range dtos {
		modelImage := models.ApplicationImage{}
		if err := m.im.toImageModel(dto.Screenshot, &modelImage); err != nil {
			return err
		}
		screenshots[i].ApplicationImage = &modelImage
		screenshots[i].Position = dto.Position
	}
	*modelsPtr = screenshots
	return nil
}

func (m applicationProfileNestedMapper) toScreenshotsDTO(models []models.ApplicationProfileScreenshot, dtosPtr *[]dtos.ApplicationProfileScreenshotGet) error {
	debugassert.Assert(dtosPtr != nil, "dtosPtr cannot be nil")
	screenshots := make([]dtos.ApplicationProfileScreenshotGet, len(models))
	for i, model := range models {
		if err := m.im.toImageDTO(model.ApplicationImage, &screenshots[i].Screenshot); err != nil {
			return err
		}
		screenshots[i].Position = model.Position
	}
	*dtosPtr = screenshots
	return nil
}

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
	debugassert.Assert(model != nil, "model cannot be nil")
	applicationID, err := services.ParseScopeIntFromString[int64](parentID, 1)
	if err != nil {
		return err
	}
	model.ApplicationId = applicationID
	model.ApplicationImage.ApplicationId = applicationID
	for i := range model.ApplicationProfileScreenshots {
		model.ApplicationProfileScreenshots[i].ApplicationImage.ApplicationId = applicationID
	}
	return nil
}

func toDownloadUrl(model *models.ApplicationImage, dto *dtos.ApplicationImageGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	var downloadDTO dtos.DownloadReady
	downloadDTO.URL = "/applications/" + strconv.FormatInt(model.ApplicationId, 10) + "/applications-images/" + strconv.FormatInt(model.ApplicationImageId, 10) + "/download"
	dto.DownloadReadyDTO = &downloadDTO
	return nil
}

func toUploadUrl(model *models.ApplicationImage, dto *dtos.ApplicationImageGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	var uploadDTO dtos.UploadPending
	uploadDTO.Method = "PUT"
	uploadDTO.URL = "/applications/" + strconv.FormatInt(model.ApplicationId, 10) + "/applications-images/" + strconv.FormatInt(model.ApplicationImageId, 10) + "/upload"
	dto.UploadPendingDTO = &uploadDTO
	return nil
}

// Converte de DTO para Model (POST)
func (m applicationProfileNestedMapper) ToPostModel(dto dtos.ApplicationProfilePost, model *models.ApplicationProfile) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Name = utils.ToNullString(dto.Name)
	model.Description = utils.ToNullString(dto.Description)
	model.ApplicationImage = &models.ApplicationImage{}
	if err := m.im.toImageModel(dto.ApplicationImage, model.ApplicationImage); err != nil {
		return err
	}
	if err := m.toScreenshotsModel(dto.ApplicationProfileScreenshots, &model.ApplicationProfileScreenshots); err != nil {
		return err
	}
	return nil
}

// Converte de DTO para Model (PATCH)
func (applicationProfileNestedMapper) ToPatchModel(dto dtos.ApplicationProfilePatch, model *models.ApplicationProfile) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Stage = utils.ToNullInt16(stageMapToModel[dto.Stage])
	return nil
}

// Converte de Model para DTO
func (m applicationProfileNestedMapper) ToDTO(model *models.ApplicationProfile, dto *dtos.ApplicationProfileGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
	dto.ApplicationProfileId = model.ApplicationProfileId
	dto.ApplicationId = model.ApplicationId
	dto.Stage = stageMapFromModel[model.Stage.Int16]
	dto.Name = utils.FromNullString(model.Name)
	dto.Description = utils.FromNullString(model.Description)
	err := m.im.toImageDTO(model.ApplicationImage, &dto.Icon)
	if err != nil {
		return err
	}
	if err := m.toScreenshotsDTO(model.ApplicationProfileScreenshots, &dto.ApplicationProfileScreenshots); err != nil {
		return err
	}
	dto.ReviewAt = utils.FromNullTimeRFC3339(model.ReviewAt)
	dto.ProductionAt = utils.FromNullTimeRFC3339(model.ProductionAt)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (applicationProfileNestedMapper) GetModelKey(id string) (map[string]any, error) {
	profileID, err := services.ParseScopeIntFromString[int64](id, 1)
	if err != nil {
		return nil, errProfileScopeIDRequired
	}
	return map[string]any{models.ColAppProfileID: profileID}, nil
}

type ApplicationProfileNestedService interface {
	services.NestedService[dtos.ApplicationProfilePost, dtos.ApplicationProfileGet, dtos.ApplicationProfilePatch, models.ApplicationProfile]
}

type applicationProfileNestedService struct {
	services.NestedService[dtos.ApplicationProfilePost, dtos.ApplicationProfileGet, dtos.ApplicationProfilePatch, models.ApplicationProfile]
	pRepo  repositories.ApplicationProfileRepository
	acRepo repositories.ApplicationConfigurationRepository
	vRepo  repositories.ApplicationVersionRepository
	cRepo  repositories.ApplicationCatalogRepository
	mapper applicationProfileNestedMapper
}

func NewApplicationProfileNestedService(p repositories.ApplicationProfileRepository, a repositories.ApplicationConfigurationRepository, v repositories.ApplicationVersionRepository, c repositories.ApplicationCatalogRepository) ApplicationProfileNestedService {
	return &applicationProfileNestedService{
		NestedService: services.NewNestedServiceImpl(p, applicationProfileNestedMapper{}),
		pRepo:         p,
		acRepo:        a,
		vRepo:         v,
		cRepo:         c,
		mapper:        applicationProfileNestedMapper{},
	}
}

// createAndArchive realiza a criação de um ApplicationProfile e garante a associação correta de filtros já existentes.
//
// O fluxo é:
// 1. Salva o profile sem filtros para evitar que o GORM tente criar filtros novos.
// 2. Restaura o slice de filtros original e associa explicitamente os filtros existentes via many2many.
// 3. Arquiva duplicatas de estágio, se houver.
//
// Essa abordagem garante que apenas a associação seja feita, sem risco de violação de constraint UNIQUE ou duplicidade de filtros.
func (s *applicationProfileNestedService) createAndArchive(txCtx context.Context, inOut *models.ApplicationProfile) error {
	debugassert.Assert(inOut != nil, "input model cannot be nil")
	filtersLen := len(inOut.Filters)
	inOut.Filters = inOut.Filters[:0] // ensure GORM doesn't create new filters if the same filter is associated multiple times
	if err := s.pRepo.Create(txCtx, inOut); err != nil {
		return err
	}
	inOut.Filters = inOut.Filters[:filtersLen] // restore original filters for association after creation
	if err := s.pRepo.AssociateFilters(txCtx, inOut); err != nil {
		return err
	}
	return s.pRepo.ArchiveStageDuplicates(txCtx, inOut)
}

func (s *applicationProfileNestedService) CreateNested(ctx context.Context, parentID string, request dtos.ApplicationProfilePost, response *dtos.ApplicationProfileGet) error {
	var model models.ApplicationProfile
	if err := s.mapper.ToPostModel(request, &model); err != nil {
		return err
	}
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}
	return s.pRepo.WithTx(ctx, func(txCtx context.Context) error {
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
	return s.pRepo.WithTx(ctx, func(txCtx context.Context) error {
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
	currentStage, err := s.pRepo.LoadCurrentStage(ctx, scopes)
	if err != nil {
		return err
	}
	return inOut.ValidateProfileStageTransition(currentStage)
}

func (s *applicationProfileNestedService) updateAndArchive(ctx context.Context, scopes map[string]any, inOut *models.ApplicationProfile) error {
	if err := s.pRepo.Update(ctx, scopes, inOut); err != nil {
		return err
	}
	return s.pRepo.ArchiveStageDuplicates(ctx, inOut)
}

func (s *applicationProfileNestedService) listApplicationConfigurationsPage(ctx context.Context, applicationID int64, page int, configs *[]models.ApplicationConfiguration) (int64, error) {
	params := contracts.ListParams{
		QueryParams: contracts.QueryParams{
			Scopes: map[string]any{models.ColAppProfileApplicationID: applicationID},
		},
		Page:    page,
		OrderBy: profileSyncOrderBy,
		Order:   profileSyncOrder,
	}
	totalItems, err := s.acRepo.List(ctx, params, configs)
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
