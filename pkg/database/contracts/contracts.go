package contracts

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

const (
	DB_DRIVER_ENV = "DB_DRIVER"
	// DB_ENDPOINT_ENV is kept for backwards compatibility.
	// Deprecated: prefer DB_HOST_ENV + DB_PORT_ENV.
	DB_ENDPOINT_ENV             = "DB_ENDPOINT"
	DB_HOST_ENV                 = "DB_HOST"
	DB_PORT_ENV                 = "DB_PORT"
	DB_SCHEMA_ENV               = "DB_SCHEMA"
	DB_NAME_ENV                 = "DB_NAME"
	DB_MODE_ENV                 = "DB_MODE"
	DB_MAX_OPEN_CONNECTIONS_ENV = "DB_MAX_OPEN_CONNECTIONS"
	DB_MAX_IDLE_CONNECTIONS_ENV = "DB_MAX_IDLE_CONNECTIONS"
	DB_CONN_MAX_LIFETIME_ENV    = "DB_CONN_MAX_LIFETIME"
	DB_CONN_MAX_IDLE_TIME_ENV   = "DB_CONN_MAX_IDLE_TIME"
)

type DatabaseDriver string

const (
	DbDriverSQLiteMemory DatabaseDriver = "sqlite_memory"
	DbDriverSQLiteDisk   DatabaseDriver = "sqlite_disk"
	DbDriverPostgres     DatabaseDriver = "postgres"
	DbDriverMySQL        DatabaseDriver = "mysql"
)

type DatabaseMode string

const (
	DatabaseModeMemory  DatabaseMode = "memory"
	DatabaseModeDisk    DatabaseMode = "disk"
	DatabaseModeDefault DatabaseMode = DatabaseModeMemory
)

const (
	DbDriverDefault              = DbDriverSQLiteMemory
	DbMaxOpenConnectionsDefault  = -1
	DbMaxIdleConnectionsDefault  = -1
	DbConnMaxLifetimeSecsDefault = -1
	DbConnMaxIdleTimeSecsDefault = -1
)

var (
	ErrUnsupportedDriver = fmt.Errorf("unsupported database driver")
)

// DatabaseConfig is a generic database configuration object.
// Specific adapters may interpret Schema/Name differently.
type DatabaseConfig struct {
	Driver              DatabaseDriver
	Host                string
	Port                string
	Schema              string
	Name                string
	Mode                DatabaseMode
	MaxOpenConnections  int
	MaxIdleConnections  int
	ConnMaxLifetimeSecs int
	ConnMaxIdleTimeSecs int
}

// DatabaseAdapter represents a lifecycle-managed database instance.
// Implementations should provide access to GORM and to the underlying SQL DB.
type DatabaseAdapter interface {
	GormDB() (*gorm.DB, error)
	SqlDB() (*sql.DB, error)
	Close() error
}

type DatabaseManager interface {
	DatabaseAdapter() DatabaseAdapter
	Migrate(dst ...interface{}) error
	HealthCheck(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
