package models

import (
	"database/sql"
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"gorm.io/gorm"
)

const (
	versionStagePending    int16 = applicationStagePending
	versionStagePilot      int16 = applicationStagePilot
	versionStageProduction int16 = applicationStageProduction
	versionStageArchived   int16 = applicationStageArchived

	colAppVersionID                           = "application_version_id"
	colAppVersionApplicationID                = "application_id"
	colAppVersionTerminalModelConfigurationID = "terminal_model_configuration_id"
	colAppVersionStage                        = "stage"
	colAppVersionDeletedAt                    = "deleted_at"

	ColAppVersionID                           = colAppVersionID
	ColAppVersionApplicationID                = colAppVersionApplicationID
	ColAppVersionTerminalModelConfigurationID = colAppVersionTerminalModelConfigurationID
	ColAppVersionStage                        = colAppVersionStage
	ColAppVersionDeletedAt                    = colAppVersionDeletedAt
)

var (
	errVersionStageInvalid             = porterrors.NewBusinessRuleError(errors.New("stage must be valid"))
	errVersionStageTransitionInvalid   = porterrors.NewBusinessRuleError(errors.New("invalid stage transition"))
	errVersionStageBackwardsTransition = porterrors.NewBusinessRuleError(errors.New("stage cannot transition backwards"))

	versionStageAllowedTransitions = map[int16]map[int16]struct{}{
		versionStagePending: {
			versionStagePilot:    {},
			versionStageArchived: {},
		},
		versionStagePilot: {
			versionStageProduction: {},
			versionStageArchived:   {},
		},
		versionStageProduction: {
			versionStageArchived: {},
		},
		versionStageArchived: {},
	}
)

type ApplicationVersion struct {
	ApplicationVersionId         int64         `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64         `gorm:"column:application_id;index:idx_appver_app_cfg,priority:1;index:idx_appver_terminal_app,priority:2"`
	TerminalModelConfigurationId int64         `gorm:"column:terminal_model_configuration_id;index:idx_appver_app_cfg,priority:2;index:idx_appver_terminal_app,priority:1"`
	Stage                        sql.NullInt16 `gorm:"column:stage;not null;default:0;index:idx_appver_stage_app,priority:2"`
	ReviewAt                     sql.NullTime
	ProductionAt                 sql.NullTime
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_application_version_history_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_version_history_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_application_version_history_del_created,priority:1;index:idx_application_version_history_del_updated,priority:1;index:idx_appver_stage_app,priority:3"`

	ApplicationCatalogs []ApplicationCatalog `gorm:"foreignKey:ApplicationVersionId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationVersion) TableName() string {
	return "application_version_history"
}

func (a ApplicationVersion) ValidateVersionStageTransition(current int16) error {
	target := a.Stage.Int16

	if target == current {
		return nil
	}

	if allowedTargets, ok := versionStageAllowedTransitions[current]; ok {
		if _, allowed := allowedTargets[target]; allowed {
			return nil
		}
	}

	if target < current {
		return errVersionStageBackwardsTransition
	}

	return errVersionStageTransitionInvalid
}

func (a ApplicationVersion) validateRequiredFieldsCreate() error {
	if a.Stage.Valid && a.Stage.Int16 != versionStagePending {
		return errVersionStageInvalid
	}

	return nil
}

func (a ApplicationVersion) BeforeCreate(tx *gorm.DB) error {
	if err := a.validateRequiredFieldsCreate(); err != nil {
		return err
	}

	return nil
}

func (a *ApplicationVersion) BeforeUpdate(tx *gorm.DB) error {
	now := sql.NullTime{Time: tx.NowFunc(), Valid: true}

	switch a.Stage.Int16 {
	case versionStagePilot:
		a.ReviewAt = now
	case versionStageProduction:
		a.ProductionAt = now
	case versionStageArchived:
		a.DeletedAt = gorm.DeletedAt{Time: tx.NowFunc(), Valid: true}
	}

	return nil
}

func (a ApplicationVersion) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}
