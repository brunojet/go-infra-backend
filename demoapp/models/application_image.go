package models

import (
	"database/sql"
	"errors"

	repositories "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"gorm.io/gorm"
)

const (
	tableApplicationImage = "application_image"

	errApplicationIDRequired = "application_id must be valid"
	errImageTypeRequired     = "image_type must be valid"
	errFileHashLength        = "file_hash must be exactly 32 bytes"
	errTransactionRequired   = "transaction is required"

	colAppImageApplicationID = "application_id"
	colAppImageFileHash      = "file_hash"
	colAppImageImageType     = "image_type"
)

type ApplicationImage struct {
	ApplicationImageId int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId      int64          `gorm:"not null;uniqueIndex:idx_application_image_application,priority:1"`
	FileName           sql.NullString `gorm:"not null;size:255"`
	FileHash           []byte         `gorm:"type:binary(32);not null;uniqueIndex:idx_application_image_application,priority:3"`
	ContentType        sql.NullString `gorm:"not null;size:255"`
	ImageType          sql.NullInt16  `gorm:"not null;uniqueIndex:idx_application_image_application,priority:2"`
	Status             int16          `gorm:"not null;default:0;index:idx_application_image_status"` // 0: Pending, 1: Processing, 2: Ready, 3: Failed
	CreatedAt          sql.NullTime   `gorm:"autoCreateTime;index:idx_application_image_del_created,priority:2"`
	UpdatedAt          sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_image_del_updated,priority:2"`
	DeletedAt          gorm.DeletedAt `gorm:"index:idx_application_image_del_created,priority:1;index:idx_application_image_del_updated,priority:1"`

	// Child-side constraint will live on ApplicationImage.Application
	Application                   *Application
	ApplicationProfiles           []ApplicationProfile           `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshot `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// Status constants for ApplicationImage
const (
	ApplicationImageStatusPending    int16 = 0
	ApplicationImageStatusProcessing int16 = 1
	ApplicationImageStatusReady      int16 = 2
	ApplicationImageStatusFailed     int16 = 3
)

func (ApplicationImage) TableName() string { return tableApplicationImage }

func (a ApplicationImage) BeforeCreate(tx *gorm.DB) (err error) {
	return repositories.AddOnConflictDoNothing(tx,
		colAppImageApplicationID,
		colAppImageFileHash,
		colAppImageImageType,
	)
}

func (a ApplicationImage) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return repositories.WhereOnConflict(tx,
		repositories.ConflictScope{ColumnName: colAppImageApplicationID, ColumnValue: a.ApplicationId},
		repositories.ConflictScope{ColumnName: colAppImageFileHash, ColumnValue: a.FileHash},
		repositories.ConflictScope{ColumnName: colAppImageImageType, ColumnValue: a.ImageType.Int16},
	)
}

func (a *ApplicationImage) GetOrCreate(tx *gorm.DB) error {
	if tx == nil {
		return errors.New(errTransactionRequired)
	}

	tx = tx.Create(a)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		tx = tx.Where("application_id = ? AND file_hash = ? AND image_type = ?", a.ApplicationId, a.FileHash, a.ImageType).First(a)
	}

	return tx.Error
}
