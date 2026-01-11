package main

import (
	"time"

	infraconfig "github.com/brunojet/go-infra-backend/infra/config"
	infradatabase "github.com/brunojet/go-infra-backend/infra/database"
	infrahttp "github.com/brunojet/go-infra-backend/infra/http"
	infraobs "github.com/brunojet/go-infra-backend/infra/observability"
	infraserver "github.com/brunojet/go-infra-backend/infra/server"
)

type appConfig struct {
	Server        infraserver.ServerParams
	HTTP          infrahttp.HTTPParams
	Observability infraobs.ObservabilityParams
	Database      infradatabase.DatabaseParams
}

func configFromEnv() appConfig {
	src := infraconfig.NewEnvSource()
	httpParams := httpParamsFromSource(src)
	return appConfig{
		Server:        serverConfigFromSource(src),
		HTTP:          httpParams,
		Observability: observabilityParamsFromSource(src, httpParams.Driver),
		Database:      databaseParamsFromSource(src),
	}
}

// Env (Server):
//   - PORT (default: 8080)
//   - READ_HEADER_TIMEOUT (default: 5s)
func serverConfigFromSource(src infraconfig.ConfigSource) infraserver.ServerParams {
	port := infraconfig.Trimmed(src, "PORT")
	if port == "" {
		port = "8080"
	}
	return infraserver.ServerParams{
		Addr:              ":" + port,
		ReadHeaderTimeout: infraconfig.Duration(src, "READ_HEADER_TIMEOUT", 5*time.Second),
	}
}

// Env (Observability):
//   - OBS_REQUEST_ID (default: true)
//   - OBS_ACCESS_LOG (default: false)
//   - OBS_TELEMETRY (default: true)
//   - if OBS_TELEMETRY is not set, OBS_ACCESS_LOG acts as a backward-compatible alias
//     for telemetry enable/disable.
//   - OBS_RECOVERY (default: true)
func observabilityParamsFromSource(src infraconfig.ConfigSource, httpDriver infrahttp.HTTPDriver) infraobs.ObservabilityParams {
	cfg := infraobs.MiddlewareConfig{RequestID: true, AccessLog: false, Telemetry: true, Recovery: true}

	cfg.RequestID = infraconfig.Bool(src, "OBS_REQUEST_ID", cfg.RequestID)
	cfg.AccessLog = infraconfig.Bool(src, "OBS_ACCESS_LOG", cfg.AccessLog)

	// Backward-compatible: OBS_ACCESS_LOG previously controlled request logs.
	// If OBS_TELEMETRY is explicitly set, it wins.
	if _, ok := src.Lookup("OBS_TELEMETRY"); ok {
		cfg.Telemetry = infraconfig.Bool(src, "OBS_TELEMETRY", cfg.Telemetry)
	} else {
		cfg.Telemetry = infraconfig.Bool(src, "OBS_ACCESS_LOG", cfg.Telemetry)
	}

	cfg.Recovery = infraconfig.Bool(src, "OBS_RECOVERY", cfg.Recovery)

	return infraobs.ObservabilityParams{Driver: httpDriver, Config: cfg}
}

// Env (HTTP):
//   - HTTP_DRIVER (default: "gin")
//   - CORS_ENABLED (default: false)
//   - CORS_ALLOW_ORIGINS (default: "*")
//   - CORS_ALLOW_METHODS (default: "GET,POST,PATCH,DELETE,OPTIONS")
//   - CORS_ALLOW_HEADERS (default: "Content-Type,Authorization")
//   - CORS_EXPOSE_HEADERS (default: "")
//   - CORS_ALLOW_CREDENTIALS (default: false)
//   - CORS_MAX_AGE (default: 10m)
func httpParamsFromSource(src infraconfig.ConfigSource) infrahttp.HTTPParams {
	// Keep defaults aligned with infra/http/config.go.
	cors := infrahttp.CORSConfig{
		Enabled:          infraconfig.Bool(src, "CORS_ENABLED", false),
		AllowOrigins:     infraconfig.SplitCSV(infraconfig.Trimmed(src, "CORS_ALLOW_ORIGINS"), []string{"*"}),
		AllowMethods:     infraconfig.SplitCSVUpper(infraconfig.Trimmed(src, "CORS_ALLOW_METHODS"), []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		AllowHeaders:     infraconfig.SplitCSV(infraconfig.Trimmed(src, "CORS_ALLOW_HEADERS"), []string{"Content-Type", "Authorization"}),
		ExposeHeaders:    infraconfig.SplitCSV(infraconfig.Trimmed(src, "CORS_EXPOSE_HEADERS"), nil),
		AllowCredentials: infraconfig.Bool(src, "CORS_ALLOW_CREDENTIALS", false),
		MaxAge:           infraconfig.Duration(src, "CORS_MAX_AGE", 10*time.Minute),
	}

	driver := infrahttp.NormalizeHTTPDriver(infraconfig.Trimmed(src, "HTTP_DRIVER"))

	return infrahttp.HTTPParams{Driver: driver, CORS: cors, Register: true}
}

// Env:
//   - DB_DRIVER (default: "mysql")
//   - DB_DSN (driver-specific)
//   - DB_HOST
//   - DB_PORT
//   - DB_NAME (or DB_DATABASE)
//   - DB_USER
//   - DB_PASSWORD
//   - DB_SCHEMA
//   - DB_OPTIONS (csv: "k=v,k2=v2")
//   - DB_MIGRATE (default: false)
func databaseParamsFromSource(src infraconfig.ConfigSource) infradatabase.DatabaseParams {
	driver := infradatabase.NormalizeDBDriver(infraconfig.Trimmed(src, "DB_DRIVER"))

	dbName := infraconfig.Trimmed(src, "DB_NAME")
	if dbName == "" {
		dbName = infraconfig.Trimmed(src, "DB_DATABASE")
	}

	params := infradatabase.DatabaseParams{
		Driver:   driver,
		DSN:      infraconfig.Trimmed(src, "DB_DSN"),
		Host:     infraconfig.Trimmed(src, "DB_HOST"),
		Port:     infraconfig.Int(src, "DB_PORT", 0),
		Database: dbName,
		Schema:   infraconfig.Trimmed(src, "DB_SCHEMA"),
		User:     infraconfig.Trimmed(src, "DB_USER"),
		Password: infraconfig.Trimmed(src, "DB_PASSWORD"),
		Options:  infraconfig.ParseKeyValueCSV(infraconfig.Trimmed(src, "DB_OPTIONS")),
		Migrate:  infraconfig.Bool(src, "DB_MIGRATE", false),
	}

	return params
}

// Funções de helpers de config para testes
var (
	TestServerConfigFromSource        = serverConfigFromSource
	TestHTTPParamsFromSource          = httpParamsFromSource
	TestDatabaseParamsFromSource      = databaseParamsFromSource
	TestObservabilityParamsFromSource = observabilityParamsFromSource
)
