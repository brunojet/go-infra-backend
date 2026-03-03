package repositories

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBusinessRuleHelpers(t *testing.T) {
	sentinel := errors.New("business-cause")
	err := NewBusinessRuleError(sentinel)
	assert.Error(t, err)
	assert.True(t, IsBusinessRuleError(err))
	assert.ErrorIs(t, err, sentinel)

	assert.False(t, IsBusinessRuleError(errors.New("plain-error")))
}

func TestMapDBError_MapDBRecordNotFound(t *testing.T) {
	err := MapDbError(gorm.ErrRecordNotFound)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestMapDBError_SQLNoRows(t *testing.T) {
	err := MapDbError(sql.ErrNoRows)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestMapDBError_Nil(t *testing.T) {
	assert.Nil(t, MapDbError(nil))
}

func TestMapDBError_UnknownPassthrough(t *testing.T) {
	sentinel := errors.New("some driver error")
	assert.Equal(t, sentinel, MapDbError(sentinel))
}

func TestMapDBError_DBUnavailable(t *testing.T) {
	err := MapDbError(sql.ErrConnDone)
	assert.ErrorIs(t, err, ErrDBUnavailable)
}

func TestMapTxError_NilTx(t *testing.T) {
	err := MapTxError(nil)
	assert.ErrorIs(t, err, ErrInvalidTx)
}

func TestMapTxError_ErrorFieldMappings(t *testing.T) {
	// sql.ErrNoRows -> ErrNotFound
	tx := &gorm.DB{Error: sql.ErrNoRows}
	err := MapTxError(tx)
	assert.ErrorIs(t, err, ErrNotFound)

	// sql.ErrConnDone -> ErrDBUnavailable
	tx = &gorm.DB{Error: sql.ErrConnDone}
	err = MapTxError(tx)
	assert.ErrorIs(t, err, ErrDBUnavailable)

	// message contains "database is closed" -> ErrDBUnavailable
	tx = &gorm.DB{Error: errors.New("database is closed: foo")}
	err = MapTxError(tx)
	assert.ErrorIs(t, err, ErrDBUnavailable)

	// unknown error -> returned as-is
	tx = &gorm.DB{Error: errors.New("unknown")}
	err = MapTxError(tx)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "unknown")
	}
}

func TestMapTxError_RowsAffectedAndSuccess(t *testing.T) {
	// RowsAffected == 0 => ErrNotFound
	tx := &gorm.DB{Error: nil, RowsAffected: 0}
	err := MapTxError(tx)
	assert.ErrorIs(t, err, ErrNotFound)

	// success
	tx = &gorm.DB{Error: nil, RowsAffected: 2}
	err = MapTxError(tx)
	assert.NoError(t, err)
}
