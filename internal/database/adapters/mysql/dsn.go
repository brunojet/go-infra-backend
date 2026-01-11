package mysql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	mysqlcfg "github.com/go-sql-driver/mysql"
)

var ErrMissingConfig = errors.New("missing mysql configuration")

// DSNParams are the minimum common parameters to build a MySQL DSN.
//
// If you already have a DSN, you can bypass this and pass it directly to New(dsn).
type DSNParams struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string

	// Options map to go-sql-driver/mysql Config.Params.
	Options map[string]string
}

// BuildDSN validates the minimum required configuration and returns a MySQL DSN.
//
// Requirements:
//   - User and Database must be provided.
//
// Defaults:
//   - Host: localhost
//   - Port: 3306
//   - parseTime: true (can be overridden via Options["parseTime"])
//   - loc: Local (can be overridden via Options["loc"], e.g. "UTC")
func BuildDSN(p DSNParams) (string, error) {
	host := strings.TrimSpace(p.Host)
	if host == "" {
		host = "localhost"
	}

	port := p.Port
	if port == 0 {
		port = 3306
	}

	user := strings.TrimSpace(p.User)
	db := strings.TrimSpace(p.Database)
	if user == "" || db == "" {
		return "", fmt.Errorf("%w: mysql requires at least User and Database", ErrMissingConfig)
	}

	params := map[string]string{}
	for k, v := range p.Options {
		kk := strings.TrimSpace(k)
		if kk == "" {
			continue
		}
		params[kk] = strings.TrimSpace(v)
	}

	// Good defaults for PT-BR and general Unicode usage.
	// Only set if not provided by Options.
	if _, ok := params["charset"]; !ok {
		params["charset"] = "utf8mb4"
	}
	// MySQL 5.7+ generally supports this; callers can override if needed.
	if _, ok := params["collation"]; !ok {
		params["collation"] = "utf8mb4_unicode_ci"
	}

	// Allow overriding defaults via Options.
	parseTime := true
	if v, ok := params["parseTime"]; ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("invalid mysql option parseTime=%q", v)
		}
		parseTime = b
		delete(params, "parseTime")
	}

	loc := time.Local
	if v, ok := params["loc"]; ok {
		// go-sql-driver/mysql accepts IANA names (e.g., "UTC") and "Local".
		l, err := time.LoadLocation(v)
		if err != nil {
			return "", fmt.Errorf("invalid mysql option loc=%q", v)
		}
		loc = l
		delete(params, "loc")
	}

	cfg := mysqlcfg.Config{
		User:      user,
		Passwd:    p.Password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%d", host, port),
		DBName:    db,
		Params:    params,
		ParseTime: parseTime,
		Loc:       loc,
	}

	return cfg.FormatDSN(), nil
}
