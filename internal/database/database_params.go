package database

import (
	"errors"
	"fmt"
	"strings"

	mysqlp "github.com/brunojet/go-infra-backend/internal/database/adapters/mysql"
	"github.com/brunojet/go-infra-backend/internal/database/types"
)

var (
	ErrUnsupportedDBDriver = errors.New("unsupported database driver")
	ErrMissingMySQLConfig  = errors.New("missing mysql configuration")
)

// DatabaseParams describes how to connect to a database.
//
// Notes:
//   - If DSN is provided, it takes precedence over the other fields.
//   - Driver must be set to a supported value (e.g. "mysql" or "sqlite").
//     Normalization/defaulting from env strings should happen at the app config layer.
//   - For sqlite, DSN can be empty to use in-memory.
//   - For mysql, either DSN must be provided or Host/User/Database must be set.
//
// Migrate is intentionally kept here as a convenience flag for higher-level bootstraps;
// the actual migration execution is performed by DatabaseManager via a provided callback.
type DatabaseParams struct {
	Driver types.DBDriver

	// DSN is the driver-specific connection string.
	DSN string

	// Common structured fields (used to build DSN when DSN is empty).
	Host     string
	Port     int
	Database string
	Schema   string
	User     string
	Password string

	// Options are driver-specific query params.
	// For mysql, this maps to go-sql-driver/mysql Config.Params.
	Options map[string]string

	Migrate bool
}

// BuildDSN returns a driver-specific DSN.
//
// If p.DSN is non-empty, it is returned as-is.
func (p DatabaseParams) BuildDSN() (string, error) {
	if strings.TrimSpace(p.DSN) != "" {
		return strings.TrimSpace(p.DSN), nil
	}

	switch p.Driver {
	case types.DBDriverSQLite:
		// The sqlite connector already falls back to in-memory when DSN is empty.
		return "", nil
	case types.DBDriverMySQL:
		return p.buildMySQLDSN()
	default:
		return "", fmt.Errorf("%w: %q", ErrUnsupportedDBDriver, p.Driver)
	}
}

func (p DatabaseParams) buildMySQLDSN() (string, error) {
	dsn, err := mysqlp.BuildDSN(mysqlp.DSNParams{
		Host:     p.Host,
		Port:     p.Port,
		Database: p.Database,
		User:     p.User,
		Password: p.Password,
		Options:  p.Options,
	})
	if err != nil {
		// Keep legacy error type for callers/tests that might inspect it.
		if errors.Is(err, mysqlp.ErrMissingConfig) {
			return "", fmt.Errorf("%w: %v", ErrMissingMySQLConfig, err)
		}
		return "", err
	}
	return dsn, nil
}
