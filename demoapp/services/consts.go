package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/models"
)

// Package-level constants for demoapp services to avoid magic literals.
const (
	// Profile service
	errTextProfileScopeIDRequired = "application_profile_id scope must be valid"
	profileSyncInitialPage        = 1
	profileSyncPageSize           = 10
	profileSyncOrderBy            = models.ColAppVersionTerminalModelConfigurationID
	profileSyncOrder              = "asc"

	// Nested/profile parent handling
	errTextNestedProfileApplicationIDScopeRequired = "application_id parent scope must be valid"

	// Hello world
	helloWorldIDKey = "id"

	// Nested/version parent error texts
	errTextNestedVersionApplicationIDScopeRequired = "application_id parent scope must be valid"
	errTextNestedVersionIDScopeRequired            = "application_version_id scope must be valid"
)

var (
	fileStatusFromModel = map[int16]dtos.FileStatus{
		models.ApplicationImageStatusPending:    dtos.ApplicationImageStatusPending,
		models.ApplicationImageStatusProcessing: dtos.ApplicationImageStatusProcessing,
		models.ApplicationImageStatusReady:      dtos.ApplicationImageStatusReady,
		models.ApplicationImageStatusFailed:     dtos.ApplicationImageStatusFailed,
	}

	fileStatusToModel = map[dtos.FileStatus]int16{
		dtos.ApplicationImageStatusPending:    models.ApplicationImageStatusPending,
		dtos.ApplicationImageStatusProcessing: models.ApplicationImageStatusProcessing,
		dtos.ApplicationImageStatusReady:      models.ApplicationImageStatusReady,
		dtos.ApplicationImageStatusFailed:     models.ApplicationImageStatusFailed,
	}

	stageMapFromModel = map[int16]dtos.ApplicationStage{
		models.ApplicationStagePending:    dtos.ApplicationStagePending,
		models.ApplicationStageReview:     dtos.ApplicationStageReview,
		models.ApplicationStagePilot:      dtos.ApplicationStagePilot,
		models.ApplicationStageProduction: dtos.ApplicationStageProduction,
		models.ApplicationStageArchived:   dtos.ApplicationStageArchived,
	}

	stageMapToModel = map[dtos.ApplicationStage]int16{
		dtos.ApplicationStagePending:    models.ApplicationStagePending,
		dtos.ApplicationStageReview:     models.ApplicationStageReview,
		dtos.ApplicationStagePilot:      models.ApplicationStagePilot,
		dtos.ApplicationStageProduction: models.ApplicationStageProduction,
		dtos.ApplicationStageArchived:   models.ApplicationStageArchived,
	}
)
