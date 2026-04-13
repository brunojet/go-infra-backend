package models

import (
	"database/sql"
	"errors"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories"
	"gorm.io/gorm"
)

type ApplicationImage struct {
	ApplicationImageId int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId      int64          `gorm:"not null;uniqueIndex:idx_application_image_application,priority:1"`
	FileHash           []byte         `gorm:"type:binary(32);not null;uniqueIndex:idx_application_image_application,priority:3"`
	FileName           sql.NullString `gorm:"not null;size:255"`
	FileSize           sql.NullInt64  `gorm:"not null"`
	FileStatus         sql.NullInt16  `gorm:"not null;default:0;index:idx_application_image_status"` // 0: Pending, 1: Processing, 2: Ready, 3: Failed
	ContentType        sql.NullString `gorm:"not null;size:255"`
	ImageType          sql.NullInt16  `gorm:"not null;default:0;uniqueIndex:idx_application_image_application,priority:2"`
	CreatedAt          sql.NullTime   `gorm:"autoCreateTime;index:idx_application_image_del_created,priority:2"`
	UpdatedAt          sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_image_del_updated,priority:2"`
	DeletedAt          gorm.DeletedAt `gorm:"index:idx_application_image_del_created,priority:1;index:idx_application_image_del_updated,priority:1"`

	// Child-side constraint will live on ApplicationImage.Application
	Application                   *Application
	ApplicationProfiles           []ApplicationProfile           `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshot `gorm:"foreignKey:ApplicationImageId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationImage) TableName() string { return tableApplicationImage }

func (a ApplicationImage) BeforeCreate(tx *gorm.DB) (err error) {
	if !a.ImageType.Valid {
		return errors.New("image_type is required")
	}
	return repositories.AddOnConflictDoNothing(tx,
		ColApplicationID,
		ColFileHash,
		ColAppImageImageType,
	)
}

func (a ApplicationImage) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

// GetOrCreate ensures that an ApplicationImage with the same ApplicationId, FileHash, and ImageType exists in the database.
// If it does not exist, it creates a new record; if it already exists, it loads the existing record into the struct.
//
// This method is intentionally called from models like ApplicationProfile and ApplicationProfileScreenshot,
// since ApplicationImage does not have its own endpoint and is always managed as a dependent resource.
// This guarantees idempotency and avoids duplicate images for the same logical file.
//
// IMPORTANT: Always call this method within an active transaction to ensure consistency when used in composite creations.
func (a *ApplicationImage) GetOrCreate(tx *gorm.DB) error {
	debugassert.Assert(tx != nil, "GetOrCreate must be called with a non-nil transaction")
	tx = tx.Create(a)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		tx = tx.Where("application_id = ? AND file_hash = ? AND image_type = ?", a.ApplicationId, a.FileHash, a.ImageType).First(a)
	}

	return tx.Error
}
