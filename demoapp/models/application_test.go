package models

import (
	"database/sql"
	"testing"

	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
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

// TestApplicationCreate_ConflictSameNameDifferentCustomer_IsIgnored verifies the model-level
// behavior: BeforeCreate uses INSERT OR IGNORE, so a name collision with a different customer
// results in a silent no-op at the GORM level (err=nil, struct not populated).
// Conflict rejection is enforced at the service layer via ErrConflictValidationFailed.
func TestApplicationCreate_ConflictSameNameDifferentCustomer_IsIgnored(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	_ = createApplication(t, gdb, "app-conflict", "cust-1")
	second := Application{
		Name:       sql.NullString{String: "app-conflict", Valid: true},
		CustomerId: sql.NullString{String: "cust-2", Valid: true},
	}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Create(&second).Error, nil
	})
	require.NoError(t, err)
	require.Zero(t, second.ApplicationId) // insert silently ignored; struct not populated

	var found Application
	require.NoError(t, gdb.Where("name = ?", "app-conflict").First(&found).Error)
	require.Equal(t, "cust-1", found.CustomerId.String) // original record preserved
}

func TestApplicationUpdate_RejectsNameCollisionWithOtherCustomer(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &Application{})

	app1 := createApplication(t, gdb, "app-a", "cust-a")
	_ = createApplication(t, gdb, "app-b", "cust-b")

	app1.Name = sql.NullString{String: "app-b", Valid: true}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&app1).Error, nil
	})
	require.ErrorIs(t, portsrepos.MapDbError(err), portsrepos.ErrConstraintViolation)
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

// TestApplicationConfiguration_ConflictSamePackageSameCfgTerm_IsIgnored verifies the
// model-level behavior: BeforeCreate uses INSERT OR IGNORE, so a (terminal_model_configuration_id,
// package_name) collision across different apps results in a silent no-op at the GORM level.
// Conflict rejection is enforced at the service layer via ErrConflictValidationFailed.
func TestApplicationConfiguration_ConflictSamePackageSameCfgTerm_IsIgnored(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{}, &Application{}, &ApplicationConfiguration{})

	cfgTerm := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-diff", 1)

	app2 := createApplication(t, gdb, "app-cfg-b", "cust-b")

	_ = createApplicationConfigurationOneShot(t, gdb, "app-cfg-a", "cust-a", cfgTerm.TerminalModelConfigurationId, "pkg.global")

	ac2 := ApplicationConfiguration{ApplicationId: app2.ApplicationId, TerminalModelConfigurationId: cfgTerm.TerminalModelConfigurationId, PackageName: sql.NullString{String: "pkg.global", Valid: true}}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Create(&ac2).Error, nil
	})
	require.NoError(t, err)

	var count int64
	require.NoError(t, gdb.Model(&ApplicationConfiguration{}).Where("terminal_model_configuration_id = ? AND package_name = ?", cfgTerm.TerminalModelConfigurationId, "pkg.global").Count(&count).Error)
	require.EqualValues(t, 1, count) // original record preserved
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

func TestApplicationConfigurationUpdate_SuccessAndErrorBranches(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{}, &Application{}, &ApplicationConfiguration{})

	cfgTerm := createTerminalModelConfigurationOneShot(t, gdb, "tm-appcfg-update", 1)

	ac1 := createApplicationConfigurationOneShot(t, gdb, "app-upd-cfg-a", "cust-a", cfgTerm.TerminalModelConfigurationId, "pkg-a")
	ac2 := createApplicationConfigurationOneShot(t, gdb, "app-upd-cfg-b", "cust-b", cfgTerm.TerminalModelConfigurationId, "pkg-b")

	// Success branch: rename to a non-conflicting package name
	ac1.PackageName = sql.NullString{String: "pkg-a-2", Valid: true}
	require.NoError(t, RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&ac1).Error, nil
	}))

	// Error branch: rename ac2 to a package that already exists for the same cfgTerm
	ac2.PackageName = sql.NullString{String: "pkg-a-2", Valid: true}
	err := RunInTransaction(t, gdb, func(tx *gorm.DB) (error, error) {
		return tx.Save(&ac2).Error, nil
	})
	require.ErrorIs(t, portsrepos.MapDbError(err), portsrepos.ErrConstraintViolation)
}
