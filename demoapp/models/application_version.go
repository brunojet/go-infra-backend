package models

import (
	"database/sql"

	"github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"gorm.io/gorm"
)

var (
	errVersionStageTransitionInvalid = errors.NewBusinessRuleError("invalid stage transition")
)

type ApplicationVersion struct {
	ApplicationVersionId         int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64          `gorm:"column:application_id;index:idx_appver_app_cfg,priority:1;index:idx_appver_terminal_app,priority:2"`
	TerminalModelConfigurationId int64          `gorm:"column:terminal_model_configuration_id;index:idx_appver_app_cfg,priority:2;index:idx_appver_terminal_app,priority:1"`
	ExternalApplicationId        sql.NullString `gorm:"column:ext_app_id;size:32"`
	ExternalApplicationVersionId sql.NullString `gorm:"column:ext_app_vrs_id;not null;size:32"`
	VersionName                  sql.NullString `gorm:"column:vrs_name;not null;size:32"`
	VersionCode                  sql.NullInt64  `gorm:"column:vrs_code;not null"`
	VersionSize                  sql.NullInt64  `gorm:"column:vrs_size;not null"`
	Stage                        sql.NullInt16  `gorm:"column:stage;not null;default:0;index:idx_appver_stage_app,priority:2"`
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

	return errVersionStageTransitionInvalid
}

func (a *ApplicationVersion) BeforeUpdate(tx *gorm.DB) error {
	now := sql.NullTime{Time: tx.NowFunc(), Valid: true}
	switch a.Stage.Int16 {
	case versionStagePilot:
		a.ReviewAt = now
	case versionStageProduction:
		a.ProductionAt = now
	case versionStageArchived:
		a.DeletedAt = gorm.DeletedAt{Time: now.Time, Valid: true}
	}
	return nil
}

func (a ApplicationVersion) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}
