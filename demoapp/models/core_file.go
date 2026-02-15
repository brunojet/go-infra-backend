package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/brunojet/go-infra-backend/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplicationImageType enumerates supported image roles for an application.
type ApplicationImageType int16

const (
	sha256Len = 32
)

const (
	ApplicationImageUnknown    ApplicationImageType = 0
	ApplicationImageIcon       ApplicationImageType = 1
	ApplicationImageBanner     ApplicationImageType = 2
	ApplicationImageScreenshot ApplicationImageType = 3
)

// FileState represents the processing state of a stored file/blob.
type FileState int16

const (
	FileStatePending    FileState = 0
	FileStateProcessing FileState = 1
	FileStateAvailable  FileState = 2
	FileStateFailed     FileState = 3
)

type FileObject struct {
	FileObjectId int64          `gorm:"primaryKey;autoIncrement"`
	Hash         []byte         `gorm:"type:binary(32);not null;uniqueIndex:ux_fileobject_hash_path,priority:1"`
	FilePath     sql.NullString `gorm:"not null;size:128;uniqueIndex:ux_fileobject_hash_path,priority:2"`
	FileName     sql.NullString `gorm:"not null;size:255"`
	FileSize     int64          `gorm:"not null"`
	FileState    FileState      `gorm:"not null;default:0;index:idx_fileobject_state"`
	ContentType  sql.NullString `gorm:"not null;size:128"`
	ErrorMessage sql.NullString `gorm:"size:1024"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index:idx_fileobject_created"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;index:idx_fileobject_updated"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_fileobject_deleted"`
}

func (FileObject) TableName() string { return "file_object" }

func (fo *FileObject) BeforeCreate(tx *gorm.DB) error {
	if len(fo.Hash) != sha256Len {
		return fmt.Errorf("invalid hash length: got %d, want %d", len(fo.Hash), sha256Len)
	}
	if fo.FileName.Valid && fo.FileName.String != "" {
		return nil
	}
	fo.FileName = sql.NullString{String: fileNameFromHash(fo.Hash), Valid: true}
	return nil
}

func createApplicationImageFromHash(db *gorm.DB, applicationID int64, fileType ApplicationImageType, hash []byte, fileSize int64, contentType, filePath, fileName string) (*ApplicationImage, *FileObject, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	existingFO, err := createOrGetFileObject(tx, hash, fileSize, contentType, filePath, fileName)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	ai, err := createOrGetApplicationImage(tx, applicationID, fileType, existingFO.FileObjectId)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return ai, existingFO, nil
}

func createOrGetApplicationImage(tx *gorm.DB, applicationID int64, fileType ApplicationImageType, fileObjectId int64) (*ApplicationImage, error) {
	ai := &ApplicationImage{
		ApplicationId: applicationID,
		FileType:      fileType,
		FileObjectId:  fileObjectId,
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(ai).Error; err != nil {
		return nil, err
	}

	query := map[string]any{
		"application_id": applicationID,
		"file_type":      fileType,
		"file_object_id": fileObjectId,
	}

	var existingAI ApplicationImage
	if err := tx.Model(&ApplicationImage{}).Where(query).First(&existingAI).Error; err != nil {
		return nil, err
	}
	return &existingAI, nil
}

func createOrGetFileObject(tx *gorm.DB, hash []byte, fileSize int64, contentType, filePath, fileName string) (*FileObject, error) {
	fo := &FileObject{
		Hash:        hash,
		FilePath:    utils.ToNullString(filePath),
		FileSize:    fileSize,
		ContentType: utils.ToNullString(contentType),
		FileName:    utils.ToNullString(fileName),
		FileState:   FileStatePending,
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(fo).Error; err != nil {
		return nil, err
	}

	query := map[string]any{
		"hash":      hash,
		"file_path": filePath,
	}

	var existingFO FileObject
	if err := tx.Model(&FileObject{}).Where(query).First(&existingFO).Error; err != nil {
		return nil, err
	}
	return &existingFO, nil
}

func fileNameFromHash(hash []byte) string {
	return base64.RawURLEncoding.EncodeToString(hash)
}

type ApplicationImage struct {
	ApplicationImageId int64                `gorm:"primaryKey;autoIncrement"`
	ApplicationId      int64                `gorm:"not null;index:idx_application_image_app,priority:1;uniqueIndex:ux_application_image_app_type_fileobj,priority:1"`
	FileObjectId       int64                `gorm:"not null;index:idx_application_image_fileobj,priority:1;uniqueIndex:ux_application_image_app_type_fileobj,priority:3"`
	FileType           ApplicationImageType `gorm:"not null;index:idx_application_image_app_type,priority:2;uniqueIndex:ux_application_image_app_type_fileobj,priority:2"`
	CreatedAt          time.Time            `gorm:"autoCreateTime;index:idx_application_image_del_created,priority:2"`
	UpdatedAt          time.Time            `gorm:"autoUpdateTime;index:idx_application_image_del_updated,priority:2"`
	DeletedAt          gorm.DeletedAt       `gorm:"index:idx_application_image_del_created,priority:1;index:idx_application_image_del_updated,priority:1"`

	Application *Application `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	FileObject  *FileObject  `gorm:"foreignKey:FileObjectId;references:FileObjectId"`
}

func (ApplicationImage) TableName() string { return "application_image" }

func CreateApplicationImage(db *gorm.DB, applicationID int64, fileType ApplicationImageType, hash []byte, fileSize int64, contentType, filePath, fileName string) (*ApplicationImage, *FileObject, error) {
	return createApplicationImageFromHash(db, applicationID, fileType, hash, fileSize, contentType, filePath, fileName)
}
