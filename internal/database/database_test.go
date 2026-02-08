package database_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/database"
	dbadpt "github.com/brunojet/go-infra-backend/internal/database/adapters"
	"github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewSQLiteDatabase_DelegatesToAdapter(t *testing.T) {
	db, err := database.NewSQLiteDatabase("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, db.Shutdown(context.Background()))
}

func TestDatabaseManager_HealthCheck(t *testing.T) {
	db, err := database.NewSQLiteDatabase("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// HealthCheck should call Ping and return nil on a healthy in-memory DB
	assert.NoError(t, db.HealthCheck(context.Background()))
	assert.NoError(t, db.Shutdown(context.Background()))
}

func TestDatabaseManager_HealthCheckReturnsError(t *testing.T) {
	db, err := database.NewSQLiteDatabase("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// HealthCheck should call Ping and return nil on a healthy in-memory DB
	assert.NoError(t, db.Shutdown(context.Background()))
	assert.Error(t, db.HealthCheck(context.Background()))
}

func TestNewSQLiteDatabase_NonexistentDirReturnsError(t *testing.T) {
	_, err := database.NewSQLiteDatabase("nonexistent_dir/subdir.db")
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "unable to stat database directory"))
	}
}

func TestNewSQLiteDatabase_FilePathCreatesDB(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dbtest")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "mydb.db")
	db, err := database.NewSQLiteDatabase(dbPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// ensure file was created
	_, statErr := os.Stat(dbPath)
	assert.NoError(t, statErr)
	assert.NoError(t, db.Shutdown(context.Background()))
}

type fakeAdapter struct {
	closed bool
}

func (f *fakeAdapter) GormDB() (*gorm.DB, error) {
	return nil, errors.New("gorm error")
}

func (f *fakeAdapter) SqlDB() (*sql.DB, error) { return nil, errors.New("gorm error") }

func (f *fakeAdapter) Close() error {
	f.closed = true
	return nil
}

func TestNewDatabaseManager_ReturnsErrorWhenAdapterGormDBFails(t *testing.T) {
	fa := &fakeAdapter{}

	m, err := database.NewDatabaseManager(fa, &failingPlugin{})
	require.Error(t, err)
	require.Nil(t, m)
	require.True(t, fa.closed, "adapter.Close should be called on GormDB error")
}

// wrapperAdapter wraps a real adapter and records Close calls.
type wrapperAdapter struct {
	inner  contracts.DatabaseAdapter
	closed bool
}

func (w *wrapperAdapter) GormDB() (*gorm.DB, error) { return w.inner.GormDB() }
func (w *wrapperAdapter) Close() error              { w.closed = true; return w.inner.Close() }
func (w *wrapperAdapter) SqlDB() (*sql.DB, error)   { return w.inner.SqlDB() }

type failingPlugin struct{}

func (f *failingPlugin) Name() string                 { return "failing" }
func (f *failingPlugin) Initialize(db *gorm.DB) error { return errors.New("init fail") }

// goodPlugin used in tests to verify Initialize is called.
type goodPlugin struct {
	called *bool
}

func (g *goodPlugin) Name() string { return "good" }

func (g *goodPlugin) Initialize(db *gorm.DB) error {
	if g.called != nil {
		*g.called = true
	}
	return nil
}

func TestNewDatabaseManager_ClosesAdapterWhenPluginInitFails(t *testing.T) {
	// create a real sqlite adapter
	real, err := dbadpt.NewSQLite("memory")
	require.NoError(t, err)

	wa := &wrapperAdapter{inner: real}

	m, err := database.NewDatabaseManager(wa, &failingPlugin{})
	require.Error(t, err)
	require.Nil(t, m)
	require.True(t, wa.closed, "wrapper adapter Close should be called when plugin init fails")
}

func TestNewDatabaseManager_SucceedsWithNilAndValidPlugin(t *testing.T) {
	real, err := dbadpt.NewSQLite("memory")
	require.NoError(t, err)

	called := false
	gp := &goodPlugin{called: &called}

	m, err := database.NewDatabaseManager(real, nil, gp)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.True(t, called, "expected plugin Initialize to be called")
}

func TestDatabaseManager_HealthCheck_ReturnsAdapterGormDBError(t *testing.T) {
	fa := &fakeAdapter{}
	// create manager without plugins so NewDatabaseManager doesn't call GormDB
	m, err := database.NewDatabaseManager(fa)
	require.NoError(t, err)
	require.NotNil(t, m)

	// HealthCheck should propagate the adapter.GormDB error
	err = m.HealthCheck(context.Background())
	require.Error(t, err)
}

// adapter that returns an empty gorm.DB which should cause gormDb.DB() to fail
type emptyGormAdapter struct{}

func (e *emptyGormAdapter) GormDB() (*gorm.DB, error) { return &gorm.DB{}, nil }
func (e *emptyGormAdapter) SqlDB() (*sql.DB, error)   { return nil, errors.New("sql db error") }
func (e *emptyGormAdapter) Close() error              { return nil }

func TestDatabaseManager_HealthCheck_ReturnsWhenGormDBDBFails(t *testing.T) {
	ea := &emptyGormAdapter{}
	m, err := database.NewDatabaseManager(ea)
	require.NoError(t, err)
	require.NotNil(t, m)

	err = m.HealthCheck(context.Background())
	require.Error(t, err)
}
