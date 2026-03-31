package services

import (
	"encoding/hex"

	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type applicationImageMapper struct{}

func (m applicationImageMapper) toImageModel(dto dtos.ApplicationImagePost, model *models.ApplicationImage) error {
	if model == nil {
		return errMapperNilModel
	}
	hashBytes, err := hex.DecodeString(dto.FileHash)
	if err != nil {
		return errMapperInvalidFileHash
	}
	model.FileHash = hashBytes
	model.FileName = utils.ToNullString(dto.FileName)
	model.ContentType = utils.ToNullString(dto.ContentType)
	return nil
}

func (m applicationImageMapper) toImageDTO(model *models.ApplicationImage, dto *dtos.ApplicationImageGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
	dto.BaseFile = dtos.BaseFile{
		BaseFilePost: dtos.BaseFilePost{
			FileName:    utils.FromNullString(model.FileName),
			ContentType: utils.FromNullString(model.ContentType),
			FileHash:    hex.EncodeToString(model.FileHash),
		},
		FileStatus: fileStatusFromModel[model.FileStatus],
	}
	switch model.FileStatus {
	case models.ApplicationImageStatusPending, models.ApplicationImageStatusFailed:
		if err := toUploadUrl(model, dto); err != nil {
			return err
		}
	case models.ApplicationImageStatusReady:
		if err := toDownloadUrl(model, dto); err != nil {
			return err
		}
	default:
	}
	return nil
}
