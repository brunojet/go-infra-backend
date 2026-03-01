package dbtest

import (
	"testing"

	internaldbtest "github.com/brunojet/go-infra-backend/internal/testutil/dbtest"
	"gorm.io/gorm"
)

func OpenMemoryDB(t *testing.T, migrate ...any) *gorm.DB {
	return internaldbtest.OpenMemoryDB(t, migrate...)
}

func RunWithSQLiteRetry(fn func() error) error {
	return internaldbtest.RunWithSQLiteRetry(fn)
}
