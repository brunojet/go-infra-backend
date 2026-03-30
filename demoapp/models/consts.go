package models

// Common
const (
	ColCreatedAt   = "created_at"
	ColUpdatedAt   = "updated_at"
	ColDeletedAt   = "deleted_at"
	ColName        = "name"
	ColDescription = "description"

	ApplicationStagePending    int16 = 0
	ApplicationStageReview     int16 = 1
	ApplicationStagePilot      int16 = 2
	ApplicationStageProduction int16 = 3
	ApplicationStageArchived   int16 = 4
)

// Filter
const (
	ColFilterTypeID = "filter_type_id"
	ColFilterID     = "filter_id"

	IntegrationTypeRFAL int16 = 1
	IntegrationTypeTEF  int16 = 2
)

// Terminal
const (
	ColTerminalModelID              = "terminal_model_id"
	ColTerminalModelConfigurationID = "terminal_model_configuration_id"
)

// Application
const (
	ColApplicationID              = "application_id"
	ColApplicationConfigurationID = "application_configuration_id"
)

// Application version
const (
	versionStagePending    int16 = ApplicationStagePending
	versionStagePilot      int16 = ApplicationStagePilot
	versionStageProduction int16 = ApplicationStageProduction
	versionStageArchived   int16 = ApplicationStageArchived

	ColAppVersionID                           = "application_version_id"
	ColAppVersionApplicationID                = "application_id"
	ColAppVersionTerminalModelConfigurationID = "terminal_model_configuration_id"
	ColAppVersionStage                        = "stage"
)

var (
	versionStageAllowedTransitions = map[int16]map[int16]struct{}{
		versionStagePending: {
			versionStagePilot:    {},
			versionStageArchived: {},
		},
		versionStagePilot: {
			versionStageProduction: {},
			versionStageArchived:   {},
		},
		versionStageProduction: {
			versionStageArchived: {},
		},
		versionStageArchived: {},
	}
)

// Application profile
const (
	tableApplicationProfileScreenshot = "application_profile_screenshot"
	tableApplicationProfileHistory    = "application_profile_history"

	profileStagePending    int16 = ApplicationStagePending
	profileStageReviewed   int16 = ApplicationStageReview
	profileStageProduction int16 = ApplicationStageProduction
	profileStageArchived   int16 = ApplicationStageArchived

	ColAppProfileID            = "application_profile_id"
	ColAppProfileApplicationID = "application_id"
	ColAppProfileStage         = "stage"
	ColAppProfileDeletedAt     = "deleted_at"

	ColAppCatalogApplicationID = "application_id"
	ColAppCatalogProfileID     = "application_profile_id"
	ColAppCatalogStage         = "stage"
)

var (
	profileStageAllowedTransitions = map[int16]map[int16]struct{}{
		profileStagePending: {
			profileStageReviewed: {},
			profileStageArchived: {},
		},
		profileStageReviewed: {
			profileStageProduction: {},
			profileStageArchived:   {},
		},
		profileStageProduction: {
			profileStageArchived: {},
		},
		profileStageArchived: {},
	}
)

// Application catalog
const (
	colCatalogApplicationID                = "application_id"
	colCatalogTerminalModelConfigurationID = "terminal_model_configuration_id"
	colCatalogStage                        = "stage"
	colCatalogApplicationVersionID         = "application_version_id"
	whereCatalogApplicationIDEq            = colCatalogApplicationID + " = ?"
	whereCatalogTerminalModelIDEq          = colCatalogTerminalModelConfigurationID + " = ?"
	whereCatalogStageEq                    = colCatalogStage + " = ?"
	errTextCatalogStageInvalid             = "invalid stage"

	catalogStageReview     int16 = ApplicationStageReview
	catalogStagePilot      int16 = ApplicationStagePilot
	catalogStageProduction int16 = ApplicationStageProduction
)
