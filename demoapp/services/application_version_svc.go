package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	utils "github.com/brunojet/go-infra-backend/pkg/utils"
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

func (applicationVersionNestedMapper) ToPostModel(dto dtos.ApplicationVersionPost, model *models.ApplicationVersion) error {
	if model == nil {
		return errMapperNilModel
	}
	return nil
}

func (applicationVersionNestedMapper) ToPatchModel(dto dtos.ApplicationVersionPatch, model *models.ApplicationVersion) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Stage = utils.ToNullInt16(stageMapToModel[dto.Stage])
	return nil
}

func (applicationVersionNestedMapper) ToDTO(model *models.ApplicationVersion, dto *dtos.ApplicationVersionGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
	applicationConfigurationId, err := utils.EncodeCompositeKey(model.ApplicationId, model.TerminalModelConfigurationId)
	if err != nil {
		return err
	}
	dto.ApplicationVersionId = model.ApplicationVersionId
	dto.ApplicationConfigurationId = applicationConfigurationId
	dto.Stage = stageMapFromModel[model.Stage.Int16]
	dto.ReviewAt = utils.FromNullTimeRFC3339(model.ReviewAt)
	dto.ProductionAt = utils.FromNullTimeRFC3339(model.ProductionAt)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	return nil
}

func (applicationVersionNestedMapper) GetModelKey(id string) (map[string]any, error) {
	versionID, err := services.ParseScopeIntFromString[int64](id, 1)
	if err != nil {
		return nil, errNestedVersionIDRequired
	}
	return map[string]any{models.ColAppVersionID: versionID}, nil
}

type ApplicationVersionNestedService interface {
	svcContracts.NestedService[dtos.ApplicationVersionPost, dtos.ApplicationVersionGet, dtos.ApplicationVersionPatch, models.ApplicationVersion]
}

type applicationVersionNestedService struct {
	svcContracts.NestedService[dtos.ApplicationVersionPost, dtos.ApplicationVersionGet, dtos.ApplicationVersionPatch, models.ApplicationVersion]
	vRepo  repo.ApplicationVersionRepository
	pRepo  repo.ApplicationProfileRepository
	cRepo  repo.ApplicationCatalogRepository
	mapper applicationVersionNestedMapper
	zero   dtos.ApplicationVersionGet
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
