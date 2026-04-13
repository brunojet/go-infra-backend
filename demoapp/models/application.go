package models

import (
	"database/sql"

	repositories "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories"
	"gorm.io/gorm"
)

const ()

type Application struct {
	ApplicationId int64          `gorm:"primaryKey;autoIncrement"`
	CustomerId    sql.NullString `gorm:"size:36;not null;index:idx_application_customer"`
	Name          sql.NullString `gorm:"not null;size:32;uniqueIndex:ux_application_name"`
	Description   sql.NullString `gorm:"size:500"`
	CreatedAt     sql.NullTime   `gorm:"autoCreateTime;index:idx_application_del_created,priority:2"`
	UpdatedAt     sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_del_updated,priority:2"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_application_del_created,priority:1;index:idx_application_del_updated,priority:1"`

	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationProfiles       []ApplicationProfile       `gorm:"foreignKey:ApplicationId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationImages         []ApplicationImage         `gorm:"foreignKey:ApplicationId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (Application) TableName() string { return tableApplication }

func (a Application) BeforeCreate(tx *gorm.DB) error {
	return repositories.AddOnConflictDoNothing(tx, ColName)
}

func (a Application) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return repositories.WhereOnConflict(tx, repositories.ConflictScope{ColumnName: ColName, ColumnValue: a.Name})
}

type ApplicationConfiguration struct {
	ApplicationId                int64          `gorm:"column:application_id;primaryKey;priority:1;index:idx_app_cfg_terminal_app,priority:2"`
	TerminalModelConfigurationId int64          `gorm:"column:terminal_model_configuration_id;primaryKey;priority:2;index:idx_app_cfg_terminal_app,priority:1;uniqueIndex:ux_app_cfg_terminal_app,priority:1"`
	PackageName                  sql.NullString `gorm:"column:package_name;not null;size:255;uniqueIndex:ux_app_cfg_terminal_app,priority:2"`
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_application_configuration_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_configuration_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_application_configuration_del_created,priority:1;index:idx_application_configuration_del_updated,priority:1"`

	Application                *Application
	TerminalModelConfiguration *TerminalModelConfiguration
	ApplicationVersions        []ApplicationVersion `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationCatalogs        []ApplicationCatalog `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationConfiguration) TableName() string { return tableApplicationConfiguration }

func (a ApplicationConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	return repositories.AddOnConflictDoNothing(tx, ColTerminalModelConfigurationID, ColPackageName)
}

func (a ApplicationConfiguration) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return repositories.WhereOnConflict(tx,
		repositories.ConflictScope{ColumnName: ColTerminalModelConfigurationID, ColumnValue: a.TerminalModelConfigurationId},
		repositories.ConflictScope{ColumnName: ColPackageName, ColumnValue: a.PackageName},
	)
}
