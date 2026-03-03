package services

import (
	"context"
	"errors"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	internalservices "github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/internal/utils"
	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

const (
	nestedVersionParentIDArity                       = 2
	errTextNestedVersionApplicationIDScopeRequired   = "application_id parent scope must be valid"
	errTextNestedVersionTerminalModelIDScopeRequired = "terminal_model_configuration_id parent scope must be valid"
	errTextNestedVersionIDScopeRequired              = "application_version_id scope must be valid"
)

var (
	errNestedVersionApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionApplicationIDScopeRequired))
	errNestedVersionTerminalIDRequired    = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionTerminalModelIDScopeRequired))
	errNestedVersionIDRequired            = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionIDScopeRequired))
)

type applicationVersionNestedMapper struct{}

func (applicationVersionNestedMapper) DecodeParentID(parentID string) (map[string]any, error) {
	values, err := utils.DecodeCompactInt64s(parentID, nestedVersionParentIDArity)
	if err != nil {
		return nil, errNestedVersionApplicationIDRequired
	}

	return map[string]any{
		models.ColAppVersionApplicationID:                values[0],
		models.ColAppVersionTerminalModelConfigurationID: values[1],
	}, nil
}

func (applicationVersionNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationVersionNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	mappedQueryScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}

	parentScopes, err := m.DecodeParentID(parentID)
	if err != nil {
		return nil, err
	}

	mergedScopes := make(map[string]any, len(parentScopes)+len(mappedQueryScopes))
	for key, value := range mappedQueryScopes {
		mergedScopes[key] = value
	}
	for key, value := range parentScopes {
		mergedScopes[key] = value
	}

	return mergedScopes, nil
}

func (m applicationVersionNestedMapper) ApplyParentScopes(parentID string, model *models.ApplicationVersion) error {
	parentScopes, err := m.DecodeParentID(parentID)
	if err != nil {
		return err
	}

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

func (applicationVersionNestedMapper) GetModelKey(id string) (map[string]any, error) {
	versionID, err := utils.StringToInt64(id)
	if err != nil || versionID <= 0 {
		return nil, errNestedVersionIDRequired
	}

	return map[string]any{models.ColAppVersionID: versionID}, nil
}

type ApplicationVersionNestedService interface {
	svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
}

type applicationVersionNestedService struct {
	svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]
	repo   repo.ApplicationVersionRepository
	mapper applicationVersionNestedMapper
}

func NewApplicationVersionNestedService(r repo.ApplicationVersionRepository) ApplicationVersionNestedService {
	return &applicationVersionNestedService{
		NestedService: internalservices.NewNestedServiceImpl(r, applicationVersionNestedMapper{}),
		repo:          r,
		mapper:        applicationVersionNestedMapper{},
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

func (s *applicationVersionNestedService) CreateNested(ctx context.Context, parentID string, dto *models.ApplicationVersion) error {
	var model models.ApplicationVersion
	s.mapper.ToModel(dto, &model)
	if err := s.mapper.ApplyParentScopes(parentID, &model); err != nil {
		return err
	}

	return s.createOneShot(ctx, &model)
}
