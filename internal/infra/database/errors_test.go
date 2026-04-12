package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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

func TestMapDBError_ConstraintViolation_UniqueConstraint(t *testing.T) {
	cause := errors.New("UNIQUE constraint failed: applications.name")
	err := MapDbError(cause)
	assert.ErrorIs(t, err, ErrConstraintViolation)
	// original cause is still reachable via Unwrap
	assert.ErrorIs(t, err, cause)
}

func TestMapDBError_UnknownPassthrough(t *testing.T) {
	sentinel := errors.New("some driver error")
	assert.Equal(t, sentinel, MapDbError(sentinel))
}

func TestMapDBError_DBUnavailable(t *testing.T) {
	err := MapDbError(sql.ErrConnDone)
	assert.ErrorIs(t, err, ErrDBUnavailable)
}
