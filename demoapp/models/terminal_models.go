package models

import (
	"database/sql"

	repositories "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"gorm.io/gorm"
)

const (
	tableTerminalModel              = "terminal_model"
	tableTerminalModelConfiguration = "terminal_model_configuration"

	colTerminalModelName = "name"
	colTerminalModelID   = "terminal_model_id"
	colIntegrationType   = "integration_type"
)

type TerminalModel struct {
	TerminalModelId int64          `gorm:"primaryKey;autoIncrement"`
	Name            sql.NullString `gorm:"not null;size:32;uniqueIndex:ux_terminal_model_name"`
	Description     sql.NullString `gorm:"size:500"`
	CreatedAt       sql.NullTime   `gorm:"autoCreateTime;index:idx_terminal_model_del_created,priority:2"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime;index:idx_terminal_model_del_updated,priority:2"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_terminal_model_del_created,priority:1;index:idx_terminal_model_del_updated,priority:1"`

	TerminalModelConfigurations []TerminalModelConfiguration `gorm:"foreignKey:TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (TerminalModel) TableName() string { return tableTerminalModel }

func (t TerminalModel) BeforeCreate(tx *gorm.DB) error {
	return repositories.AddOnConflictDoNothing(tx, colTerminalModelName)
}

type TerminalModelConfiguration struct {
	TerminalModelConfigurationId int64          `gorm:"primaryKey;autoIncrement"`
	TerminalModelId              int64          `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:1;index:idx_terminal_model_configuration_terminal,priority:1"`
	IntegrationType              sql.NullInt16  `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:2"`
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_terminal_model_configuration_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_terminal_model_configuration_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_terminal_model_configuration_del_created,priority:1;index:idx_terminal_model_configuration_del_updated,priority:1"`

	TerminalModel *TerminalModel
}

func (TerminalModelConfiguration) TableName() string { return tableTerminalModelConfiguration }

func (t TerminalModelConfiguration) BeforeCreate(tx *gorm.DB) error {
	return repositories.AddOnConflictDoNothing(tx, colTerminalModelID, colIntegrationType)
}
