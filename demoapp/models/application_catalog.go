package models

import (
	"database/sql"

	"github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"gorm.io/gorm"
)

var (
	errCatalogStageInvalid = errors.NewBusinessRuleError(errTextCatalogStageInvalid)
)

type ApplicationCatalog struct {
	ApplicationId                int64          `gorm:"column:application_id;primaryKey;priority:1;index:idx_ctlg_terminal_stage_app,priority:3"`
	TerminalModelConfigurationId int64          `gorm:"column:terminal_model_configuration_id;primaryKey;priority:2;index:idx_ctlg_terminal_stage_app,priority:1"`
	Stage                        int16          `gorm:"column:stage;primaryKey;priority:3;index:idx_ctlg_terminal_stage_app,priority:2"`
	ApplicationProfileId         int64          `gorm:"column:application_profile_id;not null"`
	ApplicationVersionId         *int64         `gorm:"column:application_version_id"`
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_application_catalog_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_catalog_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_application_catalog_del_created,priority:1;index:idx_application_catalog_del_updated,priority:1"`

	ApplicationProfile *ApplicationProfile
	ApplicationVersion *ApplicationVersion
}

func (ApplicationCatalog) TableName() string { return "application_catalog" }

func (a *ApplicationCatalog) BeforeCreate(tx *gorm.DB) (err error) {
	if _, valid := validCatalogStages[a.Stage]; !valid {
		return errCatalogStageInvalid
	}
	return repositories.AddOnConflictUpdateAll(
		tx,
		ColApplicationID,
		ColApplicationConfigurationID,
		ColStage,
	)
}

func (a ApplicationCatalog) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}
