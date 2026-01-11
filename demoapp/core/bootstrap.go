package core

import (
	coremodels "github.com/brunojet/go-infra-backend/demoapp/core/models"
	"gorm.io/gorm"
)

// Register performs core-level initialization.
//
// Today it runs DB migrations for all models owned by core.
func Register(db *gorm.DB) error {
	return db.AutoMigrate(coremodels.MigratableModels()...)
}
