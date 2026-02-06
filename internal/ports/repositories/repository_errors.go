package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrDBUnavailable = sql.ErrConnDone
	ErrInvalidTx     = errors.New("invalid transaction")
	ErrNotFound      = gorm.ErrRecordNotFound
)

func MapDbError(err error) error {
	if err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if errors.Is(err, sql.ErrConnDone) || strings.Contains(err.Error(), "database is closed") {
		return ErrDBUnavailable
	}
	return err
}

func MapTxError(tx *gorm.DB) error {
	if tx == nil {
		return ErrInvalidTx
	} else if tx.Error != nil {
		return MapDbError(tx.Error)
	} else if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
