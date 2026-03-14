package models

import (
	"database/sql"
	"testing"
	"time"

	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTerminalModelConfigurationOneShot(t *testing.T, gdb *gorm.DB, terminalModelName string, integrationType int16) TerminalModelConfiguration {
	t.Helper()

	var cfg TerminalModelConfiguration

	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		tm := buildTerminalModel(terminalModelName, "")
		if err := tx.Create(&tm).Error; err != nil {
			return err
		}
		if tm.TerminalModelId == 0 {
			if err := tx.Where("name = ?", terminalModelName).First(&tm).Error; err != nil {
				return err
			}
		}

		cfg = buildTerminalModelConfiguration(tm.TerminalModelId, integrationType)
		cfg.TerminalModel = &tm
		return tx.Create(&cfg).Error
	}))

	require.NotZero(t, cfg.TerminalModelConfigurationId)
	require.NotZero(t, cfg.TerminalModelId)
	require.NotNil(t, cfg.TerminalModel)
	require.NotZero(t, cfg.TerminalModel.TerminalModelId)

	return cfg
}

func createTerminalModel(t *testing.T, gdb *gorm.DB, name string) TerminalModel {
	return createTerminalModelWithDescription(t, gdb, name, "")
}

func buildTerminalModel(name, description string) TerminalModel {
	tm := TerminalModel{Name: sql.NullString{String: name, Valid: true}}
	if description != "" {
		tm.Description = sql.NullString{String: description, Valid: true}
	}
	return tm
}

func createTerminalModelWithDescription(t *testing.T, gdb *gorm.DB, name, description string) TerminalModel {
	t.Helper()

	tm := buildTerminalModel(name, description)
	require.NoError(t, gdb.Create(&tm).Error)
	require.NotZero(t, tm.TerminalModelId)
	return tm
}

func buildTerminalModelConfiguration(terminalModelID int64, integrationType int16) TerminalModelConfiguration {
	return TerminalModelConfiguration{
		TerminalModelId: terminalModelID,
		IntegrationType: sql.NullInt16{Int16: integrationType, Valid: true},
	}
}

func createTerminalModelConfiguration(t *testing.T, gdb *gorm.DB, terminalModelID int64, integrationType int16) TerminalModelConfiguration {
	t.Helper()

	cfg := buildTerminalModelConfiguration(terminalModelID, integrationType)
	require.NoError(t, gdb.Create(&cfg).Error)
	require.NotZero(t, cfg.TerminalModelConfigurationId)
	return cfg
}

func createApplication(t *testing.T, gdb *gorm.DB, name, customer string) Application {
	return createApplicationWithDescription(t, gdb, name, customer, "")
}

func createApplicationWithDescription(t *testing.T, gdb *gorm.DB, name, customer, description string) Application {
	t.Helper()

	var app Application
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		app = Application{
			Name:       sql.NullString{String: name, Valid: true},
			CustomerId: sql.NullString{String: customer, Valid: true},
		}
		if description != "" {
			app.Description = sql.NullString{String: description, Valid: true}
		}
		return tx.Create(&app).Error
	}))
	require.NotZero(t, app.ApplicationId)
	return app
}

func createApplicationConfiguration(t *testing.T, gdb *gorm.DB, applicationID, terminalModelConfigurationID int64, packageName string) ApplicationConfiguration {
	t.Helper()

	var ac ApplicationConfiguration
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		ac = ApplicationConfiguration{
			ApplicationId:                applicationID,
			TerminalModelConfigurationId: terminalModelConfigurationID,
			PackageName:                  sql.NullString{String: packageName, Valid: true},
		}
		return tx.Create(&ac).Error
	}))
	return ac
}

func createApplicationConfigurationOneShot(t *testing.T, gdb *gorm.DB, appName, customer string, terminalModelConfigurationID int64, packageName string) ApplicationConfiguration {
	t.Helper()

	var ac ApplicationConfiguration
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		app := Application{
			Name:       sql.NullString{String: appName, Valid: true},
			CustomerId: sql.NullString{String: customer, Valid: true},
		}
		if err := tx.Create(&app).Error; err != nil {
			return err
		}

		ac = ApplicationConfiguration{
			ApplicationId:                app.ApplicationId,
			TerminalModelConfigurationId: terminalModelConfigurationID,
			PackageName:                  sql.NullString{String: packageName, Valid: true},
		}
		return tx.Create(&ac).Error
	}))

	require.NotZero(t, ac.ApplicationId)
	require.NotZero(t, ac.TerminalModelConfigurationId)
	return ac
}

func createApplicationImage(t *testing.T, gdb *gorm.DB, applicationID int64, imageType int16, hashByte byte) ApplicationImage {
	t.Helper()

	var img ApplicationImage
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		tmp := newApplicationImage(applicationID, imageType, hashByte)
		if err := tx.Create(&tmp).Error; err != nil {
			return err
		}
		img = tmp
		return nil
	}))
	require.NotZero(t, img.ApplicationImageId)
	return img
}

func newApplicationImage(applicationID int64, imageType int16, hashByte byte) ApplicationImage {
	return ApplicationImage{
		ApplicationId:   applicationID,
		FileName:        sql.NullString{String: "icon.png", Valid: true},
		FileContentType: sql.NullString{String: "image/png", Valid: true},
		FileHash:        hash32(hashByte),
		ImageType:       sql.NullInt16{Int16: imageType, Valid: true},
	}
}

func hash32(fill byte) []byte {
	hash := make([]byte, 32)
	for i := range hash {
		hash[i] = fill
	}
	return hash
}

func createApplicationProfile(t *testing.T, gdb *gorm.DB, applicationID int64, name string, imageHashByte byte) ApplicationProfile {
	t.Helper()

	var profile ApplicationProfile
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		img := newApplicationImage(applicationID, 1, imageHashByte)
		if err := tx.Create(&img).Error; err != nil {
			return err
		}
		tmp := ApplicationProfile{
			ApplicationId:      applicationID,
			ApplicationImageId: img.ApplicationImageId,
			Name:               sql.NullString{String: name, Valid: true},
		}
		if err := tx.Create(&tmp).Error; err != nil {
			return err
		}
		profile = tmp
		return nil
	}))
	require.NotZero(t, profile.ApplicationProfileId)
	return profile
}

func createApplicationProfileWithNestedImage(t *testing.T, gdb *gorm.DB, applicationID int64, name string, imageHashByte byte) ApplicationProfile {
	t.Helper()

	var profile ApplicationProfile
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		tmp := ApplicationProfile{
			ApplicationId: applicationID,
			Name:          sql.NullString{String: name, Valid: true},
			ApplicationImage: &ApplicationImage{
				ApplicationId:   applicationID,
				FileName:        sql.NullString{String: "icon.png", Valid: true},
				FileContentType: sql.NullString{String: "image/png", Valid: true},
				FileHash:        hash32(imageHashByte),
				ImageType:       sql.NullInt16{Int16: 1, Valid: true},
			},
		}
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&tmp).Error; err != nil {
			return err
		}
		profile = tmp
		return nil
	}))
	require.NotZero(t, profile.ApplicationProfileId)
	return profile
}

func nullTimeNow() sql.NullTime {
	return sql.NullTime{Time: time.Now().UTC(), Valid: true}
}

func createFilterType(t *testing.T, gdb *gorm.DB, name string) FilterType {
	return createFilterTypeWithDescription(t, gdb, name, "")
}

func createFilterTypeWithDescription(t *testing.T, gdb *gorm.DB, name, description string) FilterType {
	t.Helper()

	ft := FilterType{Name: sql.NullString{String: name, Valid: true}}
	if description != "" {
		ft.Description = sql.NullString{String: description, Valid: true}
	}
	require.NoError(t, gdb.Create(&ft).Error)
	require.NotZero(t, ft.FilterTypeId)
	return ft
}

func createFilter(t *testing.T, gdb *gorm.DB, filterTypeID int64, name string) Filter {
	return createFilterWithDescription(t, gdb, filterTypeID, name, "")
}

func createFilterWithDescription(t *testing.T, gdb *gorm.DB, filterTypeID int64, name, description string) Filter {
	t.Helper()

	f := Filter{
		FilterTypeId: filterTypeID,
		Name:         sql.NullString{String: name, Valid: true},
	}
	if description != "" {
		f.Description = sql.NullString{String: description, Valid: true}
	}
	require.NoError(t, gdb.Create(&f).Error)
	require.NotZero(t, f.FilterId)
	return f
}

func createFilterOneShot(t *testing.T, gdb *gorm.DB, filterTypeName, filterName string) Filter {
	t.Helper()

	var f Filter
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		ft := FilterType{Name: sql.NullString{String: filterTypeName, Valid: true}}
		if err := tx.Create(&ft).Error; err != nil {
			return err
		}

		f = Filter{
			FilterTypeId: ft.FilterTypeId,
			Name:         sql.NullString{String: filterName, Valid: true},
		}
		return tx.Create(&f).Error
	}))

	require.NotZero(t, f.FilterTypeId)
	require.NotZero(t, f.FilterId)
	return f
}

func createFilterOneShotWithNestedFilterType(t *testing.T, gdb *gorm.DB, filterTypeName, filterName string) Filter {
	t.Helper()

	f := Filter{
		Name: sql.NullString{String: filterName, Valid: true},
		FilterType: &FilterType{
			Name: sql.NullString{String: filterTypeName, Valid: true},
		},
	}

	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&f).Error
	}))

	require.NotZero(t, f.FilterId)
	require.NotZero(t, f.FilterTypeId)
	require.NotNil(t, f.FilterType)
	require.NotZero(t, f.FilterType.FilterTypeId)
	return f
}

// RunInTransaction is a generic helper to run a DB operation in a transaction and return the created/modified entity.
func RunInTransaction[T any](t *testing.T, gdb *gorm.DB, op func(tx *gorm.DB) (T, error)) T {
	t.Helper()
	var result T
	require.NoError(t, gdb.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(portsrepos.ContextWithTx(tx.Statement.Context, tx))
		var err error
		result, err = op(tx)
		return err
	}))
	return result
}
