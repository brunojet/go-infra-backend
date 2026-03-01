package models

import (
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestApplicationCatalog_TableName(t *testing.T) {
	require.Equal(t, "application_catalog", (ApplicationCatalog{}).TableName())
}

func TestApplicationCatalog_BeforeCreate_UpsertUpdateAll(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t,
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationConfiguration{},
		&ApplicationImage{}, &ApplicationProfile{},
		&ApplicationVersion{}, &ApplicationCatalog{},
	)

	tmc := createTerminalModelConfigurationOneShot(t, gdb, "tm-cat-upd", 1)
	ac := createApplicationConfigurationOneShot(t, gdb, "app-cat-upd", "cust-1", tmc.TerminalModelConfigurationId, "pkg.cat.upd")

	profileV1 := createApplicationProfile(t, gdb, ac.ApplicationId, "profile-v1", 11)
	versionV1 := ApplicationVersion{ApplicationId: ac.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId}
	require.NoError(t, gdb.Create(&versionV1).Error)

	first := ApplicationCatalog{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		Stage:                        catalogStageReview,
		ApplicationVersionId:         &versionV1.ApplicationVersionId,
		ApplicationProfileId:         profileV1.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&first).Error)

	profileV2 := createApplicationProfile(t, gdb, ac.ApplicationId, "profile-v2", 12)
	versionV2 := ApplicationVersion{ApplicationId: ac.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId}
	require.NoError(t, gdb.Create(&versionV2).Error)

	dup := ApplicationCatalog{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		Stage:                        catalogStageReview,
		ApplicationVersionId:         &versionV2.ApplicationVersionId,
		ApplicationProfileId:         profileV2.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&dup).Error)

	var got ApplicationCatalog
	require.NoError(t, gdb.Where("application_id = ? AND terminal_model_configuration_id = ? AND stage = ?", ac.ApplicationId, tmc.TerminalModelConfigurationId, catalogStageReview).
		First(&got).Error)
	require.NotNil(t, got.ApplicationVersionId)
	require.Equal(t, versionV2.ApplicationVersionId, *got.ApplicationVersionId)
	require.Equal(t, profileV2.ApplicationProfileId, got.ApplicationProfileId)
}

func TestApplicationCatalog_Create_DifferentStageCreatesAnotherRow(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t,
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationConfiguration{},
		&ApplicationImage{}, &ApplicationProfile{},
		&ApplicationVersion{}, &ApplicationCatalog{},
	)

	tmc := createTerminalModelConfigurationOneShot(t, gdb, "tm-cat-stage", 1)
	ac := createApplicationConfigurationOneShot(t, gdb, "app-cat-stage", "cust-1", tmc.TerminalModelConfigurationId, "pkg.cat.stage")
	profile := createApplicationProfile(t, gdb, ac.ApplicationId, "profile-stage", 21)
	version := ApplicationVersion{ApplicationId: ac.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId}
	require.NoError(t, gdb.Create(&version).Error)

	stage1 := ApplicationCatalog{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		Stage:                        catalogStageReview,
		ApplicationVersionId:         &version.ApplicationVersionId,
		ApplicationProfileId:         profile.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&stage1).Error)

	stage2 := ApplicationCatalog{
		ApplicationId:                ac.ApplicationId,
		TerminalModelConfigurationId: tmc.TerminalModelConfigurationId,
		Stage:                        catalogStagePilot,
		ApplicationVersionId:         &version.ApplicationVersionId,
		ApplicationProfileId:         profile.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&stage2).Error)

	var count int64
	require.NoError(t, gdb.Model(&ApplicationCatalog{}).
		Where("application_id = ? AND terminal_model_configuration_id = ?", ac.ApplicationId, tmc.TerminalModelConfigurationId).
		Count(&count).Error)
	require.EqualValues(t, 2, count)
}
