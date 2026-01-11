package main

import (
	"reflect"
	"testing"
	"time"

	infraconfig "github.com/brunojet/go-infra-backend/infra/config"
	infradatabase "github.com/brunojet/go-infra-backend/infra/database"
	infrahttp "github.com/brunojet/go-infra-backend/infra/http"
)

type mapSource map[string]string

func (m mapSource) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

var _ infraconfig.ConfigSource = (mapSource)(nil)

func TestServerConfigFromSource_Defaults(t *testing.T) {
	cfg := serverConfigFromSource(mapSource{})
	if cfg.Addr != ":8080" {
		t.Fatalf("expected :8080, got %q", cfg.Addr)
	}
	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("expected 5s, got %s", cfg.ReadHeaderTimeout)
	}
}

func TestServerConfigFromSource_Custom(t *testing.T) {
	cfg := serverConfigFromSource(mapSource{
		"PORT":                "9090",
		"READ_HEADER_TIMEOUT": "2s",
	})
	if cfg.Addr != ":9090" {
		t.Fatalf("expected :9090, got %q", cfg.Addr)
	}
	if cfg.ReadHeaderTimeout != 2*time.Second {
		t.Fatalf("expected 2s, got %s", cfg.ReadHeaderTimeout)
	}
}

func TestHTTPParamsFromSource_Defaults(t *testing.T) {
	p := httpParamsFromSource(mapSource{})
	if p.Driver != infrahttp.HTTPDriverGin {
		t.Fatalf("expected gin, got %q", p.Driver)
	}
	if !p.Register {
		t.Fatalf("expected Register=true")
	}
	if p.CORS.Enabled {
		t.Fatalf("expected CORS disabled by default")
	}
	if !reflect.DeepEqual(p.CORS.AllowOrigins, []string{"*"}) {
		t.Fatalf("expected AllowOrigins [*], got %#v", p.CORS.AllowOrigins)
	}
	if !reflect.DeepEqual(p.CORS.AllowMethods, []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}) {
		t.Fatalf("expected default AllowMethods, got %#v", p.CORS.AllowMethods)
	}
	if !reflect.DeepEqual(p.CORS.AllowHeaders, []string{"Content-Type", "Authorization"}) {
		t.Fatalf("expected default AllowHeaders, got %#v", p.CORS.AllowHeaders)
	}
	if p.CORS.ExposeHeaders != nil {
		t.Fatalf("expected nil ExposeHeaders by default, got %#v", p.CORS.ExposeHeaders)
	}
	if p.CORS.AllowCredentials {
		t.Fatalf("expected AllowCredentials=false")
	}
	if p.CORS.MaxAge != 10*time.Minute {
		t.Fatalf("expected MaxAge=10m, got %s", p.CORS.MaxAge)
	}
}

func TestHTTPParamsFromSource_Overrides(t *testing.T) {
	p := httpParamsFromSource(mapSource{
		"HTTP_DRIVER":            "CHI",
		"CORS_ENABLED":           "0",
		"CORS_ALLOW_ORIGINS":     " https://a.com, https://b.com ",
		"CORS_ALLOW_METHODS":     "get,post",
		"CORS_ALLOW_HEADERS":     "x-a, x-b",
		"CORS_EXPOSE_HEADERS":    "x-out",
		"CORS_ALLOW_CREDENTIALS": "true",
		"CORS_MAX_AGE":           "30s",
	})
	if p.Driver != infrahttp.HTTPDriverChi {
		t.Fatalf("expected chi, got %q", p.Driver)
	}
	if p.CORS.Enabled {
		t.Fatalf("expected CORS disabled")
	}
	if !reflect.DeepEqual(p.CORS.AllowOrigins, []string{"https://a.com", "https://b.com"}) {
		t.Fatalf("unexpected AllowOrigins: %#v", p.CORS.AllowOrigins)
	}
	if !reflect.DeepEqual(p.CORS.AllowMethods, []string{"GET", "POST"}) {
		t.Fatalf("unexpected AllowMethods: %#v", p.CORS.AllowMethods)
	}
	if !reflect.DeepEqual(p.CORS.AllowHeaders, []string{"x-a", "x-b"}) {
		t.Fatalf("unexpected AllowHeaders: %#v", p.CORS.AllowHeaders)
	}
	if !reflect.DeepEqual(p.CORS.ExposeHeaders, []string{"x-out"}) {
		t.Fatalf("unexpected ExposeHeaders: %#v", p.CORS.ExposeHeaders)
	}
	if !p.CORS.AllowCredentials {
		t.Fatalf("expected AllowCredentials=true")
	}
	if p.CORS.MaxAge != 30*time.Second {
		t.Fatalf("expected MaxAge=30s, got %s", p.CORS.MaxAge)
	}
}

func TestDatabaseParamsFromSource_Defaults(t *testing.T) {
	p := databaseParamsFromSource(mapSource{})
	if p.Driver != infradatabase.DBDriverMySQL {
		t.Fatalf("expected mysql, got %q", p.Driver)
	}
	if p.Migrate {
		t.Fatalf("expected Migrate=false")
	}
	if p.Port != 0 {
		t.Fatalf("expected Port=0, got %d", p.Port)
	}
	if p.Options != nil {
		t.Fatalf("expected nil Options, got %#v", p.Options)
	}
}

func TestDatabaseParamsFromSource_FallbackDatabaseName(t *testing.T) {
	p := databaseParamsFromSource(mapSource{
		"DB_DATABASE": "store",
	})
	if p.Database != "store" {
		t.Fatalf("expected store, got %q", p.Database)
	}
}

func TestDatabaseParamsFromSource_OptionsAndMigrate(t *testing.T) {
	p := databaseParamsFromSource(mapSource{
		"DB_DRIVER":  "sqlite",
		"DB_OPTIONS": "charset=utf8mb4,loc=UTC,invalid",
		"DB_MIGRATE": "false",
	})
	if p.Driver != infradatabase.DBDriverSQLite {
		t.Fatalf("expected sqlite, got %q", p.Driver)
	}
	if p.Migrate {
		t.Fatalf("expected Migrate=false")
	}
	want := map[string]string{"charset": "utf8mb4", "loc": "UTC"}
	if !reflect.DeepEqual(p.Options, want) {
		t.Fatalf("expected %#v, got %#v", want, p.Options)
	}
}

func TestObservabilityParamsFromSource_Defaults(t *testing.T) {
	p := observabilityParamsFromSource(mapSource{}, "chi")
	if p.Driver != "chi" {
		t.Fatalf("expected driver chi, got %q", p.Driver)
	}
	if !p.Config.RequestID || !p.Config.Telemetry || !p.Config.Recovery {
		t.Fatalf("expected all defaults true, got %+v", p.Config)
	}
}

func TestObservabilityParamsFromSource_TelemetryFallbackToAccessLog(t *testing.T) {
	p := observabilityParamsFromSource(mapSource{
		"OBS_ACCESS_LOG": "0",
	}, "gin")
	if p.Config.Telemetry {
		t.Fatalf("expected Telemetry=false when OBS_ACCESS_LOG=0 and OBS_TELEMETRY not set")
	}
}

func TestObservabilityParamsFromSource_TelemetryWinsWhenExplicit(t *testing.T) {
	p := observabilityParamsFromSource(mapSource{
		"OBS_TELEMETRY":  "0",
		"OBS_ACCESS_LOG": "1",
	}, "gin")
	if p.Config.Telemetry {
		t.Fatalf("expected Telemetry=false when OBS_TELEMETRY=0 even if OBS_ACCESS_LOG=1")
	}
}
