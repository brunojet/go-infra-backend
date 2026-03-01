package models

import (
	"database/sql"
	"errors"

	repositories "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"gorm.io/gorm"
)

const (
	tableApplication              = "application"
	tableApplicationConfiguration = "application_configuration"

	errNameAndCustomerRequired  = "name and customer_id must be valid"
	errAppIDAndCustomerRequired = "application_id and customer_id must be valid"
	errPackageNameRequired      = "package_name must be valid"

	colName                         = "name"
	colCustomerID                   = "customer_id"
	colApplicationID                = "application_id"
	colPackageName                  = "package_name"
	colTerminalModelConfigurationID = "terminal_model_configuration_id"

	whereNameEq          = "name = ?"
	whereApplicationIDEq = "application_id = ?"
	wherePackageNameEq   = "package_name = ?"
)

type Application struct {
	ApplicationId int64          `gorm:"primaryKey;autoIncrement"`
	CustomerId    sql.NullString `gorm:"size:32;not null;index:idx_application_customer"`
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

// validateAppOwnerCreate checks whether there is an existing Application with the
// same name owned by a different customer. candidateCustomer is the customer
// identifier to validate (useful for Create and Update flows). It requires the
// caller to run inside a transaction so SELECT ... FOR UPDATE is effective.
func (a Application) validateAppOwnerCreate(tx *gorm.DB) error {
	if !a.Name.Valid || !a.CustomerId.Valid {
		return errors.New(errNameAndCustomerRequired)
	}
	candidateCustomer := a.CustomerId.String
	return repositories.ValidateTxWithUpdateLock(tx, repositories.LockValidationSpec[Application]{
		SelectColumns: []string{colCustomerID},
		WhereSQL:      whereNameEq,
		WhereArgs:     []any{a.Name.String},
		BlockIfFound: func(existing *Application) error {
			if existing.CustomerId.String != candidateCustomer {
				return gorm.ErrCheckConstraintViolated
			}
			return nil
		},
	})
}

// ValidateAppOwnerUpdate validates updates that would change customer ownership or name.
// It uses the Application instance's `CustomerId` as the candidate value and
// reuses ValidateCustomer without reflection.
func (a Application) validateAppOwnerUpdate(tx *gorm.DB) error {
	if a.ApplicationId == 0 || !a.CustomerId.Valid {
		return errors.New(errAppIDAndCustomerRequired)
	}
	candidateCustomer := a.CustomerId.String
	if err := repositories.ValidateTxWithUpdateLock(tx, repositories.LockValidationSpec[Application]{
		SelectColumns: []string{colCustomerID},
		WhereSQL:      whereApplicationIDEq,
		WhereArgs:     []any{a.ApplicationId},
		BlockIfFound: func(existing *Application) error {
			if existing.CustomerId.String != candidateCustomer {
				return gorm.ErrCheckConstraintViolated
			}
			return nil
		},
	}); err != nil {
		return err
	}

	// Also enforce name ownership rule during updates (name collision across
	// different customers must be blocked).
	return a.validateAppOwnerCreate(tx)
}

func (a Application) BeforeCreate(tx *gorm.DB) error {
	if err := a.validateAppOwnerCreate(tx); err != nil {
		return err
	}

	return repositories.AddOnConflictDoNothing(tx, colName)
}

func (a *Application) BeforeUpdate(tx *gorm.DB) error {
	if err := a.validateAppOwnerUpdate(tx); err != nil {
		return err
	}
	return nil
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
	Filters                    []Filter             `gorm:"many2many:application_configuration_filter;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (ApplicationConfiguration) TableName() string { return tableApplicationConfiguration }

func (a ApplicationConfiguration) validateAppConfigurationOwner(tx *gorm.DB) error {
	if !a.PackageName.Valid {
		return errors.New(errPackageNameRequired)
	}

	return repositories.ValidateTxWithUpdateLock(tx, repositories.LockValidationSpec[ApplicationConfiguration]{
		SelectColumns: []string{colApplicationID},
		WhereSQL:      wherePackageNameEq,
		WhereArgs:     []any{a.PackageName.String},
		BlockIfFound: func(existing *ApplicationConfiguration) error {
			if existing.ApplicationId != a.ApplicationId {
				return gorm.ErrCheckConstraintViolated
			}
			return nil
		},
	})
}

func (a ApplicationConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	if err := a.validateAppConfigurationOwner(tx); err != nil {
		return err
	}

	return repositories.AddOnConflictDoNothing(tx, colTerminalModelConfigurationID, colPackageName)
}

func (a *ApplicationConfiguration) BeforeUpdate(tx *gorm.DB) error {
	if err := a.validateAppConfigurationOwner(tx); err != nil {
		return err
	}
	return nil
}
