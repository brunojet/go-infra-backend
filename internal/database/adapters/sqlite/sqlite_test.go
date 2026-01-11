package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnector_OpenAndClose_InMemory(t *testing.T) {
	c := New("") // deve usar file::memory:?cache=shared
	db, err := c.Open()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, c.DB())
	err = c.Close()
	assert.NoError(t, err)
}

func TestConnector_OpenAndClose_File(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.db")
	c := New(filePath)
	db, err := c.Open()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, c.DB())
	// O arquivo deve ser criado
	_, statErr := os.Stat(filePath)
	assert.NoError(t, statErr)
	err = c.Close()
	assert.NoError(t, err)
}

func TestConnector_Close_Idempotent(t *testing.T) {
	c := New("")
	_, err := c.Open()
	assert.NoError(t, err)
	assert.NoError(t, c.Close())
	// Close novamente não deve dar erro
	assert.NoError(t, c.Close())
}

func TestNormalizeDSN(t *testing.T) {
	dsn, path, ok := normalizeDSN("")
	assert.Equal(t, defaultInMemoryDSN, dsn)
	assert.False(t, ok)

	dsn, path, ok = normalizeDSN(":memory:")
	assert.Equal(t, defaultInMemoryDSN, dsn)
	assert.False(t, ok)

	dsn, path, ok = normalizeDSN("file:test.db?mode=memory")
	assert.Equal(t, "file:test.db?mode=memory", dsn)
	assert.False(t, ok)

	dsn, path, ok = normalizeDSN("/tmp/test.db?foo=bar")
	assert.Equal(t, "/tmp/test.db?foo=bar", dsn)
	assert.Equal(t, "/tmp/test.db", path)
	assert.True(t, ok)

	dsn, path, ok = normalizeDSN("/tmp/test.db")
	assert.Equal(t, "/tmp/test.db", dsn)
	assert.Equal(t, "/tmp/test.db", path)
	assert.True(t, ok)
}
