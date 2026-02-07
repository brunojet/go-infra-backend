package contracts

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// DatabaseAdapter represents a lifecycle-managed database instance.
// Implementations should provide access to the underlying GORM DB, the
// underlying SQL DB and a safe Close method to release resources.
type DatabaseAdapter interface {
	// GormDB returns the underlying *gorm.DB used for ORM operations.
	GormDB() (*gorm.DB, error)
	// SqlDB returns the underlying *sql.DB used for lower-level operations.
	SqlDB() (*sql.DB, error)
	Close() error
}

type DatabaseManager interface {
	DatabaseAdapter() DatabaseAdapter
	Shutdown(ctx context.Context) error
	HealthCheck(ctx context.Context) error
}
