package models

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	FileContentType    sql.NullString `gorm:"not null;size:255"`
	// Store raw hash bytes (32 bytes). Use binary(32) for DB storage and
	// let GORM handle []byte mapping. Avoid sql.NullByte which doesn't exist.
	FileHash  []byte         `gorm:"type:binary(32);not null;uniqueIndex:idx_application_image_application,priority:3"`
	ImageType sql.NullInt16  `gorm:"not null;uniqueIndex:idx_application_image_application,priority:2"`
	CreatedAt sql.NullTime   `gorm:"autoCreateTime;index:idx_application_image_del_created,priority:2"`
	UpdatedAt sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_image_del_updated,priority:2"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_application_image_del_created,priority:1;index:idx_application_image_del_updated,priority:1"`

	// Child-side constraint will live on ApplicationImage.Application
	Application                   *Application
	ApplicationProfiles           []ApplicationProfile           `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshot `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationImage) TableName() string { return tableApplicationImage }

func (a ApplicationImage) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ApplicationId == 0 {
		return errors.New(errApplicationIDRequired)
	}
	if !a.ImageType.Valid {
		return errors.New(errImageTypeRequired)
	}
	if len(a.FileHash) != 32 {
		return errors.New(errFileHashLength)
	}
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: colAppImageApplicationID}, {Name: colAppImageFileHash}, {Name: colAppImageImageType}},
		DoNothing: true,
	})
	return nil
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
