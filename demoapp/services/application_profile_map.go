package services

import (
	"strconv"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/pkg/ports/services"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type applicationProfileNestedMapper struct {
	fm filterNestedMapper
	im applicationImageMapper
}

func (m applicationProfileNestedMapper) toFiltersModel(filterIds []int64, model *models.ApplicationProfile) error {
	debugassert.Assert(model != nil, "model cannot be nil")
	model.Filters = make([]models.Filter, len(filterIds))
	m.fm.toAssociativeModel(filterIds, &model.Filters)
	return nil
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
	m.toFiltersModel(dto.FilterIds, model)
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
	mappedQueryScopes[models.ColApplicationID] = applicationProfileID
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
