package models

import (
	"database/sql"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestApplicationImage_TableName(t *testing.T) {
	require.Equal(t, tableApplicationImage, (ApplicationImage{}).TableName())
}

func TestApplicationImage_BeforeCreate_ValidationBranches(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{})
	app := createApplication(t, gdb, "app-img-validation", "cust-1")
	invalidApplication := newApplicationImage(0, 1, 1)
	err := gdb.Create(&invalidApplication).Error
	require.Error(t, err)

	invalidImageType := newApplicationImage(app.ApplicationId, 1, 4)
	invalidImageType.ImageType = sql.NullInt16{}
	err = gdb.Create(&invalidImageType).Error
	require.Error(t, err)
}

func TestApplicationImage_BeforeCreate_UpsertDoNothingCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{})
	app := createApplication(t, gdb, "app-img-upsert", "cust-1")
	first := createApplicationImage(t, gdb, app.ApplicationId, 1, 9)

	second := newApplicationImage(app.ApplicationId, 1, 9)
	second.FileName = sql.NullString{String: "other-name.png", Valid: true}
	require.NoError(t, gdb.Create(&second).Error)

	var count int64
	require.NoError(t, gdb.Model(&ApplicationImage{}).
		Where("application_id = ? AND file_hash = ? AND image_type = ?", app.ApplicationId, first.FileHash, 1).
		Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestApplicationImage_CreateOrGet_CreateAndGetBranches(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{})
	app := createApplication(t, gdb, "app-img-cog", "cust-1")

	created := newApplicationImage(app.ApplicationId, 2, 7)
	created.FileName = sql.NullString{String: "icon-a.png", Valid: true}
	require.NoError(t, created.GetOrCreate(gdb))
	require.NotZero(t, created.ApplicationImageId)

	duplicate := newApplicationImage(app.ApplicationId, 2, 7)
	duplicate.FileName = sql.NullString{String: "icon-b.png", Valid: true}
	require.NoError(t, duplicate.GetOrCreate(gdb))
	require.Equal(t, created.ApplicationImageId, duplicate.ApplicationImageId)

	var count int64
	require.NoError(t, gdb.Model(&ApplicationImage{}).
		Where("application_id = ? AND file_hash = ? AND image_type = ?", app.ApplicationId, created.FileHash, 2).
		Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestApplicationImage_CreateOrGet_RequiresTx(t *testing.T) {
	img := ApplicationImage{}
	err := img.GetOrCreate(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), errTransactionRequired)
}
