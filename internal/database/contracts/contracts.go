package contracts

import (
	"gorm.io/gorm"
)

// Database represents a lifecycle-managed database instance.
// Implementations should provide access to the underlying GORM DB and
// a safe Close method to release resources.
type Database interface {
	// GormDB returns the underlying *gorm.DB used for ORM operations.
	GormDB() *gorm.DB

	Close() error
}
