package models

import (
	"database/sql"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestApplicationVersion_TableName(t *testing.T) {
	require.Equal(t, "application_version_history", (ApplicationVersion{}).TableName())
}

func TestApplicationVersion_PreloadCatalogs(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t,
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationConfiguration{},
		&ApplicationImage{}, &ApplicationProfile{},
		&ApplicationVersion{}, &ApplicationCatalog{},
	)

	tmc := createTerminalModelConfigurationOneShot(t, gdb, "tm-appver-preload", 1)
	ac := createApplicationConfigurationOneShot(t, gdb, "app-appver-preload", "cust-1", tmc.TerminalModelConfigurationId, "pkg.appver.preload")
	prof := createApplicationProfile(t, gdb, ac.ApplicationId, "profile-appver", 7)

	version := ApplicationVersion{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		ExternalApplicationVersionId: sql.NullString{String: "v1", Valid: true},
		VersionName:                  sql.NullString{String: "Version 1", Valid: true},
		VersionCode:                  sql.NullInt64{Int64: 1, Valid: true},
		VersionSize:                  sql.NullInt64{Int64: 100, Valid: true},
	}
	require.NoError(t, gdb.Create(&version).Error)

	catalog := ApplicationCatalog{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		Stage:                        CatalogStageReview,
		ApplicationVersionId:         &version.ApplicationVersionId,
		ApplicationProfileId:         prof.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&catalog).Error)

	var loaded ApplicationVersion
	require.NoError(t, gdb.Preload("ApplicationCatalogs").
		Where("application_version_id = ?", version.ApplicationVersionId).
		First(&loaded).Error)
	require.Len(t, loaded.ApplicationCatalogs, 1)
	require.Equal(t, catalog.Stage, loaded.ApplicationCatalogs[0].Stage)
}
