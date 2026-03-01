package models

import (
	"database/sql"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestApplicationProfile_TableNames(t *testing.T) {
	require.Equal(t, tableApplicationProfileHistory, (ApplicationProfile{}).TableName())
	require.Equal(t, tableApplicationProfileScreenshot, (ApplicationProfileScreenshot{}).TableName())
}

func TestApplicationProfile_BeforeCreate_ValidatesOnlyStageAndDates(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{}, &ApplicationProfile{})
	app := createApplication(t, gdb, "app-profile-validation", "cust-1")
	img := createApplicationImage(t, gdb, app.ApplicationId, 1, 10)

	invalidStage := ApplicationProfile{
		ApplicationId: app.ApplicationId,
		Name:          sql.NullString{String: "main", Valid: true},
		Stage:         sql.NullInt16{Int16: profileStageReviewed, Valid: true},
	}
	err := gdb.Create(&invalidStage).Error
	require.Error(t, err)
	require.ErrorIs(t, err, errProfileStageInvalid)

	validWithImage := ApplicationProfile{
		ApplicationId:      app.ApplicationId,
		ApplicationImageId: img.ApplicationImageId,
		Name:               sql.NullString{String: "main", Valid: true},
	}
	require.NoError(t, gdb.Create(&validWithImage).Error)
}

func TestApplicationProfile_BeforeCreate_AlwaysStartsPending(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{}, &ApplicationProfile{})
	app := createApplication(t, gdb, "app-profile-pending", "cust-1")
	img := createApplicationImage(t, gdb, app.ApplicationId, 1, 11)

	invalidReview := ApplicationProfile{
		ApplicationId:      app.ApplicationId,
		ApplicationImageId: img.ApplicationImageId,
		Name:               sql.NullString{String: "main", Valid: true},
		ReviewAt:           nullTimeNow(),
	}
	err := gdb.Create(&invalidReview).Error
	require.Error(t, err)
	require.ErrorIs(t, err, errProfileReviewAtNotAllowed)

	invalidProduction := ApplicationProfile{
		ApplicationId:      app.ApplicationId,
		ApplicationImageId: img.ApplicationImageId,
		Name:               sql.NullString{String: "main", Valid: true},
		ProductionAt:       nullTimeNow(),
	}
	err = gdb.Create(&invalidProduction).Error
	require.Error(t, err)
	require.ErrorIs(t, err, errProfileProductionAtNotAllowed)

	valid := ApplicationProfile{
		ApplicationId:      app.ApplicationId,
		ApplicationImageId: img.ApplicationImageId,
		Name:               sql.NullString{String: "main", Valid: true},
	}
	require.NoError(t, gdb.Create(&valid).Error)

	var stored ApplicationProfile
	require.NoError(t, gdb.First(&stored, valid.ApplicationProfileId).Error)
	require.True(t, stored.Stage.Valid)
	require.EqualValues(t, profileStagePending, stored.Stage.Int16)
}

func TestApplicationProfile_BeforeCreate_ResolvesNestedImage(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{}, &ApplicationProfile{})
	app := createApplication(t, gdb, "app-profile-nested-image", "cust-1")

	profile := createApplicationProfileWithNestedImage(t, gdb, app.ApplicationId, "main", 14)
	require.NotZero(t, profile.ApplicationImageId)

	var imgCount int64
	require.NoError(t, gdb.Model(&ApplicationImage{}).
		Where("application_id = ?", app.ApplicationId).
		Count(&imgCount).Error)
	require.EqualValues(t, 1, imgCount)
}

func TestApplicationProfileScreenshot_BeforeCreate_ResolvesNestedImage(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{}, &ApplicationImage{}, &ApplicationProfile{}, &ApplicationProfileScreenshot{})
	app := createApplication(t, gdb, "app-profile-screenshot", "cust-1")
	profile := createApplicationProfile(t, gdb, app.ApplicationId, "main", 19)

	screenshot := ApplicationProfileScreenshot{
		ApplicationProfileId: profile.ApplicationProfileId,
		Position:             1,
		ApplicationImage: &ApplicationImage{
			ApplicationId:   app.ApplicationId,
			FileName:        sql.NullString{String: "screen.png", Valid: true},
			FileContentType: sql.NullString{String: "image/png", Valid: true},
			FileHash:        hash32(20),
			ImageType:       sql.NullInt16{Int16: 2, Valid: true},
		},
	}
	require.NoError(t, gdb.Create(&screenshot).Error)
	require.NotZero(t, screenshot.ApplicationImageId)
}
