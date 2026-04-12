package repositories

import (
	dberrs "github.com/brunojet/go-infra-backend/internal/infra/database/errors"
	"gorm.io/gorm"
)

// MapTxError inspects the GORM transaction result and maps it to a port sentinel.
// Lives here rather than in infra/database because it is specific to the GORM
// repository pattern: it interprets *gorm.DB fields (Error, RowsAffected) that
// are only meaningful in the context of repository operations.
func MapTxError(tx *gorm.DB) error {
	if tx == nil {
		return dberrs.ErrInvalidTx
	} else if tx.Error != nil {
		return dberrs.MapDbError(tx.Error)
	} else if tx.RowsAffected == 0 {
		return dberrs.ErrNotFound
	}
	return nil
}
