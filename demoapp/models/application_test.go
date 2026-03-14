package models

import (
	"database/sql"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestApplicationCreate_UpsertSameNameSameCustomer(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	_ = createApplicationWithDescription(t, gdb, "app-upsert", "cust-1", "v1")
	a := Application{
		Name:        sql.NullString{String: "app-upsert", Valid: true},
		CustomerId:  sql.NullString{String: "cust-1", Valid: true},
		Description: sql.NullString{String: "v2", Valid: true},
	}
	RunInTransaction(t, gdb, func(tx *gorm.DB) (Application, error) {
		return a, tx.Create(&a).Error
	})

	var count int64
	require.NoError(t, gdb.Model(&Application{}).Where("name = ?", "app-upsert").Count(&count).Error)
	require.EqualValues(t, 1, count)

	var got Application
	require.NoError(t, gdb.Where("name = ?", "app-upsert").First(&got).Error)
	require.Equal(t, "v1", got.Description.String)
	require.Equal(t, "cust-1", got.CustomerId.String)
}

func TestApplicationCreate_RejectsSameNameDifferentCustomer(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	_ = createApplication(t, gdb, "app-conflict", "cust-1")
	second := Application{
		Name:       sql.NullString{String: "app-conflict", Valid: true},
		CustomerId: sql.NullString{String: "cust-2", Valid: true},
	}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Create(&second).Error, nil
	})
	require.ErrorIs(t, err, gorm.ErrCheckConstraintViolated)
}

func TestApplicationUpdate_RejectsNameCollisionWithOtherCustomer(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	app1 := createApplication(t, gdb, "app-a", "cust-a")
	_ = createApplication(t, gdb, "app-b", "cust-b")

	app1.Name = sql.NullString{String: "app-b", Valid: true}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&app1).Error, nil
	})
	require.ErrorIs(t, err, gorm.ErrCheckConstraintViolated)
}

func TestApplicationConfiguration_PackageCanRepeatInsideSameApp(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{}, &Application{}, &ApplicationConfiguration{})

	cfgTerm1 := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-same", 1)
	cfgTerm2 := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-same", 2)

	ac1 := createApplicationConfigurationOneShot(t, gdb, "app-cfg-one", "cust-1", cfgTerm1.TerminalModelConfigurationId, "pkg.shared")
	_ = createApplicationConfiguration(t, gdb, ac1.ApplicationId, cfgTerm2.TerminalModelConfigurationId, "pkg.shared")

	var count int64
	require.NoError(t, gdb.Model(&ApplicationConfiguration{}).Where("application_id = ? AND package_name = ?", ac1.ApplicationId, "pkg.shared").Count(&count).Error)
	require.EqualValues(t, 2, count)
}

func TestApplicationConfiguration_RejectsSamePackageAcrossDifferentApps(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{}, &Application{}, &ApplicationConfiguration{})

	cfgTerm := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-diff", 1)

	app2 := createApplication(t, gdb, "app-cfg-b", "cust-b")

	_ = createApplicationConfigurationOneShot(t, gdb, "app-cfg-a", "cust-a", cfgTerm.TerminalModelConfigurationId, "pkg.global")

	ac2 := ApplicationConfiguration{ApplicationId: app2.ApplicationId, TerminalModelConfigurationId: cfgTerm.TerminalModelConfigurationId, PackageName: sql.NullString{String: "pkg.global", Valid: true}}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Create(&ac2).Error, nil
	})
	require.ErrorIs(t, err, gorm.ErrCheckConstraintViolated)
}

func TestApplicationCreate_InvalidRequiredFields(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	invalid := Application{
		Name:       sql.NullString{},
		CustomerId: sql.NullString{String: "cust-x", Valid: true},
	}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Create(&invalid).Error, nil
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), errNameAndCustomerRequired)
}

func TestApplicationUpdate_SuccessSameOwnerNoNameCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	app := createApplicationWithDescription(t, gdb, "app-update-ok", "cust-ok", "before")

	app.Description = sql.NullString{String: "after", Valid: true}
	require.NoError(t, RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&app).Error, nil
	}))

	var got Application
	require.NoError(t, gdb.First(&got, app.ApplicationId).Error)
	require.Equal(t, "after", got.Description.String)
}

func TestApplicationUpdate_InvalidRequiredFields(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	app := createApplication(t, gdb, "app-invalid-update", "cust-a")

	       app.CustomerId = sql.NullString{}
	       err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		       return tx.Save(&app).Error, nil
	       })
	       require.Error(t, err)
	       require.Contains(t, err.Error(), errAppIDAndCustomerRequired)
}

func TestApplicationConfigurationUpdate_SuccessAndErrorBranches(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{}, &Application{}, &ApplicationConfiguration{})

	cfgTerm1 := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-update", 1)
	cfgTerm2 := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-update", 2)

	ac1 := createApplicationConfigurationOneShot(t, gdb, "app-upd-cfg-a", "cust-a", cfgTerm1.TerminalModelConfigurationId, "pkg-a")
	ac2 := createApplicationConfigurationOneShot(t, gdb, "app-upd-cfg-b", "cust-b", cfgTerm2.TerminalModelConfigurationId, "pkg-b")

	// Success branch of BeforeUpdate
	ac1.PackageName = sql.NullString{String: "pkg-a-2", Valid: true}
	require.NoError(t, RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&ac1).Error, nil
	}))

	       // Error branch: invalid package_name
	       ac1.PackageName = sql.NullString{}
	       err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		       return tx.Save(&ac1).Error, nil
	       })
	       require.Error(t, err)
	       require.Contains(t, err.Error(), errPackageNameRequired)

	       // Error branch: package collides with another app
	       ac2.PackageName = sql.NullString{String: "pkg-a-2", Valid: true}
	       err = RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		       return tx.Save(&ac2).Error, nil
	       })
	       require.ErrorIs(t, err, gorm.ErrCheckConstraintViolated)
}
