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

	// Create a FilterType with two Filters
	ft := createFilterType(t, gdb, "colors")
	_ = createFilter(t, gdb, ft.FilterTypeId, "red")
	_ = createFilter(t, gdb, ft.FilterTypeId, "blue")

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
	suffix := t.Name()

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
		apps[i].app = Application{
			Name:       sql.NullString{String: fmt.Sprintf("%s-%s", apps[i].name, suffix), Valid: true},
			CustomerId: sql.NullString{String: "customer-1", Valid: true},
		}
		require.NoError(t, gdb.Create(&apps[i].app).Error)

		// Image for profile icon
		img := ApplicationImage{ApplicationId: apps[i].app.ApplicationId, FileName: sql.NullString{String: "icon.png", Valid: true}, FileContentType: sql.NullString{String: "image/png", Valid: true}, FileHash: make([]byte, 32), ImageType: sql.NullInt16{Int16: 1, Valid: true}}
		require.NoError(t, gdb.Create(&img).Error)

		// Profile
		apps[i].profile = ApplicationProfile{ApplicationId: apps[i].app.ApplicationId, ApplicationImageId: img.ApplicationImageId, Name: sql.NullString{String: fmt.Sprintf("profile-%d", i+1), Valid: true}, Stage: sql.NullInt16{Int16: 0, Valid: true}}
		require.NoError(t, gdb.Create(&apps[i].profile).Error)

		// ApplicationConfiguration
		acfg := ApplicationConfiguration{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId, PackageName: sql.NullString{String: fmt.Sprintf("pkg-%d-%s", i+1, suffix), Valid: true}}
		require.NoError(t, gdb.Create(&acfg).Error)

		apps[i].profile.Stage = sql.NullInt16{Int16: 1, Valid: true}
		require.NoError(t, gdb.Save(&apps[i].profile).Error)

		// ApplicationVersion
		av := ApplicationVersion{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId}
		require.NoError(t, gdb.Create(&av).Error)

		// FilterType + Filter
		ft := createFilterType(t, gdb, fmt.Sprintf("ft-%d", i+1))
		f := createFilter(t, gdb, ft.FilterTypeId, apps[i].filterName)

		// Associate filter to profile (many2many)
		require.NoError(t, gdb.Model(&apps[i].profile).Association("Filters").Append(&f))

		// Catalog entry
		cat := ApplicationCatalog{ApplicationId: apps[i].app.ApplicationId, TerminalModelConfigurationId: tmc.TerminalModelConfigurationId, Stage: 1, ApplicationVersionId: &av.ApplicationVersionId, ApplicationProfileId: apps[i].profile.ApplicationProfileId}
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
	require.Equal(t, fmt.Sprintf("%s-%s", "app-one", suffix), apps[0].app.Name.String)

	c2 := findByFilter("blue")
	require.NotNil(t, c2)
	require.Equal(t, fmt.Sprintf("%s-%s", "app-two", suffix), apps[1].app.Name.String)
}

func createCatalogWithAssociations(t *testing.T, gdb *gorm.DB, appName, packageName, terminalModelName, filterType, filterName string) *ApplicationCatalog {
	acfg := ApplicationConfiguration{
		PackageName: sql.NullString{String: packageName, Valid: true},
		Application: &Application{
			Name:       sql.NullString{String: appName, Valid: true},
			CustomerId: sql.NullString{String: "123456", Valid: true},
		},
		TerminalModelConfiguration: &TerminalModelConfiguration{
			TerminalModel:   &TerminalModel{Name: sql.NullString{String: terminalModelName, Valid: true}},
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
		Stage:         sql.NullInt16{Int16: 0, Valid: true},
		ApplicationImage: &ApplicationImage{
			ApplicationId:   acfg.ApplicationId,
			FileName:        sql.NullString{String: "i.png", Valid: true},
			FileContentType: sql.NullString{String: "image/png", Valid: true},
			FileHash:        make([]byte, 32),
			ImageType:       sql.NullInt16{Int16: 1, Valid: true},
		},
		Filters: []Filter{
			{
				FilterType: &FilterType{Name: sql.NullString{String: filterType, Valid: true}},
				Name:       sql.NullString{String: filterName, Valid: true},
			},
		},
	}
	require.NoError(t, gdb.Create(&prof).Error)

	av := ApplicationVersion{
		ApplicationId:                acfg.ApplicationId,
		TerminalModelConfigurationId: acfg.TerminalModelConfiguration.TerminalModelConfigurationId,
	}
	require.NoError(t, gdb.Create(&av).Error)

	cat := ApplicationCatalog{
		ApplicationId:                acfg.ApplicationId,
		TerminalModelConfigurationId: acfg.TerminalModelConfiguration.TerminalModelConfigurationId,
		Stage:                        1,
		ApplicationVersionId:         &av.ApplicationVersionId,
		ApplicationProfileId:         prof.ApplicationProfileId,
	}
	require.NoError(t, gdb.Create(&cat).Error)

	// One-shot chained preload: load only the catalog we just created so subsequent
	// calls don't accidentally return the first row in the table.
	var got ApplicationCatalog
	err := gdb.
		Preload("ApplicationProfile.Filters").
		Preload("ApplicationProfile.ApplicationImage").
		Preload("ApplicationVersion").
		Preload("Application").
		Where("application_id = ? AND terminal_model_configuration_id = ? AND stage = ?", acfg.ApplicationId, acfg.TerminalModelConfiguration.TerminalModelConfigurationId, 1).
		First(&got).Error
	require.NoError(t, err)
	require.NotNil(t, got.ApplicationProfile)
	require.Len(t, got.ApplicationProfile.Filters, 1)
	require.NotNil(t, got.ApplicationVersion)
	require.NotNil(t, got.Application)
	return &got
}

func TestPreloadOneShot(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, dbm)

	err = dbm.Migrate(
		&TerminalModel{}, &TerminalModelConfiguration{},
		&Application{}, &ApplicationImage{}, &ApplicationProfile{},
		&ApplicationConfiguration{}, &ApplicationVersion{}, &ApplicationCatalog{},
		&FilterType{}, &Filter{},
	)
	require.NoError(t, err)

	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	cat1 := createCatalogWithAssociations(t, gdb, "app-oneshot-one", "pkg-oneshot-one", "tm-oneshot-1", "ft-oneshot-one", "green")
	require.NotNil(t, cat1)

	cat2 := createCatalogWithAssociations(t, gdb, "app-oneshot-two", "pkg-oneshot-two", "tm-oneshot-2", "ft-oneshot-two", "blue")
	require.NotNil(t, cat2)
}

func TestApplicationBeforeUpdate_RejectsDifferentCustomerOwnership(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)
	require.NotNil(t, dbm)

	err = dbm.Migrate(&Application{})
	require.NoError(t, err)

	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	app := Application{
		Name:       sql.NullString{String: "app-update-owner", Valid: true},
		CustomerId: sql.NullString{String: "partner-a", Valid: true},
	}
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&app).Error
	}))

	app.CustomerId = sql.NullString{String: "partner-b", Valid: true}
	err = gdb.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&app).Error
	})
	require.ErrorIs(t, err, gorm.ErrCheckConstraintViolated)
}

func TestTerminalModel_UniqueNameAndUpdateValidation(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)

	require.NoError(t, dbm.Migrate(&TerminalModel{}))
	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	tm1 := TerminalModel{Name: sql.NullString{String: "tm-unique", Valid: true}}
	require.NoError(t, gdb.Create(&tm1).Error)

	tm2 := TerminalModel{
		Name:        sql.NullString{String: "tm-unique", Valid: true},
		Description: sql.NullString{String: "updated-by-upsert", Valid: true},
	}
	require.NoError(t, gdb.Create(&tm2).Error)

	var got TerminalModel
	require.NoError(t, gdb.Where("name = ?", "tm-unique").First(&got).Error)
	require.Equal(t, tm1.Description.String, got.Description.String)

	tm3 := TerminalModel{Name: sql.NullString{String: "tm-3", Valid: true}}
	require.NoError(t, gdb.Create(&tm3).Error)
	tm3.Name = sql.NullString{}
	err = gdb.Save(&tm3).Error
	require.Error(t, err)
}

func TestTerminalModelConfiguration_UniquePairAndUpdateValidation(t *testing.T) {
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeMemory))

	dbm, err := database.NewDatabaseManagerFromEnv()
	require.NoError(t, err)

	require.NoError(t, dbm.Migrate(&TerminalModel{}, &TerminalModelConfiguration{}))
	gdb, err := dbm.DatabaseAdapter().GormDB()
	require.NoError(t, err)

	tm := TerminalModel{Name: sql.NullString{String: "tm-config", Valid: true}}
	require.NoError(t, gdb.Create(&tm).Error)

	cfg1 := TerminalModelConfiguration{TerminalModelId: tm.TerminalModelId, IntegrationType: sql.NullInt16{Int16: 1, Valid: true}}
	require.NoError(t, gdb.Create(&cfg1).Error)

	cfg2 := TerminalModelConfiguration{TerminalModelId: tm.TerminalModelId, IntegrationType: sql.NullInt16{Int16: 1, Valid: true}}
	require.NoError(t, gdb.Create(&cfg2).Error)

	cfg3 := TerminalModelConfiguration{TerminalModelId: tm.TerminalModelId, IntegrationType: sql.NullInt16{Int16: 2, Valid: true}}
	require.NoError(t, gdb.Create(&cfg3).Error)
	cfg3.IntegrationType = sql.NullInt16{}
	err = gdb.Save(&cfg3).Error
	require.Error(t, err)
}
