package models

import (
	"database/sql"
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	repositories "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"gorm.io/gorm"
)

const (
	colCatalogApplicationID                = "application_id"
	colCatalogTerminalModelConfigurationID = "terminal_model_configuration_id"
	colCatalogStage                        = "stage"
	colCatalogApplicationVersionID         = "application_version_id"
	whereCatalogApplicationIDEq            = colCatalogApplicationID + " = ?"
	whereCatalogTerminalModelIDEq          = colCatalogTerminalModelConfigurationID + " = ?"
	whereCatalogStageEq                    = colCatalogStage + " = ?"
	errTextCatalogVersionRequiredPilotProd = "application_version_id is required when stage is pilot or production"
	errTextCatalogStageInvalid             = "invalid stage"

	catalogStageReview     int16 = applicationStageReview
	catalogStagePilot      int16 = applicationStagePilot
	catalogStageProduction int16 = applicationStageProduction
)

var (
	errCatalogVersionRequiredPilotProd = porterrors.NewBusinessRuleError(errors.New(errTextCatalogVersionRequiredPilotProd))
	errCatalogStageInvalid             = porterrors.NewBusinessRuleError(errors.New(errTextCatalogStageInvalid))
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

	Application                *Application
	TerminalModelConfiguration *TerminalModelConfiguration
	ApplicationProfile         *ApplicationProfile
	ApplicationVersion         *ApplicationVersion
}

func (ApplicationCatalog) TableName() string { return "application_catalog" }

func (a *ApplicationCatalog) BeforeCreate(tx *gorm.DB) (err error) {
	if err := a.hydrateApplicationVersionID(tx); err != nil {
		return err
	}
	if err := a.validateRequiredFieldsCreate(); err != nil {
		return err
	}
	return repositories.AddOnConflictUpdateAll(
		tx,
		colCatalogApplicationID,
		colCatalogTerminalModelConfigurationID,
		colCatalogStage,
	)
}

func (a *ApplicationCatalog) validateRequiredFieldsCreate() error {
	switch a.Stage {
	case catalogStageReview, catalogStagePilot, catalogStageProduction:
		return nil
	default:
		return errCatalogStageInvalid
	}
}

func (a *ApplicationCatalog) hydrateApplicationVersionID(tx *gorm.DB) error {
	if a.ApplicationVersionId != nil {
		return nil
	}

	existing := make([]sql.NullInt64, 0, 1)
	if err := tx.Model(&ApplicationCatalog{}).
		Where(whereCatalogApplicationIDEq, a.ApplicationId).
		Where(whereCatalogTerminalModelIDEq, a.TerminalModelConfigurationId).
		Where(whereCatalogStageEq, a.Stage).
		Limit(1).
		Pluck(colCatalogApplicationVersionID, &existing).Error; err != nil {
		return err
	}

	if len(existing) > 0 && existing[0].Valid {
		existingVersionID := existing[0].Int64
		a.ApplicationVersionId = &existingVersionID
	}

	return nil
}
