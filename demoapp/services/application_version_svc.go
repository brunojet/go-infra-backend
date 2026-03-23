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

// Converte de DTO para Model

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

// Converte de DTO para Model
func (applicationVersionNestedMapper) ToModel(_ *dtos.ApplicationVersionDTO, _ *models.ApplicationVersion) {
}

// Converte de Model para DTO
func (applicationVersionNestedMapper) ToDTO(model *models.ApplicationVersion, dto *dtos.ApplicationVersionDTO) {
	dto.ApplicationVersionId = model.ApplicationVersionId
	dto.Stage = utils.FromNullInt16(model.Stage)
	dto.ReviewAt = utils.FromNullTimeRFC3339(model.ReviewAt)
	dto.ProductionAt = utils.FromNullTimeRFC3339(model.ProductionAt)
	dto.CreatedAt = utils.FromNullTimeRFC3339(model.CreatedAt)
	dto.UpdatedAt = utils.FromNullTimeRFC3339(model.UpdatedAt)
	dto.DeletedAt = utils.FromNullTimeRFC3339(model.DeletedAt)
	// parentID (ApplicationConfigurationId) pode ser gerado a partir dos IDs se necessário
}

func (applicationVersionNestedMapper) GetModelKey(id string) (map[string]any, error) {
	versionID, err := services.ParseScopeIntFromString[int64](id, 1)
	if err != nil {
		return nil, errNestedVersionIDRequired
	}
	return map[string]any{models.ColAppVersionID: versionID}, nil
}

type ApplicationVersionNestedService interface {
	svcContracts.NestedService[dtos.ApplicationVersionDTO, models.ApplicationVersion]
}

type applicationVersionNestedService struct {
	svcContracts.NestedService[dtos.ApplicationVersionDTO, models.ApplicationVersion]
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

func (s *applicationVersionNestedService) CreateNested(ctx context.Context, parentID string, dto *dtos.ApplicationVersionDTO) error {
	var model models.ApplicationVersion
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

func (s *applicationVersionNestedService) Update(ctx context.Context, id string, inOut *dtos.ApplicationVersionDTO) error {
	scopes, err := s.mapper.GetModelKey(id)
	if err != nil {
		return err
	}
	var model models.ApplicationVersion
	s.mapper.ToModel(inOut, &model)
	if err := s.vRepo.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.validateVersionStageTransition(txCtx, scopes, &model); err != nil {
			return err
		}
		if err := s.updateVersion(txCtx, scopes, &model); err != nil {
			return err
		}
		return s.syncCatalogFromVersion(txCtx, model)
	}); err != nil {
		return err
	}
	s.mapper.ToDTO(&model, inOut)
	return nil
}

func (s *applicationVersionNestedService) validateVersionStageTransition(ctx context.Context, scopes map[string]any, inOut *models.ApplicationVersion) error {
	currentStage, err := s.vRepo.LoadCurrentStage(ctx, scopes)
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
