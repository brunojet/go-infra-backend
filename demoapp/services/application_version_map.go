package services

import (
	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/internal/utils"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
)

type applicationVersionNestedMapper struct{}

func (applicationVersionNestedMapper) ToPostModel(dto dtos.ApplicationVersionPost, model *models.ApplicationVersion) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	return nil
}

func (applicationVersionNestedMapper) ToPatchModel(dto dtos.ApplicationVersionPatch, model *models.ApplicationVersion) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Stage = utils.ToNullInt16(stageMapToModel[dto.Stage])
	return nil
}

func (applicationVersionNestedMapper) ToDTO(model *models.ApplicationVersion, dto *dtos.ApplicationVersionGet) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	debugassert.Assert(dto != nil, "dto cannot be nil")
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

func (applicationVersionNestedMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (m applicationVersionNestedMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	var applicationID, terminalModelConfigurationID int64
	if err := utils.DecodeCompositeKey(parentID, &applicationID, &terminalModelConfigurationID); err != nil {
		return nil, errNestedVersionApplicationIDRequired
	}
	queryScopes[models.ColApplicationID] = applicationID
	queryScopes[models.ColApplicationConfigurationID] = terminalModelConfigurationID
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
