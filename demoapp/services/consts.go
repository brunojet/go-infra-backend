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

	// Version service
	errTextVersionScopeIDRequired = "application_version_id scope must be valid"

	// Nested/version parent error texts
	errTextNestedVersionApplicationIDScopeRequired   = "application_id parent scope must be valid"
	errTextNestedVersionTerminalModelIDScopeRequired = "terminal_model_configuration_id parent scope must be valid"
	errTextNestedVersionIDScopeRequired              = "application_version_id scope must be valid"
)
