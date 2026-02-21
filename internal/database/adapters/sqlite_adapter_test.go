package adapters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSQLiteInMemory_ReturnsDatabase(t *testing.T) {
	db, err := NewSQLite("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// ensure it satisfies the contract
	var _ dbcontracts.DatabaseAdapter = db

	// ensure GormDB is available
	gdb, err := db.GormDB()
	assert.NoError(t, err)
	assert.NotNil(t, gdb)

	// Close should succeed
	assert.NoError(t, db.Close())
}

func TestBuildDSNFromPath_MemoryAndEmpty(t *testing.T) {
	dsn, err := buildDSNFromPath("")
	assert.NoError(t, err)
	assert.Equal(t, "file::memory:?mode=memory&cache=shared&_pragma=foreign_keys(1)", dsn)

	dsn, err = buildDSNFromPath("use-memory")
	assert.NoError(t, err)
	assert.Equal(t, "file::memory:?mode=memory&cache=shared&_pragma=foreign_keys(1)", dsn)
}

func TestBuildDSNFromPath_FilePathSuccess(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "app.db")

	dsn, err := buildDSNFromPath(dbPath)
	assert.NoError(t, err)
	abs, _ := filepath.Abs(dbPath)
	want := "file:" + filepath.ToSlash(abs) + "?_pragma=foreign_keys(1)"
	assert.Equal(t, want, dsn)
}

func TestBuildDSNFromPath_FilePathNonexistentDir(t *testing.T) {
	tmp := t.TempDir()
	// create a nested path that does not exist
	dbPath := filepath.Join(tmp, "nope", "will.db")
	_, err := buildDSNFromPath(dbPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to stat database directory")
}

func TestBuildDSNFromPath_InvalidInputs(t *testing.T) {
	// plain string that doesn't match memory or .db
	_, err := buildDSNFromPath("just-a-string")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unable to configure dsn"))

	// contains file: but still should be rejected by this helper
	dsn, err := buildDSNFromPath("file:my.db")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(dsn, "file:my.db"))
}

func TestNewSQLite_ContinuesWithNilPlugin(t *testing.T) {
	// sanity-check NewSQLite handles nil plugins by continuing
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "nodb.db")
	// ensure directory exists
	assert.NoError(t, os.MkdirAll(filepath.Dir(dbPath), 0o755))

	db, err := NewSQLite(dbPath)
	// NewSQLite will attempt to open the file; expect no error
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, db.Close())
}

func TestSQLiteAdapter_GormSqlCloseLifecycle(t *testing.T) {
	a, err := NewSQLite("memory")
	require.NoError(t, err)
	require.NotNil(t, a)

	// GormDB should be available
	gdb, err := a.GormDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)

	// SqlDB should be available and Pingable
	sdb, err := a.SqlDB()
	require.NoError(t, err)
	require.NotNil(t, sdb)
	require.NoError(t, sdb.Ping())

	// Close should succeed
	require.NoError(t, a.Close())

	// After Close both GormDB and SqlDB should return errors
	_, err = a.GormDB()
	require.Error(t, err)
	_, err = a.SqlDB()
	require.Error(t, err)

	// Close again is idempotent
	require.NoError(t, a.Close())
}
