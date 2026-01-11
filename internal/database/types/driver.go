package types

import "strings"

// DBDriver identifies the database adapter/connector implementation.
//
// It is intentionally string-based to keep env/config wiring simple.
type DBDriver string

const (
	DBDriverMySQL  DBDriver = "mysql"
	DBDriverSQLite DBDriver = "sqlite"
)

// NormalizeDBDriver trims/normalizes the input and applies the default.
//
// Empty values default to mysql.
func NormalizeDBDriver(raw string) DBDriver {
	v := DBDriver(strings.TrimSpace(strings.ToLower(raw)))
	if v == "" {
		return DBDriverMySQL
	}
	return v
}

func (d DBDriver) IsSupported() bool {
	switch d {
	case DBDriverMySQL, DBDriverSQLite:
		return true
	default:
		return false
	}
}
