package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type TerminalModel struct {
	TerminalModelId int64          `gorm:"primaryKey;autoIncrement"`
	Name            sql.NullString `gorm:"not null;uniqueIndex:ux_terminal_model_name"`
	Description     sql.NullString
	CreatedAt       sql.NullTime   `gorm:"autoCreateTime;index:idx_terminal_model_del_created,priority:2"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime;index:idx_terminal_model_del_updated,priority:2"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_terminal_model_del_created,priority:1;index:idx_terminal_model_del_updated,priority:1"`

	TerminalModelConfigurations []TerminalModelConfiguration `gorm:"foreignKey:TerminalModelId"`
}

func (TerminalModel) TableName() string { return "terminal_model" }

type TerminalModelConfiguration struct {
	TerminalModeConfigurationId int64          `gorm:"primaryKey;autoIncrement"`
	TerminalModelId             int64          `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:1;index:idx_terminal_model_configuration_terminal,priority:1"`
	IntegrationType             sql.NullInt16  `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:2"`
	CreatedAt                   sql.NullTime   `gorm:"autoCreateTime;index:idx_terminal_model_configuration_del_created,priority:2"`
	UpdatedAt                   sql.NullTime   `gorm:"autoUpdateTime;index:idx_terminal_model_configuration_del_updated,priority:2"`
	DeletedAt                   gorm.DeletedAt `gorm:"index:idx_terminal_model_configuration_del_created,priority:1;index:idx_terminal_model_configuration_del_updated,priority:1"`

	TerminalModel             *TerminalModel             `gorm:"foreignKey:TerminalModelId;references:TerminalModelId"`
	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModeConfigurationId"`
}

func (TerminalModelConfiguration) TableName() string { return "terminal_model_configuration" }

type Application struct {
	ApplicationId int64          `gorm:"primaryKey;autoIncrement"`
	Name          sql.NullString `gorm:"not null;uniqueIndex:ux_application_name"`
	Description   sql.NullString
	CreatedAt     sql.NullTime   `gorm:"autoCreateTime;index:idx_application_del_created,priority:2"`
	UpdatedAt     sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_del_updated,priority:2"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_application_del_created,priority:1;index:idx_application_del_updated,priority:1"`

	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationProfiles       []ApplicationProfile       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationCatalogs       []ApplicationCatalog       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
}

func (Application) TableName() string { return "application" }

type ApplicationConfigurationPk struct {
	ApplicationId                int64                       `gorm:"column:application_id;primaryKey,priority:1;index:idx_app_cfg_terminal_app,priority:2"`
	TerminalModelConfigurationId int64                       `gorm:"column:terminal_model_configuration_id;primaryKey,priority:2;index:idx_app_cfg_terminal_app,priority:1"`
	Application                  *Application                `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	TerminalModelConfiguration   *TerminalModelConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModeConfigurationId"`
}

type ApplicationConfiguration struct {
	ApplicationConfigurationPk `gorm:"embedded"`
	PackageName                sql.NullString `gorm:"column:package_name;not null;size:255"`
	CreatedAt                  sql.NullTime   `gorm:"autoCreateTime;index:idx_application_configuration_del_created,priority:2"`
	UpdatedAt                  sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_configuration_del_updated,priority:2"`
	DeletedAt                  gorm.DeletedAt `gorm:"index:idx_application_configuration_del_created,priority:1;index:idx_application_configuration_del_updated,priority:1"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }

type ApplicationProfile struct {
	ApplicationProfileId int64 `gorm:"primaryKey;autoIncrement"`
	ApplicationId        int64 `gorm:"index"`
	ReviewAt             sql.NullTime
	ProductionAt         sql.NullTime
	CreatedAt            sql.NullTime   `gorm:"autoCreateTime;index:idx_application_profile_history_del_created,priority:2"`
	UpdatedAt            sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_profile_history_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_application_profile_history_del_created,priority:1;index:idx_application_profile_history_del_updated,priority:1"`

	//Relacionamentos
	Application *Application `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
}

func (ApplicationProfile) TableName() string {
	return "application_profile_history"
}

type ApplicationVersion struct {
	ApplicationVersionId         int64 `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64 `gorm:"column:application_id;index:idx_appver_app_cfg,priority:1;index:idx_appver_terminal_app,priority:2"`
	TerminalModelConfigurationId int64 `gorm:"column:terminal_model_configuration_id;index:idx_appver_app_cfg,priority:2;index:idx_appver_terminal_app,priority:1"`
	ReviewAt                     sql.NullTime
	ProductionAt                 sql.NullTime
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_application_version_history_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_version_history_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_application_version_history_del_created,priority:1;index:idx_application_version_history_del_updated,priority:1"`

	ApplicationConfiguration *ApplicationConfiguration `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;references:ApplicationId,TerminalModelConfigurationId"`
}

func (ApplicationVersion) TableName() string {
	return "application_version_history"
}

type ApplicationCatalogPk struct {
	ApplicationId                int64 `gorm:"column:application_id;primaryKey,priority:1;index:idx_ctlg_terminal_stage_app,priority:3"`
	TerminalModelConfigurationId int64 `gorm:"column:terminal_model_configuration_id;primaryKey,priority:2;index:idx_ctlg_terminal_stage_app,priority:1"`
	Stage                        int16 `gorm:"column:stage;primaryKey,priority:3;index:idx_ctlg_terminal_stage_app,priority:2"`
}

type ApplicationCatalog struct {
	ApplicationCatalogPk `gorm:"embedded"`
	ApplicationVersionId int64          `gorm:"column:application_version_id"`
	ApplicationProfileId int64          `gorm:"column:application_profile_id"`
	CreatedAt            sql.NullTime   `gorm:"autoCreateTime;index:idx_application_catalog_del_created,priority:2"`
	UpdatedAt            sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_catalog_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_application_catalog_del_created,priority:1;index:idx_application_catalog_del_updated,priority:1"`

	ApplicationConfiguration *ApplicationConfiguration `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;references:ApplicationId,TerminalModelConfigurationId"`
	ApplicationVersion       *ApplicationVersion       `gorm:"foreignKey:ApplicationVersionId;references:ApplicationVersionId"`
	ApplicationProfile       *ApplicationProfile       `gorm:"foreignKey:ApplicationProfileId;references:ApplicationProfileId"`
}

func (ApplicationCatalog) TableName() string { return "application_catalog" }
