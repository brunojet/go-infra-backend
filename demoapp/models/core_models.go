package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type TerminalModel struct {
	TerminalModelId int64          `gorm:"primaryKey;autoIncrement"`
	Name            sql.NullString `gorm:"not null;uniqueIndex:ux_terminal_model_name"`
	Description     sql.NullString

	TerminalModelConfigurations []TerminalModelConfiguration `gorm:"foreignKey:TerminalModelId"`
}

func (TerminalModel) TableName() string { return "terminal_model" }

type TerminalModelConfiguration struct {
	TerminalModeConfigurationId int64         `gorm:"primaryKey;autoIncrement"`
	TerminalModelId             int64         `gorm:"not null"`
	IntegrationType             sql.NullInt16 `gorm:"not null"`

	TerminalModel             *TerminalModel             `gorm:"foreignKey:TerminalModelId;references:TerminalModelId"`
	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModeConfigurationId"`
}

func (TerminalModelConfiguration) TableName() string { return "terminal_model_configuration" }

type Application struct {
	ApplicationId int64          `gorm:"primaryKey;autoIncrement"`
	Name          sql.NullString `gorm:"not null;uniqueIndex:ux_application_name"`
	Description   sql.NullString
	CreatedAt     sql.NullTime   `gorm:"autoCreateTime;index"`
	UpdatedAt     sql.NullTime   `gorm:"autoUpdateTime;index"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationProfiles       []ApplicationProfile       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationCatalogs       []ApplicationCatalog       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
}

func (Application) TableName() string { return "application" }

type ApplicationConfiguration struct {
	ApplicationId                int64          `gorm:"primaryKey:pk_cfg,priority:1;index:uk_cfg_rev,priority:3;uniqueIndex:uk_cfg_pk,priority:1"`
	TerminalModelConfigurationId int64          `gorm:"primaryKey:pk_cfg,priority:2;index:uk_cfg_rev,priority:2;uniqueIndex:uk_cfg_pk,priority:2"`
	PackageName                  sql.NullString `gorm:"not null;size:255"`

	Application                *Application                `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	TerminalModelConfiguration *TerminalModelConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModeConfigurationId"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }

type ApplicationProfile struct {
	ApplicationProfileId int64 `gorm:"primaryKey;autoIncrement" json:"application_profile_id"`
	ReviewAt             sql.NullTime
	ProductionAt         sql.NullTime

	//Relacionamentos
	ApplicationId int64        `gorm:"index"`
	Application   *Application `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
}

func (ApplicationProfile) TableName() string {
	return "application_profile_history"
}

type ApplicationVersion struct {
	ApplicationVersionId         int64 `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64
	TerminalModelConfigurationId int64

	ReviewAt     sql.NullTime
	ProductionAt sql.NullTime

	Application                Application                 `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	TerminalModelConfiguration *TerminalModelConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModeConfigurationId"`
}

func (ApplicationVersion) TableName() string {
	return "application_version_history"
}

type ApplicationCatalog struct {
	ApplicationId                int64 `gorm:"primaryKey:pk_ctlg,priority:3;index:uk_ctlg_rev,priority:1"`
	TerminalModelConfigurationId int64 `gorm:"primaryKey:pk_ctlg,priority:2;index:uk_ctlg_rev,priority:2"`
	Stage                        int16 `gorm:"primaryKey:pk_ctlg,priority:1;index:uk_ctlg_rev,priority:3"`
	ApplicationVersionId         int64
	ApplicationProfileId         int64
	ApplicationVersion           *ApplicationVersion `gorm:"foreignKey:ApplicationVersionId;references:ApplicationVersionId"`
	ApplicationProfile           *ApplicationProfile `gorm:"foreignKey:ApplicationProfileId;references:ApplicationProfileId"`
}
