package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type ApplicationVersion struct {
	ApplicationVersionId         int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64          `gorm:"column:application_id;index:idx_appver_app_cfg,priority:1;index:idx_appver_terminal_app,priority:2"`
	TerminalModelConfigurationId int64          `gorm:"column:terminal_model_configuration_id;index:idx_appver_app_cfg,priority:2;index:idx_appver_terminal_app,priority:1"`
	ExternalApplicationId        sql.NullString `gorm:"column:ext_app_id;size:36"`
	ExternalApplicationVersionId sql.NullString `gorm:"column:ext_app_vrs_id;not null;size:36"`
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
	return tableApplicationVersionHistory
}

// ValidateVersionStageTransition checks if the transition from the current stage to the target stage is allowed.
//
// Returns nil if the transition is valid, or an error if the transition is not permitted by the allowed transitions map.
// This method enforces business rules for version stage progression and should be called before updating the stage.
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
	return errInvalidStageTransition
}

// BeforeUpdate handles stage transitions for ApplicationVersion.
//
// When a version is archived (stage = Archived), its DeletedAt is set and it becomes immutable.
// This method does NOT allow reactivating ("un-archiving") a version by simply updating its stage.
// To make a version active again, a new publication (new record) must be explicitly created.
// This ensures a clear audit trail and prevents accidental or implicit reactivation of old versions.
func (a *ApplicationVersion) BeforeUpdate(tx *gorm.DB) error {
	now := sql.NullTime{Time: tx.NowFunc(), Valid: true}
	switch a.Stage.Int16 {
	case versionStagePilot:
		if !a.ReviewAt.Valid {
			a.ReviewAt = now
		}
	case versionStageProduction:
		if !a.ProductionAt.Valid {
			a.ProductionAt = now
		}
	case versionStageArchived:
		if !a.DeletedAt.Valid {
			a.DeletedAt = gorm.DeletedAt{Time: now.Time, Valid: true}
		}
	}
	return nil
}

func (a ApplicationVersion) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}
