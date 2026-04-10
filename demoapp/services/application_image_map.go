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
	model.FileSize = utils.ToNullInt64(dto.FileSize)
	model.ContentType = utils.ToNullString(dto.ContentType)
	return nil
}

func (m applicationImageMapper) toImageDTO(model *models.ApplicationImage, dto *dtos.ApplicationImageGet) error {
	if model == nil || dto == nil {
		return errMapperNilModel
	}
	dto.BaseFile = dtos.BaseFile{
		BaseFilePost: dtos.BaseFilePost{
			FileHash:    hex.EncodeToString(model.FileHash),
			FileName:    utils.FromNullString(model.FileName),
			FileSize:    utils.FromNullInt64(model.FileSize),
			ContentType: utils.FromNullString(model.ContentType),
		},
		FileStatus: fileStatusFromModel[model.FileStatus.Int16],
	}
	switch model.FileStatus.Int16 {
	case models.FileStatusPending, models.FileStatusFailed:
		if err := toUploadUrl(model, dto); err != nil {
			return err
		}
	case models.FileStatusReady:
		if err := toDownloadUrl(model, dto); err != nil {
			return err
		}
	default:
	}
	return nil
}
