package services

import "github.com/brunojet/go-infra-backend/demoapp/models"

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
	stageMapFromModel = map[int16]string{
		models.ApplicationStagePending:    "pending",
		models.ApplicationStageReview:     "review",
		models.ApplicationStagePilot:      "pilot",
		models.ApplicationStageProduction: "production",
		models.ApplicationStageArchived:   "archived",
	}

	stageMapToModel = map[string]int16{
		"pending":    models.ApplicationStagePending,
		"review":     models.ApplicationStageReview,
		"pilot":      models.ApplicationStagePilot,
		"production": models.ApplicationStageProduction,
		"archived":   models.ApplicationStageArchived,
	}
)
