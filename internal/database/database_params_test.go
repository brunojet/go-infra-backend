package database

import (
	"errors"
	"strings"
	"testing"

	mysqlp "github.com/brunojet/go-infra-backend/internal/database/adapters/mysql"
	"github.com/brunojet/go-infra-backend/internal/database/types"
)

func TestDatabaseParams_BuildDSN_EmptyDriverIsUnsupported(t *testing.T) {
	p := DatabaseParams{}
	_, err := p.BuildDSN()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrUnsupportedDBDriver) {
		t.Fatalf("expected ErrUnsupportedDBDriver, got %v", err)
	}
}

func TestDatabaseParams_BuildDSN_PrefersExplicitDSN(t *testing.T) {
	p := DatabaseParams{Driver: types.DBDriverMySQL, DSN: "  some-dsn  "}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dsn != "some-dsn" {
		t.Fatalf("expected trimmed DSN, got %q", dsn)
	}
}

func TestDatabaseParams_BuildDSN_SqliteEmpty(t *testing.T) {
	p := DatabaseParams{Driver: types.DBDriverSQLite}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dsn != "" {
		t.Fatalf("expected empty DSN for sqlite (connector handles in-memory), got %q", dsn)
	}
}

func TestDatabaseParams_BuildDSN_MySQLStructured(t *testing.T) {
	p := DatabaseParams{
		Driver:   types.DBDriverMySQL,
		Host:     "db",
		Port:     3307,
		Database: "store",
		User:     "root",
		Password: "secret",
	}
	dsn, err := p.BuildDSN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(dsn, "tcp(db:3307)") {
		t.Fatalf("expected tcp host:port in DSN, got %q", dsn)
	}
	if !strings.Contains(dsn, "/store") {
		t.Fatalf("expected db name in DSN, got %q", dsn)
	}
	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Fatalf("expected charset=utf8mb4 default, got %q", dsn)
	}
	// parseTime is a driver option; some DSN formatters omit it when it matches the default.
}

func TestDatabaseParams_BuildDSN_MySQLMissingUserOrDatabase(t *testing.T) {
	p := DatabaseParams{Driver: types.DBDriverMySQL, Host: "db"}
	_, err := p.BuildDSN()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDatabaseParams_BuildDSN_MySQLError(t *testing.T) {
	old := buildMySQLDSNFunc
	buildMySQLDSNFunc = func(mysqlp.DSNParams) (string, error) {
		return "", errors.New("fail")
	}
	defer func() { buildMySQLDSNFunc = old }()
	p := DatabaseParams{Driver: types.DBDriverMySQL, Host: "db", User: "u", Database: "d"}
	_, err := p.BuildDSN()
	if err == nil || err.Error() != "fail" && !strings.Contains(err.Error(), "fail") {
		t.Fatalf("expected fail error, got %v", err)
	}
}

func TestDatabaseParams_BuildDSN_MySQLError_MissingConfig(t *testing.T) {
	old := buildMySQLDSNFunc
	buildMySQLDSNFunc = func(mysqlp.DSNParams) (string, error) {
		return "", mysqlp.ErrMissingConfig
	}
	defer func() { buildMySQLDSNFunc = old }()
	p := DatabaseParams{Driver: types.DBDriverMySQL, Host: "db", User: "u", Database: "d"}
	_, err := p.BuildDSN()
	if err == nil || !errors.Is(err, ErrMissingMySQLConfig) {
		t.Fatalf("expected ErrMissingMySQLConfig, got %v", err)
	}
}
