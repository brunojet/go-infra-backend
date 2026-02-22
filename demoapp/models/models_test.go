package models

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/database"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMigrations(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	db, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.Migrate(&Application{}, &TerminalModel{}, &TerminalModelConfiguration{})
	require.NoError(t, err)

	err = db.Migrate(&AuditEvent{}, &AuditFieldChange{})
	require.NoError(t, err)
}

func TestPreloadBidirectional(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, dbm)

	// Migrate only the tables used in this test
	err = dbm.Migrate(&FilterType{}, &Filter{})
	require.NoError(t, err)

	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	// Create a FilterType with two Filters (parent -> children)
	ft := FilterType{
		Name: sql.NullString{String: "colors", Valid: true},
		Filters: []Filter{
			{Name: sql.NullString{String: "red", Valid: true}},
			{Name: sql.NullString{String: "blue", Valid: true}},
		},
	}

	err = gdb.Create(&ft).Error
	require.NoError(t, err)
	require.NotZero(t, ft.FilterTypeId)

	// Parent -> children preload
	var loadedFt FilterType
	err = gdb.Preload("Filters").First(&loadedFt, ft.FilterTypeId).Error
	require.NoError(t, err)
	require.Len(t, loadedFt.Filters, 2)

	// Child -> parent preload
	var filters []Filter
	err = gdb.Preload("FilterType").Find(&filters, "filter_type_id = ?", ft.FilterTypeId).Error
	require.NoError(t, err)
	require.Len(t, filters, 2)
	for _, f := range filters {
		require.NotNil(t, f.FilterType)
		require.Equal(t, loadedFt.FilterTypeId, f.FilterType.FilterTypeId)
	}
}

func TestCatalogFilterSearch(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, dbm)

	// Migrate all involved models
	err = dbm.Migrate(
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationImage{}, &ApplicationProfile{},
		&ApplicationConfiguration{}, &ApplicationVersion{}, &ApplicationCatalog{},
		&FilterType{}, &Filter{},
	)
	require.NoError(t, err)

	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	// Create a terminal model + configuration used by both apps
	tm := TerminalModel{Name: sql.NullString{String: "tm", Valid: true}}
	require.NoError(t, gdb.Create(&tm).Error)

	tmc := TerminalModelConfiguration{TerminalModelId: tm.TerminalModelId, IntegrationType: sql.NullInt16{Int16: 1, Valid: true}}
	require.NoError(t, gdb.Create(&tmc).Error)

	// create two apps, each with its own filter type + filter and profile
	apps := []struct {
		name       string
		filterName string
		app        Application
		profile    ApplicationProfile
		catalog    ApplicationCatalog
	}{
		{name: "app-one", filterName: "red"},
		{name: "app-two", filterName: "blue"},
	}

	for i := range apps {
		// Application
		apps[i].app = Application{Name: sql.NullString{String: apps[i].name, Valid: true}}
		require.NoError(t, gdb.Create(&apps[i].app).Error)

		// Image for profile icon
		img := ApplicationImage{ApplicationId: apps[i].app.ApplicationId, FileName: sql.NullString{String: "icon.png", Valid: true}, FileContentType: sql.NullString{String: "image/png", Valid: true}, FileHash: make([]byte, 32), ImageType: sql.NullInt16{Int16: 1, Valid: true}}
		require.NoError(t, gdb.Create(&img).Error)

		// Profile
		apps[i].profile = ApplicationProfile{ApplicationId: apps[i].app.ApplicationId, ApplicationImageId: img.ApplicationImageId, Name: sql.NullString{String: fmt.Sprintf("profile-%d", i+1), Valid: true}}
		require.NoError(t, gdb.Create(&apps[i].profile).Error)

		// ApplicationConfiguration
		acfg := ApplicationConfiguration{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId, PackageName: sql.NullString{String: "pkg", Valid: true}}
		require.NoError(t, gdb.Create(&acfg).Error)

		// ApplicationVersion
		av := ApplicationVersion{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId}
		require.NoError(t, gdb.Create(&av).Error)

		// FilterType + Filter
		ft := FilterType{Name: sql.NullString{String: fmt.Sprintf("ft-%d", i+1), Valid: true}}
		require.NoError(t, gdb.Create(&ft).Error)
		f := Filter{FilterTypeId: ft.FilterTypeId, Name: sql.NullString{String: apps[i].filterName, Valid: true}}
		require.NoError(t, gdb.Create(&f).Error)

		// Associate filter to profile (many2many)
		require.NoError(t, gdb.Model(&apps[i].profile).Association("Filters").Append(&f))

		// Catalog entry
		cat := ApplicationCatalog{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId, Stage: 1, ApplicationVersionId: av.ApplicationVersionId, ApplicationProfileId: apps[i].profile.ApplicationProfileId}
		require.NoError(t, gdb.Create(&cat).Error)
		apps[i].catalog = cat
	}

	// Load catalogs with preloaded profile+filters and search by filter name in Go
	var catalogs []ApplicationCatalog
	require.NoError(t, gdb.Preload("ApplicationProfile").Preload("ApplicationProfile.Filters").Find(&catalogs).Error)

	// helper to find catalog by filter
	findByFilter := func(filterName string) *ApplicationCatalog {
		for i := range catalogs {
			ap := catalogs[i].ApplicationProfile
			if ap == nil {
				continue
			}
			for _, fl := range ap.Filters {
				if fl.Name.Valid && fl.Name.String == filterName {
					return &catalogs[i]
				}
			}
		}
		return nil
	}

	c1 := findByFilter("red")
	require.NotNil(t, c1)
	require.Equal(t, "app-one", apps[0].app.Name.String)

	c2 := findByFilter("blue")
	require.NotNil(t, c2)
	require.Equal(t, "app-two", apps[1].app.Name.String)
}

func TestPreloadOneShot(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, dbm)

	// Migrate involved models
	err = dbm.Migrate(
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationImage{}, &ApplicationProfile{},
		&ApplicationConfiguration{}, &ApplicationVersion{}, &ApplicationCatalog{},
		&FilterType{}, &Filter{},
	)
	require.NoError(t, err)

	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	acfg := ApplicationConfiguration{
		PackageName: sql.NullString{String: "pkg-one", Valid: true},
		Application: &Application{
			Name: sql.NullString{String: "one-app", Valid: true},
		},
		TerminalModelConfiguration: &TerminalModelConfiguration{
			TerminalModel:   &TerminalModel{Name: sql.NullString{String: "tm", Valid: true}},
			IntegrationType: sql.NullInt16{Int16: 1, Valid: true},
		},
	}

	require.NoError(t, gdb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&acfg).Error)
	require.NotZero(t, acfg.ApplicationId)
	require.NotNil(t, acfg.TerminalModelConfiguration)

	// Second one-shot: create profile with image and filters, attaching to the existing application
	prof := ApplicationProfile{
		ApplicationId: acfg.ApplicationId,
		Name:          sql.NullString{String: "main", Valid: true},
		ApplicationImage: &ApplicationImage{
			ApplicationId:   acfg.ApplicationId,
			FileName:        sql.NullString{String: "i.png", Valid: true},
			FileContentType: sql.NullString{String: "image/png", Valid: true},
			FileHash:        make([]byte, 32),
			ImageType:       sql.NullInt16{Int16: 1, Valid: true},
		},
		Filters: []Filter{
			{
				FilterType: &FilterType{Name: sql.NullString{String: "ft-one", Valid: true}},
				Name:       sql.NullString{String: "green", Valid: true},
			},
		},
	}

	av := ApplicationVersion{
		ApplicationId:                acfg.ApplicationId,
		TerminalModelConfigurationId: acfg.TerminalModelConfiguration.TerminalModelConfigurationId,
	}

	cat := ApplicationCatalog{
		ApplicationId:                acfg.ApplicationId,
		TerminalModelConfigurationId: acfg.TerminalModelConfiguration.TerminalModelConfigurationId,
		Stage:                        1,
		ApplicationVersion:           &av,
		ApplicationProfile:           &prof,
	}
	require.NoError(t, gdb.Create(&cat).Error)

	// One-shot chained preload: ApplicationProfile (and its Filters), ApplicationVersion, Application
	var catalogs []ApplicationCatalog
	err = gdb.
		Preload("ApplicationProfile.Filters").
		Preload("ApplicationProfile.ApplicationImage").
		Preload("ApplicationVersion").
		Preload("Application").Find(&catalogs).Error
	require.NoError(t, err)
	require.Len(t, catalogs, 1)

	got := catalogs[0]
	require.NotNil(t, got.ApplicationProfile)
	require.Len(t, got.ApplicationProfile.Filters, 1)
	require.NotNil(t, got.ApplicationVersion)
	require.NotNil(t, got.Application)
}
